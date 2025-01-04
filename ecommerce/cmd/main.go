package main

import (
	"github.com/YugenDev/global-mobility-test/internal/config"
	"github.com/YugenDev/global-mobility-test/internal/handlers"
	"github.com/YugenDev/global-mobility-test/internal/repositories"
	"github.com/YugenDev/global-mobility-test/internal/routes"
	"github.com/YugenDev/global-mobility-test/internal/services"
	"github.com/labstack/echo/v4"
)

func main() {
	config.ConnectDatabase()
	productRepo := repositories.NewProductRepository()
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to Global Mobility Apex ecommerce 🚀")
	})

	routes.ProductRoutes(e, productHandler)

	e.Logger.Fatal(e.Start(":8080"))
}
