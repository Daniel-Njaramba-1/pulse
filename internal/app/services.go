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
}

type CustomerServices struct {
	authentication *customerSvc.Authentication
}

func NewAdminServices(db *sqlx.DB) *AdminServices {
	authentication := adminSvc.NewAuthentication(db)
	brandService := adminSvc.NewBrandService(db)
	categoryService := adminSvc.NewCategoryService(db)
	productService := adminSvc.NewProductService(db, categoryService, brandService)

	return &AdminServices{
		authentication: authentication,
		brandService: brandService,
		categoryService: categoryService,
		productService: productService,
	}
}

func NewCustomerServices(db *sqlx.DB) *CustomerServices {
	return &CustomerServices{
		authentication: customerSvc.NewAuthentication(db),
	}
}
