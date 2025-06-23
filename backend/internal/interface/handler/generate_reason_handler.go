package handler

import (
	"backend/internal/domain/model"
	"backend/internal/usecase"
	"backend/internal/util"
	"context"
	"net/http"

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

		uid, err := util.ExtractUserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
		}

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
