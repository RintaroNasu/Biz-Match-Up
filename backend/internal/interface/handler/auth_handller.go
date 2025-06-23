package handler

import (
	"backend/internal/domain/model"
	"backend/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	Usecase *usecase.AuthUsecase
}

func (h *AuthHandler) SignUp(c echo.Context) error {
	var req usecase.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}

	res, err := h.Usecase.SignUp(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, model.AuthResponse{
		Message: res.Message,
		User:    res.User,
		Token:   res.Token})
}

func (h *AuthHandler) SignIn(c echo.Context) error {
	var req usecase.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
	}

	res, err := h.Usecase.SignIn(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, model.AuthResponse{
		Message: res.Message,
		User:    res.User,
		Token:   res.Token})
}
