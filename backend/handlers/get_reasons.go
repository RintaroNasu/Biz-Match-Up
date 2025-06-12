package handlers

import (
	"backend/models"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetReasons(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		const bearer = "Bearer "
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, bearer) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing bearer token"})
		}

		tokenStr := strings.TrimPrefix(auth, bearer)
		claims := jwt.MapClaims{}
		secret := []byte(os.Getenv("JWT_SECRET"))

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !tkn.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		uid, ok := claims["id"].(float64)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid claims"})
		}

		var reasons []models.Reason
		if err := db.Where("user_id = ?", uint(uid)).Find(&reasons).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch reasons"})
		}

		return c.JSON(http.StatusOK, reasons)
	}
}
