package app

import (
	"github.com/Daniel-Njaramba-1/pulse/internal/services/adminSvc"
	"github.com/Daniel-Njaramba-1/pulse/internal/services/customerSvc"
	"github.com/jmoiron/sqlx"
)

type AdminServices struct {
	authentication *adminSvc.Authentication
	brandService *adminSvc.BrandService
	categoryService *adminSvc.CategoryService
	productService *adminSvc.ProductService
	dashboardService *adminSvc.DashboardService
}

type CustomerServices struct {
	authentication *customerSvc.Authentication
	productService *customerSvc.ProductService
	cartService *customerSvc.CartService
	orderService *customerSvc.OrderService
	paymentService *customerSvc.PaymentService
	reviewService *customerSvc.ReviewService
	wishlistService *customerSvc.WishlistService
}

func NewAdminServices(db *sqlx.DB) *AdminServices {
	authentication := adminSvc.NewAuthentication(db)
	brandService := adminSvc.NewBrandService(db)
	categoryService := adminSvc.NewCategoryService(db)
	productService := adminSvc.NewProductService(db, categoryService, brandService)
	dashboardService := adminSvc.NewDashboardService(db)

	return &AdminServices{
		authentication: authentication,
		brandService: brandService,
		categoryService: categoryService,
		productService: productService,
		dashboardService: dashboardService,
	}
}

func NewCustomerServices(db *sqlx.DB) *CustomerServices {
	authentication := customerSvc.NewAuthentication(db)
	productService := customerSvc.NewProductService(db)
	cartService := customerSvc.NewCartService(db)
	orderService := customerSvc.NewOrderService(db)
	paymentService := customerSvc.NewPaymentService(db)
	reviewService := customerSvc.NewReviewService(db)
	wishlistService := customerSvc.NewWishlistService(db)


	return &CustomerServices{
		authentication: authentication,
		productService: productService,
		cartService: cartService,
		orderService: orderService,
		paymentService: paymentService,
		reviewService: reviewService,
		wishlistService: wishlistService,
	}
}
