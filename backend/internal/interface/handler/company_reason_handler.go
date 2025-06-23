package handler

import (
	"backend/internal/usecase"
	"backend/internal/util"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ReasonHandler struct {
	Usecase *usecase.CompanyReasonUsecase
}

type postReasonsRequest struct {
	Content     string `json:"content"`
	CompanyName string `json:"companyName"`
	CompanyUrl  string `json:"companyUrl"`
}

func (h *ReasonHandler) SaveCompanyReason(c echo.Context) error {
	var req postReasonsRequest
	if err := c.Bind(&req); err != nil || strings.TrimSpace(req.Content) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}

	uid, err := util.ExtractUserIDFromToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	err = h.Usecase.SaveCompanyReason(usecase.CompanyReasonRequest{
		Content:     req.Content,
		UserID:      uint(uid),
		CompanyName: req.CompanyName,
		CompanyUrl:  req.CompanyUrl,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to save reason"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "志望理由を保存しました",
	})
}

func (h *ReasonHandler) GetCompanyReasons(c echo.Context) error {
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
	fmt.Println("uid:", uid)
	reasons, err := h.Usecase.GetCompanyReasons(uint(uid))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to fetch reasons"})
	}
	fmt.Println("取得した理由:", reasons)
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"reasons": reasons,
	})
}
