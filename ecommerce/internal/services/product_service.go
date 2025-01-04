package services

import (
	"github.com/YugenDev/global-mobility-test/internal/models"
	"github.com/YugenDev/global-mobility-test/internal/repositories"
	"github.com/YugenDev/global-mobility-test/internal/utils"
	"github.com/labstack/echo/v4"
)

type ProductService struct {
	Repository *repositories.ProductRepository
}

func NewProductService() *ProductService {
	return &ProductService{
		Repository: repositories.NewProductRepository(),
	}
}

func (s *ProductService) CreateProduct(c echo.Context, product *models.Product) error {
	if product.Name == "" {
		return utils.ErrProductNameRequired
	}
	if product.Description == "" {
		return utils.ErrProductDescriptionRequired
	}
	if product.Price <= 0 {
		return utils.ErrProductPriceInvalid
	}
	if product.Stock < 0 {
		return utils.ErrProductStockInvalid
	}

	if product.ProductID == "" {
		product.ProductID = utils.GenerateUniqueID()
	} else {
		existingProduct, err := s.Repository.GetProductByID(product.ProductID)
		if err == nil && existingProduct.ProductID != "" {
			return utils.ErrProductIDAlreadyExists
		}
	}

	_, err := s.Repository.CreateProduct(c, product)
	return err
}

func (s *ProductService) GetAll() ([]models.Product, error) {
	products, err := s.Repository.GetAllProducts()
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, utils.ErrNoProductsFound
	}
	return products, nil
}

func (s *ProductService) GetByID(id string) (models.Product, error) {
	if id == "" {
		return models.Product{}, utils.ErrProductIDRequired
	}
	return s.Repository.GetProductByID(id)
}

func (s *ProductService) UpdateProduct(c echo.Context, id string, product *models.Product) error {
	if id == "" {
		return utils.ErrProductIDRequired
	}

	existingProduct, err := s.Repository.GetProductByID(id)
	if err != nil {
		return err
	}

	if product.ProductID != "" && product.ProductID != id {
		return utils.ErrProductIDCannotBeChanged
	}

	if product.Name != "" {
		existingProduct.Name = product.Name
	}
	if product.Description != "" {
		existingProduct.Description = product.Description
	}
	if product.Price != 0 {
		if product.Price <= 0 {
			return utils.ErrProductPriceInvalid
		}
		existingProduct.Price = product.Price
	}
	if product.Stock != 0 {
		if product.Stock < 0 {
			return utils.ErrProductStockInvalid
		}
		existingProduct.Stock = product.Stock
	}

	_, err = s.Repository.UpdateProduct(c, id, &existingProduct)
	return err
}

func (s *ProductService) DeleteProduct(c echo.Context, id string) error {
	if id == "" {
		return utils.ErrProductIDRequired
	}

	product, err := s.Repository.GetProductByID(id)
	if err != nil {
		return err
	}
	if product.ProductID == "" {
		return utils.ErrNoProductsFound
	}

	_, err = s.Repository.DeleteProduct(c, id)
	return err
}
