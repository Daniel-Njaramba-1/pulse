package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type Customer struct {
	Id					int			`db:"id" json:"id"`
	Username			string		`db:"username" json:"username"`
	Email				string		`db:"email" json:"email"`
	PasswordHash		string		`db:"password_hash" json:"-"`
    Password        	string      `db:"-" json:"password"`
	IsActive			bool		`db:"is_active" json:"is_active"`
	CreatedAt			time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt			time.Time	`db:"updated_at" json:"updated_at"`
}

func (customer *Customer) FeedGetId() *int {
	return &customer.Id
}

func (customer *Customer) FeedCreateQuery() string {
	return `
		INSERT INTO customers (username, email, password_hash, is_active)
		VALUES (:username, :email, :password_hash, :is_active)
		RETURNING id
	`
}

func (customer *Customer) FeedGetByIdQuery() string {
	return `
		SELECT id, username, email, is_active, created_at, updated_at
		FROM customers
		WHERE id = $1
	`
}

func (customer *Customer) FeedGetAllQuery() string {
	return `
		SELECT id, username, email, is_active, created_at, updated_at
		FROM customers
	`
}

func (customer *Customer) FeedUpdateDetailsQuery() string {
	return `
		UPDATE customers
		SET username = :username,
			email = :email,
		WHERE id = :id
	`
}

func (customer *Customer) FeedDeleteQuery() string {
	return `
		DELETE FROM customers
		WHERE id = :id
	`
}

func (customer *Customer) FeedDeactivateQuery() string {
	return `
		UPDATE customers
		SET is_active = FALSE
		WHERE id = :id
	`
}

func (customer *Customer) FeedReactivateQuery() string {
	return `
		UPDATE customers
		SET is_active = TRUE
		WHERE id = :id
	`
}

func GetCustomerByNameQuery(ctx context.Context, db *sqlx.DB, name string, customer *Customer) error {
	query := `	
		SELECT id, username, email, is_active, created_at, updated_at
		FROM customers 
		WHERE username = $1
	`
	return db.Get(customer, query, name)
}

func SearchCustomerByNameQuery(ctx context.Context, db *sqlx.DB, name string, customer *Customer) error {
	query :=`	
		SELECT id, username, email, is_active, created_at, updated_at
		FROM customers 
		WHERE username ILIKE $1
	`
	return db.Select(customer, query, name)
}

// for comparing with password attempts on login
func GetCustomerPasswordQuery(ctx context.Context, db *sqlx.DB, username string, customer *Customer) error {
	query :=`		
		SELECT password_hash
		FROM customers 
		WHERE username = $1
	`
	return db.Get(customer, query, username)
}

// update password using ID
func UpdateCustomerPasswordQuery(db *sqlx.DB, id int, customer *Customer) error {
	query :=`	
		UPDATE customers
		SET password_hash = :password_hash
		WHERE id = :id
	`
	_, err := db.NamedExec(query, customer)
	return err
}