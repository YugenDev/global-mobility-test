package handlers

import (
	"net/http"

	"github.com/YugenDev/global-mobility-test/internal/models"
	"github.com/YugenDev/global-mobility-test/internal/services"
	"github.com/YugenDev/global-mobility-test/internal/utils"
	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	Service services.IProductService
}

func NewProductHandler(service services.IProductService) *ProductHandler {
	return &ProductHandler{
		Service: service,
	}
}

func (h *ProductHandler) GetAllProducts(c echo.Context) error {
	products, err := h.Service.GetAll()
	if err != nil {
		if err == utils.ErrNoProductsFound {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": utils.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductIDRequired.Error()})
	}

	product, err := h.Service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": utils.ErrNoProductsFound.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrInvalidRequestPayload.Error()})
	}

	if product.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductNameRequired.Error()})
	}
	if product.Description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductDescriptionRequired.Error()})
	}
	if product.Price <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductPriceInvalid.Error()})
	}
	if product.Stock < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductStockInvalid.Error()})
	}

	if product.ProductID != "" {
		existingProduct, err := h.Service.GetByID(product.ProductID)
		if err == nil && existingProduct.ProductID != "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductIDAlreadyExists.Error()})
		}
	}

	err := h.Service.CreateProduct(c, &product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": utils.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductIDRequired.Error()})
	}

	existingProduct, err := h.Service.GetByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": utils.ErrNoProductsFound.Error()})
	}

	var product models.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrInvalidRequestPayload.Error()})
	}

	if product.Price < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductPriceInvalid.Error()})
	}
	if product.Stock < 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductStockInvalid.Error()})
	}

	if product.ProductID != "" && product.ProductID != existingProduct.ProductID {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductIDCannotBeChanged.Error()})
	}

	err = h.Service.UpdateProduct(c, id, &product)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": utils.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": utils.ErrProductIDRequired.Error()})
	}

	if _, err := h.Service.GetByID(id); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": utils.ErrNoProductsFound.Error()})
	}

	err := h.Service.DeleteProduct(c, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": utils.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
