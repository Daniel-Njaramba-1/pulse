package db

import (
	"fmt"
	"log"

	"github.com/Daniel-Njaramba-1/pulse/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // _ is used because import is not used explicitly in code but required by package sqlx
)

// DBConfig holds the database configuration
type DBConfig struct {
    Host     string
    User     string
    Password string
    Dbname   string
}

// LoadDBConfig loads the database configuration from environment variables based on the mode.
func LoadDBConfig() (*DBConfig, error) {
    log.Printf("loading DB config")
    return &DBConfig{
        Host:     config.GetEnv("DB_HOST"),
        User:     config.GetEnv("DB_USER"),
        Password: config.GetEnv("DB_PASSWORD"),
        Dbname:   config.GetEnv("DB_NAME"),
    }, nil
}

// InitDB initializes the database connection using the provided configuration.
func InitDB(cfg *DBConfig) (*sqlx.DB, error) {
    connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.User,cfg.Password, cfg.Dbname)

    db, err := sqlx.Connect("postgres", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
    }
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
    }
    log.Println("Connected to Postgres")
    return db, nil
}

// ConnDB loads the configuration and initializes the database connection based on the mode.
func ConnDB() (*sqlx.DB, error) {
    cfg, err := LoadDBConfig()
    if err != nil {
        return nil, fmt.Errorf("load config failed: %v", err)
    }
    db, err := InitDB(cfg)
    if err != nil {
        return nil, fmt.Errorf("connection to Postgres failed: %v", err)
    }
    return db, nil
}

// CloseDB closes the database connection.
func CloseDB(db *sqlx.DB) {
    if err := db.Close(); err != nil {
        log.Printf("error closing database: %v", err)
    }
}

