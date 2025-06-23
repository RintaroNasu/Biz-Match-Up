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

		uid, err := util.ExtractUserIDFromToken(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
		}

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
