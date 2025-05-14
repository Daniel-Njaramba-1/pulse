package customerSvc

import (
	"context"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/hashing"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/logging"
	"github.com/jmoiron/sqlx"
)

type Authentication struct {
	db *sqlx.DB
}

func NewAuthentication(db *sqlx.DB) *Authentication {
	return &Authentication{db: db}
}

func (a *Authentication) RegisterCustomer(ctx context.Context, customer *repo.Customer) (string, *repo.Customer, error) {
	if customer.Username == "" || customer.Email == "" || customer.Password == "" {
		return "", nil, errors.New("missing required fields")
	}

	hashedPassword, err := hashing.HashPassword(customer.Password)
	if err != nil {
		return "", nil, err
	}
	customer.PasswordHash = hashedPassword
	customer.IsActive = true
	customer.Password = ""

	// Begin transaction
	tx, err := a.db.BeginTxx(ctx, nil)
	logging.LogInfo("Transaction started")
	if err != nil {
		return "", nil, err
	}
	defer tx.Rollback()

	// Create all records with a single query using CTE (Common Table Expressions)
	query := `
		WITH new_customer AS (
			INSERT INTO customers (username, email, password_hash, is_active, is_email_verified)
			VALUES (:username, :email, :password_hash, :is_active, FALSE)
			RETURNING id
		),
		new_profile AS (
			INSERT INTO customer_profiles (customer_id, firstname, lastname, phone, address)
			SELECT id, NULL, NULL, NULL, NULL FROM new_customer
		),
		new_cart AS (
			INSERT INTO carts (customer_id, is_active)
			SELECT id, TRUE FROM new_customer
		),
		new_wishlist AS (
			INSERT INTO wishlists (customer_id, is_active)
			SELECT id, TRUE FROM new_customer
		)
		SELECT id FROM new_customer;
	`
	
	rows, err := tx.NamedQuery(query, customer)
	if err != nil {
		logging.LogError("Error executing registration: %v", err)
		return "", nil, err
	}
	defer rows.Close()
	
	if rows.Next() {
		if err := rows.Scan(&customer.Id); err != nil {
			return "", nil, err
		}
	} else {
		return "", nil, errors.New("failed to create customer")
	}
	
	// Commit transaction
	if err = tx.Commit(); err != nil {
		return "", nil, err
	}

	// Generate token
	token, err := CreateCustomerToken(customer.Id, customer.Username)
	if err != nil {
		return "", nil, err
	}

	return token, customer, nil
}


func (a *Authentication) LoginCustomer(ctx context.Context, username string, password string) (string, *repo.Customer, error) {
	var customer repo.Customer

	query := `SELECT id, username, email, password_hash FROM customers WHERE username = $1 AND is_active = TRUE`
	err := a.db.GetContext(ctx, &customer, query, username)
	if err != nil {
		return "", nil, err
	}

	if !hashing.VerifyPassword(password, customer.PasswordHash) {
		return "", nil, errors.New("invalid password")
	}

	token, err := CreateCustomerToken(customer.Id, customer.Username)
	if err != nil {
		return "", nil, err
	}

	return token, &customer, nil
}

// later implement using github.com/wneessen/go-mail
// func (a *Authentication) SendVerificationEmail(ctx context.Context, email string) (bool, error) {
	
// }

// func (a *Authentication) VerifyEmail(ctx context.Context, email string) (bool, error) {
	
// }

// func (a *Authentication) SendPasswordResetEmail(ctx context.Context, email string) (bool, error) {
	
// }

// func (a *Authentication) ResetCustomerPassword(ctx context.Context, newPassword string) {

// }
