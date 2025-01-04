package main

import (
	"github.com/YugenDev/global-mobility-test/internal/config"
	"github.com/YugenDev/global-mobility-test/internal/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	config.ConnectDatabase()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Welcome to Global Mobility Apex ecommerce 🚀")
	})

	routes.ProductRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
