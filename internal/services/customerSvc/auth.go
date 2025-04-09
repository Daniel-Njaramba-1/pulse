package customerSvc

import (
	"context"
	"errors"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/Daniel-Njaramba-1/pulse/internal/util/hashing"
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
	customer.Password = ""

	// Begin transaction
	tx, err := a.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", nil, err
	}
	defer tx.Rollback()

	// Create Customer
	customerQuery := `
		INSERT INTO customers (username, email, password_hash, is_active)
		VALUES(:username, :email, :password_hash, :is_active)
		RETURNING id
	`
	rows, err := tx.NamedQuery(customerQuery, customer)
	if err != nil {
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

	// Create Customer Profile
	profile := &repo.CustomerProfile{
		CustomerId: customer.Id,
		FirstName:  "",
		LastName:   "",
		Phone:      "",
		Address:    "",
	}
	profileQuery := `
		INSERT INTO customer_profiles (customer_id, firstname, lastname, phone, address)
		VALUES (:customer_id, :firstname, :lastname, :phone, :address)
	`
	_, err = tx.NamedExecContext(ctx, profileQuery, profile)
	if err != nil {
		return "", nil, err
	}

	// Create Cart
	cart := &repo.Cart{
		CustomerId: customer.Id,
		IsActive: true,
	}
	cartQuery := `
		INSERT INTO carts (customer_id, is_processed)
		VALUES (:customer_id, :is_processed)
	`
	_, err = tx.NamedExecContext(ctx, cartQuery, cart)
	if err != nil {
		return "", nil, err
	}

	// Create Wishlist
	wishlist := &repo.Wishlist{
		CustomerId: customer.Id,
		IsActive:   true,
	}
	wishlistQuery := `
		INSERT INTO wishlists (customer_id, is_active)
		VALUES (:customer_id, :is_active)
	`
	_, err = tx.NamedExecContext(ctx, wishlistQuery, wishlist)
	if err != nil {
		return "", nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return "", nil, err
	}

	// Generate token
	token, err := CreateCustomerToken(customer.Username)
	if err != nil {
		return "", nil, err
	}

	return token, customer, nil
}


func (a *Authentication) LoginCustomer(ctx context.Context, username string, password string) (string, *repo.Customer, error) {
	var customer repo.Customer
	err := repo.GetCustomerByNameQuery(ctx, a.db, username, &customer)
	if err != nil {
		return "", nil, err
	}

	err = repo.GetCustomerPasswordQuery(ctx, a.db, username, &customer)
	if err != nil {
		return "", nil, err
	}

	if !hashing.VerifyPassword(password, customer.PasswordHash) {
		return "", nil, errors.New("invalid password")
	}

	token, err := CreateCustomerToken(customer.Username)
	if err != nil {
		return "", nil, err
	}

	return token, &customer, nil
}

func (a *Authentication) ResetCustomerPassword(ctx context.Context) {

}
