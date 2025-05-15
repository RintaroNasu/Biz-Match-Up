package main

import (
	"backend/db"
	"backend/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := db.NewDB()
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Echo + GORM!")
	})
	apiGroup := e.Group("/api")
	apiGroup.POST("/signup", handlers.SignUp(db))
	apiGroup.POST("/signin", handlers.SignIn(db))
	apiGroup.POST("/company-scrape", handlers.CompanyScrape(db))
	apiGroup.POST("/profile/:id", handlers.EditUserProfile(db))
	
	e.Logger.Fatal(e.Start(":8080"))
}
