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

type PostReasonsRequest struct {
	Content     string `json:"content"`
	CompanyName string `json:"companyName"`
	CompanyUrl  string `json:"companyUrl"`
}

func PostReasons(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// リクエストのバインド
		var req PostReasonsRequest
		if err := c.Bind(&req); err != nil || strings.TrimSpace(req.Content) == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
		}

		// JWTトークンの検証
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

		// ユーザー取得
		var user models.User
		if err := db.First(&user, uint(uid)).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}

		// 志望理由を保存
		reason := models.Reason{
			Content:     req.Content,
			UserID:      user.ID,
			CompanyName: req.CompanyName,
			CompanyUrl:  req.CompanyUrl,
		}

		if err := db.Create(&reason).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to save reason"})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"success": true,
			"message": "志望理由を保存しました",
		})
	}
}
