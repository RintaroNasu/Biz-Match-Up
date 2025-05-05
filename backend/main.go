package main

import (
	"backend/db"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	db.NewDB()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Echo + GORM!")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
