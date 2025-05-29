package adminSvc

import (
	"context"

	"github.com/Daniel-Njaramba-1/pulse/internal/repo"
	"github.com/jmoiron/sqlx"
)

type CustomerService struct {
	db *sqlx.DB
}

func NewCustomerService(db *sqlx.DB) *CustomerService {
	return &CustomerService{db: db}
}

func (s *CustomerService) GetAllCustomers(ctx context.Context) ([]*repo.Customer, error) {
	var customers []*repo.Customer
	query := `
		SELECT id, username, email, is_active
		FROM customers
	`
	err := s.db.SelectContext(ctx, &customers, query)
	if err != nil {
		return nil, err
	}
	return customers, nil
}
