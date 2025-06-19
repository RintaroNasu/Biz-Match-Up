package handler

import (
	"backend/internal/domain/model"
	"backend/internal/usecase"
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type GenerateReasonHandler struct {
	Usecase *usecase.GenerateReasonUsecase
}

func (h *GenerateReasonHandler) GenerateReason(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req model.GenerateReasonsRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request"})
		}

		// JWTからユーザーIDを取得
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
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid claims"})
		}

		// ユーザー取得
		var user model.User
		if err := db.First(&user, uint(uid)).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}

		reason, err := h.Usecase.GenerateReason(context.Background(), &user, req.MatchResult, req.Questions)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"success": true,
			"reason":  reason,
		})
	}
}
