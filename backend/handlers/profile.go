package handlers

import (
  "fmt"
  "net/http"
  "strconv"

  "backend/models"

  "github.com/labstack/echo/v4"
  "gorm.io/gorm"
)

type UserProfileUpdate struct {
  Name               string `json:"name"`
  DesiredJobType     string `json:"desiredJobType"`
  DesiredLocation    string `json:"desiredLocation"`
  DesiredCompanySize string `json:"desiredCompanySize"`
  CareerAxis1        string `json:"careerAxis1"`
  CareerAxis2        string `json:"careerAxis2"`
  SelfPr             string `json:"selfPr"`
}

func EditUserProfile(db *gorm.DB) echo.HandlerFunc{
    return func(c echo.Context) error{
        // ユーザーIDを取得
      idStr := c.Param("id")
      id, err := strconv.Atoi(idStr)
      if err != nil {
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid user ID"})
      }
      fmt.Println("id:", id)

      // リクエストボディをバインド
    var body UserProfileUpdate
    if err := c.Bind(&body); err != nil {
      return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid request body"})
    }
    fmt.Println("body:", body)
  
    user := models.User{
      Name:               &body.Name,
      DesiredJobType:     &body.DesiredJobType,
      DesiredLocation:    &body.DesiredLocation,
      DesiredCompanySize: &body.DesiredCompanySize,
      CareerAxis1:        &body.CareerAxis1,
      CareerAxis2:        &body.CareerAxis2,
      SelfPr:             &body.SelfPr,
    }
    fmt.Println("user:", user)

    if err := db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error; err != nil {
      return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "error": err.Error()})
    }

    return c.JSON(http.StatusOK, echo.Map{"success": true, "user": user})
    }
}
