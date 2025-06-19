package handler

import (
	"backend/internal/domain/model"
	"backend/internal/usecase"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CompanyScrapeHandler struct {
	Usecase *usecase.CompanyScrapeUsecase
}

type CompanyScrapeRequest struct {
	CompanyURL string `json:"companyUrl"`
}

func (h *CompanyScrapeHandler) CompanyScrape(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CompanyScrapeRequest
		if err := c.Bind(&req); err != nil || req.CompanyURL == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		// トークン取得・検証
		auth := c.Request().Header.Get("Authorization")
		const bearer = "Bearer "
		if !strings.HasPrefix(auth, bearer) {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "missing bearer token"})
		}
		tokenStr := strings.TrimPrefix(auth, bearer)

		claims := jwt.MapClaims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !tkn.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		uid, ok := claims["id"].(float64)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid token claims"})
		}
		fmt.Println("User ID:", uid)

		var user model.User
		if err := db.First(&user, uint(uid)).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}

		resp, err := h.Usecase.ScrapeCompanyInfo(context.Background(), req.CompanyURL, &user)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{"success": true, "matchResult": resp})
	}
}
