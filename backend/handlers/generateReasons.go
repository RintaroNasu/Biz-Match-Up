package handlers

import (
	"backend/models"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type QuestionAnswers struct {
	ReasonInterest    string `json:"reasonInterest"`
	AttractiveService string `json:"attractiveService"`
	RelatedExperience string `json:"relatedExperience"`
}

type GenerateReasonsRequest struct {
	MatchResult []MatchItem     `json:"matchResult"`
	Questions   QuestionAnswers `json:"questions"`
}

func GenerateReasons(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		// JWTトークンの検証
		var req GenerateReasonsRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
		}

		fmt.Println("マッチ結果:", req.MatchResult)
		fmt.Println("質問への回答:", req.Questions)

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

		cli := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		prompt := strings.TrimSpace(`
			以下のユーザープロフィール、企業分析結果、本人の自由記述をもとに、その企業への志望理由を300文字以内で日本語で作成してください。

			【ユーザー情報】
			名前: ` + coalesce(user.Name) + `
			志望職種: ` + coalesce(user.DesiredJobType) + `
			志望勤務地: ` + coalesce(user.DesiredLocation) + `
			志望企業の規模: ` + coalesce(user.DesiredCompanySize) + `
			就活軸①: ` + coalesce(user.CareerAxis1) + `
			就活軸②: ` + coalesce(user.CareerAxis2) + `
			自己PR: ` + coalesce(user.SelfPr) + `

			【企業とのマッチ分析結果】
			` + formatMatchResult(req.MatchResult) + `

			【本人記入の自由回答】
			・この企業に興味を持った理由: ` + req.Questions.ReasonInterest + `
			・魅力を感じた製品・サービス: ` + req.Questions.AttractiveService + `
			・それに関連する自身の経験: ` + req.Questions.RelatedExperience + `

			この情報をもとに、エンジニアとしてのキャリア視点から一貫性のある志望理由を300文字程度で生成してください。
		`)
		resp, err := cli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{Role: "user", Content: prompt},
			},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "OpenAIエラー: " + err.Error()})
		}

		generatedReason := resp.Choices[0].Message.Content

		return c.JSON(http.StatusOK, echo.Map{
			"success": true,
			"reason":  generatedReason,
		})
	}

}

func formatMatchResult(items []MatchItem) string {
	var sb strings.Builder
	for _, item := range items {
		sb.WriteString(fmt.Sprintf("・%s（スコア: %d）: %s\n", item.Axis, item.Score, item.Reason))
	}
	return sb.String()
}
