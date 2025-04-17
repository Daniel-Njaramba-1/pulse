package customerSvc

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type WishlistService struct {
	db *sqlx.DB
}

func NewWishlistService(db *sqlx.DB) *WishlistService {
	return &WishlistService{db: db}
}

func (s *WishlistService) AddToWishlist (ctx context.Context, productId int, userId int) error {
	
	return nil
}