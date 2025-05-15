package handlers

import (
  "net/http"
  "os"
  "time"

  "backend/models"

  "github.com/golang-jwt/jwt/v5"
  "github.com/labstack/echo/v4"
  "golang.org/x/crypto/bcrypt"
  "gorm.io/gorm"
)

type registerReq struct {
  Email              string  `json:"email"`
  Password           string  `json:"password"`
  Name               *string `json:"name"`
  DesiredJobType     *string `json:"desiredJobType"`
  DesiredLocation    *string `json:"desiredLocation"`
  DesiredCompanySize *string `json:"desiredCompanySize"`
  CareerAxis1        *string `json:"careerAxis1"`
  CareerAxis2        *string `json:"careerAxis2"`
  SelfPr             *string `json:"selfPr"`
}

type loginReq struct {
  Email    string `json:"email"`
  Password string `json:"password"`
}

type authRes struct {
  Message string      `json:"message"`
  User    interface{} `json:"user"`
  Token   string      `json:"token"`
}

func SignUp(db *gorm.DB) echo.HandlerFunc {
  return func(c echo.Context) error {
    var req registerReq
    if err := c.Bind(&req); err != nil {
      return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
    }

    var existing models.User
    if err := db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
      return c.JSON(http.StatusConflict, echo.Map{"error": "already registered"})
    }

    hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    user := models.User{
      Email:              req.Email,
      Password:           string(hash),
      Name:               req.Name,
      DesiredJobType:     req.DesiredJobType,
      DesiredLocation:    req.DesiredLocation,
      DesiredCompanySize: req.DesiredCompanySize,
      CareerAxis1:        req.CareerAxis1,
      CareerAxis2:        req.CareerAxis2,
      SelfPr:             req.SelfPr,
    }

    if err := db.Create(&user).Error; err != nil {
      return c.JSON(http.StatusInternalServerError, echo.Map{"error": "db error"})
    }

    token, _ := generateToken(user.ID, user.Email)

    return c.JSON(http.StatusOK, authRes{
      Message: "サインアップ成功",
      User:    user,
      Token:   token,
    })
  }
}

func SignIn(db *gorm.DB) echo.HandlerFunc {
  return func(c echo.Context) error {
    var req loginReq
    if err := c.Bind(&req); err != nil {
      return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
    }

    var user models.User
    if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
      return c.JSON(http.StatusNotFound, echo.Map{"error": "ユーザーが見つかりません"})
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
      return c.JSON(http.StatusUnauthorized, echo.Map{"error": "パスワードが間違っています"})
    }

    token, _ := generateToken(user.ID, user.Email)

    return c.JSON(http.StatusOK, authRes{
      Message: "ログイン成功",
      User:    user,
      Token:   token,
    })
  }
}

func generateToken(id uint, email string) (string, error) {
  claims := jwt.MapClaims{
    "id":    id,
    "email": email,
    "exp":   time.Now().Add(time.Hour).Unix(),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
  secret := os.Getenv("JWT_SECRET")
  return token.SignedString([]byte(secret))
}
