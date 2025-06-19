package main

import (
	"backend/internal/db"
	"backend/internal/infrastructure"
	"backend/internal/interface/handler"
	"backend/internal/usecase"
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
	authRepo := &infrastructure.AuthRepositoryImpl{DB: db}
	userRepo := &infrastructure.UserRepositoryImpl{DB: db}
	companyReasonRepo := &infrastructure.CompanyReasonRepositoryImpl{DB: db}
	companyScrapeRepo := &infrastructure.CompanyScrapeRepositoryImpl{DB: db}
	generateReason := &infrastructure.GenerateReasonRepositoryImpl{DB: db}

	authUsecase := &usecase.AuthUsecase{Repo: authRepo}
	userUsecase := &usecase.UserUsecase{Repo: userRepo}
	companyReasonUsecase := &usecase.CompanyReasonUsecase{Repo: companyReasonRepo}
	companyScrapeUsecase := &usecase.CompanyScrapeUsecase{Repo: companyScrapeRepo}
	generateReasonUsecase := &usecase.GenerateReasonUsecase{Repo: generateReason}

	authHandler := &handler.AuthHandler{Usecase: authUsecase}
	userHandler := &handler.UserHandler{Usecase: userUsecase}
	companyReasonHandler := &handler.ReasonHandler{Usecase: companyReasonUsecase}
	companyScrapeHandler := &handler.CompanyScrapeHandler{Usecase: companyScrapeUsecase}
	generateReasonHandler := &handler.GenerateReasonHandler{Usecase: generateReasonUsecase}

	apiGroup := e.Group("/api")
	apiGroup.POST("/signup", authHandler.SignUp)
	apiGroup.POST("/signin", authHandler.SignIn)

	apiGroup.POST("/company-scrape", companyScrapeHandler.CompanyScrape(db))
	apiGroup.POST("/generate-reasons", generateReasonHandler.GenerateReason(db))
	apiGroup.POST("/reason", companyReasonHandler.SaveCompanyReason)
	apiGroup.GET("/reasons", companyReasonHandler.GetCompanyReasons)
	apiGroup.POST("/profile/:id", userHandler.EditUserProfile)
	apiGroup.GET("/profile/:id", userHandler.GetUserProfile)

	e.Logger.Fatal(e.Start(":8080"))
}
