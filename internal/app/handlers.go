package app

import (
	"github.com/Daniel-Njaramba-1/pulse/internal/api/adminHdl"
	"github.com/Daniel-Njaramba-1/pulse/internal/api/customerHdl"
)

type AdminHdl struct {
	AuthHandler *adminHdl.AuthHandler
	BrandHandler *adminHdl.BrandHandler
	CategoryHandler *adminHdl.CategoryHandler
	ProductHandler *adminHdl.ProductHandler
	DashboardHandler *adminHdl.DashboardHandler
}

type CustomerHdl struct {
	AuthHandler *customerHdl.AuthHandler
	ProductHandler *customerHdl.ProductHandler
	CartHandler *customerHdl.CartHandler
	OrderHandler *customerHdl.OrderHandler
	PaymentHandler *customerHdl.PaymentHandler
	ReviewHandler *customerHdl.ReviewHandler
	WishlistHandler *customerHdl.WishlistHandler
}

func NewAdminHdl(adminSvc *AdminServices) *AdminHdl {
	return &AdminHdl{
		AuthHandler: adminHdl.NewAuthHandler(adminSvc.authentication),
		BrandHandler: adminHdl.NewBrandHandler(adminSvc.brandService),
		CategoryHandler: adminHdl.NewCategoryHandler(adminSvc.categoryService),
		ProductHandler: adminHdl.NewProductHandler(adminSvc.productService),
		DashboardHandler: adminHdl.NewDashboardHandler(adminSvc.dashboardService),
	}
}

func NewCustomerHdl(customerSvc *CustomerServices) *CustomerHdl {
	return &CustomerHdl{
		AuthHandler: customerHdl.NewAuthHandler(customerSvc.authentication),
		ProductHandler: customerHdl.NewProductHandler(customerSvc.productService),
		CartHandler: customerHdl.NewCartHandler(customerSvc.cartService),
		OrderHandler: customerHdl.NewOrderHandler(customerSvc.orderService),
		PaymentHandler: customerHdl.NewPaymentHandler(customerSvc.paymentService),
		ReviewHandler: customerHdl.NewReviewHandler(customerSvc.reviewService),
		WishlistHandler: customerHdl.NewWishlistHandler(customerSvc.wishlistService),
	}
}
