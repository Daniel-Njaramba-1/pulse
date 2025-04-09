package repo

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type CustomerProfile struct {
	Id					int			`db:"id" json:"id"`
	CustomerId			int			`db:"customer_id" json:"customer_id"`
	FirstName			string		`db:"firstname" json:"firstname"`
	LastName			string		`db:"lastname" json:"lastname"`
	Phone				string		`db:"phone" json:"phone"`
	Address				string		`db:"address" json:"address"`
	CreatedAt			time.Time	`db:"created_at" json:"created_at"`
	UpdatedAt			time.Time	`db:"updated_at" json:"updated_at"`
}


func (customerProfile *CustomerProfile) FeedGetId() *int {
	return &customerProfile.Id
}

func (customerProfile *CustomerProfile) FeedCreateQuery() string {
	return `
		INSERT INTO customer_profiles (customer_id, firstname, lastname, phone, address)
		VALUES (:customer_id, :firstname, :lastname, :phone, :address)
		RETURNING id
	`
}

func (customerProfile *CustomerProfile) FeedGetByIdQuery() string {
	return `
		SELECT id, customer_id, firstname, lastname, phone, address, created_at, updated_at
		FROM customer_profiles
		WHERE id = $1
	`
}

func (customerProfile *CustomerProfile) FeedGetAllQuery() string {
	return `
		SELECT id, customer_id, firstname, lastname, phone, address, created_at, updated_at
		FROM customer_profiles
	`
}

func (customerProfile *CustomerProfile) FeedUpdateDetailsQuery() string {
	return `
		UPDATE customer_profiles
		SET customer_id = :customer_id,
			firstname = :firstname,
			lastname = :lastname,
			phone = :phone,
			address = :address,
		WHERE id = :id
	`
}

func (customerProfile *CustomerProfile) FeedDeleteQuery() string {
	return `
		DELETE FROM customer_profiles
		WHERE id = :id
	`
}

func GetCustomerProfileByCustomerQuery(ctx context.Context, db *sqlx.DB, name string, customerProfile *CustomerProfile) error {
	query := `	
		SELECT id, customer_id, firstname, lastname, phone, address, created_at, updated_at
		FROM customer_profiles 
		WHERE customer_id = $1
	`
	return db.Get(customerProfile, query, name)
}



