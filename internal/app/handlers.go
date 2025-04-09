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
}

type CustomerHdl struct {
	AuthHandler *customerHdl.AuthHandler
}

func NewAdminHdl(adminSvc *AdminServices) *AdminHdl {
	return &AdminHdl{
		AuthHandler: adminHdl.NewAuthHandler(adminSvc.authentication),
		BrandHandler: adminHdl.NewBrandHandler(adminSvc.brandService),
		CategoryHandler: adminHdl.NewCategoryHandler(adminSvc.categoryService),
		ProductHandler: adminHdl.NewProductHandler(adminSvc.productService),
	}
}

func NewCustomerHdl(customerSvc *CustomerServices) *CustomerHdl {
	return &CustomerHdl{
		AuthHandler: customerHdl.NewAuthHandler(customerSvc.authentication),
	}
}
