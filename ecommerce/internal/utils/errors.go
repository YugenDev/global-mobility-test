package utils

import "errors"

var (
	ErrProductNameRequired        = errors.New("product name is required")
	ErrProductDescriptionRequired = errors.New("product description is required")
	ErrProductPriceInvalid        = errors.New("product price must be greater than zero")
	ErrProductStockInvalid        = errors.New("product stock cannot be negative")
	ErrProductIDRequired          = errors.New("product ID is required")
	ErrNoProductsFound            = errors.New("no products found")
	ErrInternalServer             = errors.New("internal server error")
	ErrInvalidRequestPayload      = errors.New("invalid request body")
	ErrProductIDAlreadyExists     = errors.New("product ID already exists")
	ErrProductIDCannotBeChanged   = errors.New("product ID cannot be changed")
	ErrDatabaseNotInitialized     = errors.New("database collection not initialized")
	ErrNullProductData            = errors.New("product data cannot be nil")
)
