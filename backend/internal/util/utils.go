package util

import (
	"context"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
)

func GenerateToken(id uint, email string) (string, error) {
	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}

func Coalesce(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func CallOpenAIWithPrompt(ctx context.Context, prompt string) (string, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	ctxTimeout, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctxTimeout, openai.ChatCompletionRequest{
		Model: openai.GPT4TurboPreview,
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func ExtractUserIDFromToken(c echo.Context) (uint, error) {
	const bearer = "Bearer "
	auth := c.Request().Header.Get("Authorization")

	if !strings.HasPrefix(auth, bearer) {
		return 0, errors.New("missing bearer token")
	}
	tokenStr := strings.TrimPrefix(auth, bearer)

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	uid, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("invalid claims")
	}
	return uint(uid), nil
}
