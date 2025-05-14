package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Daniel-Njaramba-1/pulse/internal/config"
	"github.com/Daniel-Njaramba-1/pulse/internal/pricing"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// DBConfig holds the database configuration
type DBConfig struct {
	Host     string
	User     string
	Password string
	Dbname   string
}

// LoadDBConfig loads the database configuration from environment variables
func LoadDBConfig() (*DBConfig, error) {
	log.Printf("Loading DB config")
	return &DBConfig{
		Host:     config.GetEnv("DB_HOST"),
		User:     config.GetEnv("DB_USER"),
		Password: config.GetEnv("DB_PASSWORD"),
		Dbname:   config.GetEnv("DB_NAME"),
	}, nil
}

// BuildConnStr generates a connection string from config
func BuildConnStr(cfg *DBConfig) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Dbname)
}

// ClientManager handles SSE client connections
type ClientManager struct {
	clients    map[chan string]bool
	register   chan chan string
	unregister chan chan string
	broadcast  chan string
	mutex      sync.Mutex
}

// NewClientManager creates a new SSE client manager
func NewClientManager() *ClientManager {
	return &ClientManager{
		clients:    make(map[chan string]bool),
		register:   make(chan chan string),
		unregister: make(chan chan string),
		broadcast:  make(chan string),
		mutex:      sync.Mutex{},
	}
}

// Global manager instance
var Manager = NewClientManager()

// Run starts the client manager's main loop
func (m *ClientManager) Run() {
	for {
		select {
		case client := <-m.register:
			m.mutex.Lock()
			m.clients[client] = true
			m.mutex.Unlock()
			log.Println("Registered SSE client")
		case client := <-m.unregister:
			m.mutex.Lock()
			if _, ok := m.clients[client]; ok {
				delete(m.clients, client)
				close(client)
			}
			m.mutex.Unlock()
			log.Println("Unregistered SSE client")
		case msg := <-m.broadcast:
			m.mutex.Lock()
			for client := range m.clients {
				select {
				case client <- msg:
				default:
					close(client)
					delete(m.clients, client)
				}
			}
			m.mutex.Unlock()
		}
	}
}

func (m *ClientManager) RegisterChannel(ch chan string) {
    m.register <- ch
}

func (m *ClientManager) UnregisterChannel(ch chan string) {
	m.unregister <- ch
}

// InitDB initializes the database connection using the provided configuration
func InitDB(ctx context.Context, cfg *DBConfig) (*sqlx.DB, error) {
	connStr := BuildConnStr(cfg)

	db, err := sqlx.ConnectContext(ctx, "postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}
	
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}
	
	log.Println("Connected to Postgres")
	return db, nil
}

// ConnDB loads the configuration and initializes the database connection
func ConnDB(ctx context.Context) (*sqlx.DB, error) {
	cfg, err := LoadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("load config failed: %w", err)
	}
	
	db, err := InitDB(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("connection to Postgres failed: %w", err)
	}
	
	return db, nil
}

// StartSaleListener adjusts price after a sale notification
func StartSaleListener(connStr string, modelService *pricing.ModelService) { // Accept ModelService
	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("Error in sales listener event: %v", err)
			return
		}
	})
	defer listener.Close()

	err := listener.Listen("sale") // Listen on 'sale' channel
	if err != nil {
		log.Printf("Error setting up sales listener: %v", err)
		return
	}
	log.Println("Listening for sales")

	for {
		select {
		case n := <-listener.Notify:
			if n == nil {
				continue
			}
			log.Printf("Received sale notification: %v", n.Extra)

			var sale struct {
				ID        int     `json:"id"`
				ProductID int     `json:"product_id"`
				SalePrice float64 `json:"sale_price"`
				Quantity  int     `json:"quantity"`
				CreatedAt string  `json:"created_at"`
			}

			// Use n.Extra which contains the payload from the trigger
			if err := json.Unmarshal([]byte(n.Extra), &sale); err != nil {
				log.Printf("Error parsing sale notification: %v", err)
				continue
			}

			// Adjust price for the product using the ModelService
			ctx := context.Background()
			newPrice, confidence, err := modelService.AdjustPrice(ctx, sale.ProductID) // Use modelService
			if err != nil {
				log.Printf("Error adjusting price for product %d: %v", sale.ProductID, err)
				continue
			}
			log.Printf("Adjusted price for product %d to %.2f with confidence %.2f", sale.ProductID, newPrice, confidence)

		case <-time.After(90 * time.Second):
			// Ping the listener to keep the connection alive
			go func() {
				if err := listener.Ping(); err != nil {
					log.Printf("Error pinging sales listener: %v", err)
				}
			}()
		}
	}
}

