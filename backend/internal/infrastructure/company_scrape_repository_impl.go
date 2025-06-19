package infrastructure

import (
	"backend/internal/domain/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

type CompanyScrapeRepositoryImpl struct {
	DB *gorm.DB
}

var (
	reImportant = regexp.MustCompile(`(?i)recruit|job|vision|engineer|business|services|technology|interview|message|culture|about|company|profile|corporate|access|location|history`)
	reTag       = regexp.MustCompile(`<[^>]*>`)
	reSpace     = regexp.MustCompile(`\s+`)
)

type pageInfo struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

func (r *CompanyScrapeRepositoryImpl) ScrapeCompany(ctx context.Context, companyURL string, user *model.User) ([]model.MatchItem, error) {
	// スクレイピング
	cdpCtx, cancel := chromedp.NewContext(ctx)
	defer cancel()

	var html string
	if err := chromedp.Run(cdpCtx,
		chromedp.Navigate(companyURL),
		chromedp.OuterHTML("html", &html),
	); err != nil {
		return nil, errors.New("scraping failed")
	}

	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	base, _ := url.Parse(companyURL)

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

	var targets []string
	for link := range linkSet {
		if reImportant.MatchString(link) {
			targets = append(targets, link)
		}
		if len(targets) == 5 {
			break
		}
	}

	var pages []pageInfo
	for _, link := range targets {
		var body string
		if err := chromedp.Run(cdpCtx,
			chromedp.Navigate(link),
			chromedp.InnerHTML("body", &body, chromedp.ByQuery),
		); err == nil {
			text := reSpace.ReplaceAllString(reTag.ReplaceAllString(body, ""), " ")
			if len(text) > 2000 {
				text = text[:2000]
			}
			pages = append(pages, pageInfo{URL: link, Text: text})
		}
	}

	var sb strings.Builder
	for _, p := range pages {
		sb.WriteString("▼ URL: " + p.URL + "\n" + p.Text + "\n\n")
	}
	formattedPages := sb.String()

	prompt := buildPrompt(user, formattedPages)

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
		return nil, errors.New("OpenAI request failed")
	}

	var result []model.MatchItem
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &result); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v", err)
	}

	return result, nil
}

func buildPrompt(user *model.User, content string) string {
	return strings.TrimSpace(`
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
` + content + `

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
}

func coalesce(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
