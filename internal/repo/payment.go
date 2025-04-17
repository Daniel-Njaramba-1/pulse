package repo

import "time"

type PaymentStatus string

const (
	PaymentStatusPending	PaymentStatus = "pending"
	PaymentStatusFailed 	PaymentStatus = "failed"
	PaymentStatusCompleted 	PaymentStatus = "completed"
)

type Payment struct {
	Id				int				`db:"id" json:"id"`
	OrderId			int				`db:"order_id" json:"order_id"`
	PaymentMethod	string			`db:"payment_method" json:"payment_method"`	
	Amount			float64			`db:"amount" json:"amount"`
	Status			PaymentStatus	`db:"status" json:"status"`
	TransactionId	string			`db:"transaction_id" json:"transaction_id"`
	CreatedAt		time.Time		`db:"created_at" json:"created_at"`
	UpdatedAt		time.Time		`db:"updated_at" json:"updated_at"`
}