// StartPriceAdjustmentListener initializes the PostgreSQL notification listener
func StartPriceAdjustmentListener(connStr string) {
	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Printf("Error in price adjustments listener event: %v", err)
			return
		}
	})
	defer listener.Close()

	err := listener.Listen("price_adjustment")
	if err != nil {
		log.Printf("Error setting up price adjustments listener: %v", err)
		return
	}
	log.Println("Listening for price adjustments")

	// Database connection for processing notifications
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening DB connection for price adjustment processing: %v", err)
		return
	}
	defer db.Close()

	for {
		select {
		case n := <-listener.Notify:
			if n == nil {
				continue
			}
			
			// Process the notification data
			log.Printf("Received price adjustment notification: %v", n.Extra)
			
			// Parse the original notification
			var adjustment struct {
				ID          int     `json:"id"`
				ProductID   int     `json:"product_id"`
				OldPrice    float64 `json:"old_price"`
				NewPrice    float64 `json:"new_price"`
				CreatedAt   string  `json:"created_at"`
			}
			
			if err := json.Unmarshal([]byte(n.Extra), &adjustment); err != nil {
				log.Printf("Error parsing price adjustment notification: %v", err)
				continue
			}
			
			// Create optimized payload
			optimizedData := struct {
				ProductID    int     `json:"product_id"`
				NewPrice     float64 `json:"new_price"`
				ChangedAt    string  `json:"changed_at"`
				PriceChange  float64 `json:"price_change"`
				ChangeType   string  `json:"change_type"`
				ProductName  string  `json:"product_name"`
			}{
				ProductID:   adjustment.ProductID,
				NewPrice:    adjustment.NewPrice,
				ChangedAt:   adjustment.CreatedAt,
				PriceChange: adjustment.NewPrice - adjustment.OldPrice,
				ChangeType:  getChangeType(adjustment.NewPrice, adjustment.OldPrice),
			}
			
			// Get product name
			err := db.QueryRow("SELECT name FROM products WHERE id = $1", adjustment.ProductID).Scan(&optimizedData.ProductName)
			if err != nil {
				log.Printf("Error fetching product name: %v", err)
				optimizedData.ProductName = fmt.Sprintf("Product #%d", adjustment.ProductID)
			}
			
			// Convert optimized data to JSON
			optimizedJSON, err := json.Marshal(optimizedData)
			if err != nil {
				log.Printf("Error creating optimized payload: %v", err)
				continue
			}
			
			// Broadcast the optimized data
			Manager.broadcast <- string(optimizedJSON)
			
		case <-time.After(90 * time.Second):
			go func() {
				if err := listener.Ping(); err != nil {
					log.Printf("Error pinging price adjustments listener: %v", err)
				}
			}()
		}
	}
}

// Helper function to determine price change type
func getChangeType(newPrice, oldPrice float64) string {
	if newPrice > oldPrice {
		return "increase"
	} else if newPrice < oldPrice {
		return "decrease"
	}
	return "unchanged"
}

// CloseDB safely closes the database connection
func CloseDB(db *sqlx.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}
}