package routes

import (
	"github.com/YugenDev/global-mobility-test/internal/handlers"
	"github.com/labstack/echo/v4"
)

func ProductRoutes(e *echo.Echo, handler *handlers.ProductHandler) {

	e.POST("/products", handler.CreateProduct)
	e.GET("/products", handler.GetAllProducts)
	e.GET("/products/:id", handler.GetProductByID)
	e.PUT("/products/:id", handler.UpdateProduct)
	e.DELETE("/products/:id", handler.DeleteProduct)
}
