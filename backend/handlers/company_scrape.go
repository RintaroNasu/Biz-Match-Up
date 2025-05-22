package handlers

import (
	"backend/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type companyReq struct {
	CompanyURL string `json:"companyUrl"`
}

type pageInfo struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}
type MatchItem struct {
	Axis   string `json:"axis"`
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

var reImportant = regexp.MustCompile(`(?i)recruit|job|vision|engineer|business|services|technology|interview|message|culture|about|company|profile|corporate|access|location|history`)
var reTag = regexp.MustCompile(`<[^>]*>`)
var reSpace = regexp.MustCompile(`\s+`)

func CompanyScrape(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		// JWTトークンの検証
		var req companyReq
		if err := c.Bind(&req); err != nil || req.CompanyURL == "" {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid json"})
		}

		const bearer = "Bearer "
		auth := c.Request().Header.Get("Authorization")
		if !strings.HasPrefix(auth, bearer) {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "missing bearer token"})
		}
		println("auth:", auth)

		tokenStr := strings.TrimPrefix(auth, bearer)
		println("tokenStr:", tokenStr)

		claims := jwt.MapClaims{}
		println("claims:", claims)

		secret := []byte(os.Getenv("JWT_SECRET"))
		println("secret:", secret)

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		println("tkn:", tkn)
		if err != nil || !tkn.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid token"})
		}

		uid, ok := claims["id"].(float64)
		println("uid:", uid)
		if !ok {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "invalid claims"})
		}

		// ユーザー取得
		var user models.User
		if err := db.First(&user, uint(uid)).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{"error": "user not found"})
		}
		fmt.Println("user:", user)

		// スクレイピング
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		fmt.Println("ctx:", ctx)

		var html string
		if err := chromedp.Run(ctx,
			chromedp.Navigate(req.CompanyURL),
			chromedp.OuterHTML("html", &html),
		); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "scrape failed"})
		}

		println("html:", html)

		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
		fmt.Println("doc:", doc)
		base, _ := url.Parse(req.CompanyURL)
		fmt.Println("base:", base)

		linkSet := map[string]struct{}{}
		doc.Find("header a").Each(func(_ int, s *goquery.Selection) {
			if href, ok := s.Attr("href"); ok && !strings.HasPrefix(href, "javascript") {
				u, _ := url.Parse(href)
				if !u.IsAbs() {
					href = base.ResolveReference(u).String()
				}
				linkSet[href] = struct{}{}
			}
		})
		fmt.Println("linkSet:", linkSet)

		// 重要URLを抽出
		var targets []string
		for link := range linkSet {
			if reImportant.MatchString(link) {
				targets = append(targets, link)
			}
			if len(targets) == 5 {
				break
			}
		}

		fmt.Println("targets:", targets)

		// 各ページの本文を取得
		var pages []pageInfo
		for _, link := range targets {
			var body string
			if err := chromedp.Run(ctx,
				chromedp.Navigate(link),
				chromedp.InnerHTML("body", &body, chromedp.ByQuery),
			); err == nil {
				text := reSpace.ReplaceAllString(
					reTag.ReplaceAllString(body, ""), " ",
				)
				if len(text) > 2000 {
					text = text[:2000]
				}
				pages = append(pages, pageInfo{URL: link, Text: text})
			}
		}
		fmt.Println("pages:", pages)

		// 取得したページの情報を整形
		var sb strings.Builder
		for _, p := range pages {
			sb.WriteString("▼ URL: " + p.URL + "\n" + p.Text + "\n\n")
		}
		formattedPages := sb.String()
		fmt.Println("formattedPages:", formattedPages)
		prompt := strings.TrimSpace(`
			以下のユーザープロフィールと企業のWebページ情報をもとに、次の評価軸でマッチ度（1〜5の数値）と理由をJSON形式で出力してください。
			
			【評価軸】
			1. 専門領域での開発
			2. 裁量権・自由度
			3. 勤務地
			4. 会社の規模
			
			【ユーザー情報】
			名前: ` + coalesce(user.Name) + `
			志望職種: ` + coalesce(user.DesiredJobType) + `
			志望勤務地: ` + coalesce(user.DesiredLocation) + `
			志望企業の規模: ` + coalesce(user.DesiredCompanySize) + `
			就活軸①: ` + coalesce(user.CareerAxis1) + `
			就活軸②: ` + coalesce(user.CareerAxis2) + `
			自己PR: ` + coalesce(user.SelfPr) + `
			
			【企業情報（Webページから抽出）】
			` + formattedPages + `
			
			以下のようなJSON形式で出力してください:
			コードブロックなどで囲まず、そのままJSON形式で出力してください。
			
			[
				{
					"axis": "専門領域での開発",
					"score": 4,
					"reason": "△△"
				},
				{
					"axis": "裁量権・自由度",
					"score": 2,
					"reason": "△△"
				},
				{
					"axis": "勤務地",
					"score": 1,
					"reason": "△△"
				},
				{
					"axis": "会社の規模",
					"score": 5,
					"reason": "△△"
				}
			]
		`)

		fmt.Println("prompt:", prompt)

		// OpenAI APIを使用してマッチ度を評価
		cli := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		ctx, cancel2 := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel2()

		resp, err := cli.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{Role: "user", Content: prompt},
			},
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "openai error"})
		}
		rawContent := resp.Choices[0].Message.Content

		var parsedResult []MatchItem
		if err := json.Unmarshal([]byte(rawContent), &parsedResult); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"success": false,
				"error":   "JSONパース失敗: " + err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"success":     true,
			"matchResult": parsedResult,
		})
	}
}

func coalesce(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
