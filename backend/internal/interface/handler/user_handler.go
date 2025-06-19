package handler

import (
	"backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Usecase *usecase.UserUsecase
}

type UpdateUserProfileRequest struct {
	Name               string `json:"name"`
	DesiredJobType     string `json:"desiredJobType"`
	DesiredLocation    string `json:"desiredLocation"`
	DesiredCompanySize string `json:"desiredCompanySize"`
	CareerAxis1        string `json:"careerAxis1"`
	CareerAxis2        string `json:"careerAxis2"`
	SelfPr             string `json:"selfPr"`
}

func (h *UserHandler) EditUserProfile(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user ID"})
	}

	var req UpdateUserProfileRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
	}

	user, err := h.Usecase.EditUserProfile(id, usecase.UpdateUserProfileRequest(req))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"success": true, "user": user})
}

func (h *UserHandler) GetUserProfile(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user ID"})
	}

	user, err := h.Usecase.GetUserProfile(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{"user": user})
}
