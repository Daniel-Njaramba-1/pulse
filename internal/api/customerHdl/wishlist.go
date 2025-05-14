package customerHdl

import "github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"

type WishlistHandler struct {
	wishlistService *customerSvc.WishlistService
}

func NewWishlistHandler(wishlistService *customerSvc.WishlistService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
	}
}