package infrastructure

import (
	"backend/internal/domain/model"
	"backend/internal/util"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type GenerateReasonRepositoryImpl struct {
	DB *gorm.DB
}

func NewGenerateReasonRepositoryImpl() *GenerateReasonRepositoryImpl {
	return &GenerateReasonRepositoryImpl{}
}

func (r *GenerateReasonRepositoryImpl) GenerateReasonFromAI(ctx context.Context, user *model.User, match []model.MatchItem, qa model.QuestionAnswers) (string, error) {
	prompt := buildCompanyReasonPrompt(user, match, qa)

	result, err := util.CallOpenAIWithPrompt(ctx, prompt)
	if err != nil {
		return "", err
	}
	return result, nil
}

func buildCompanyReasonPrompt(user *model.User, match []model.MatchItem, qa model.QuestionAnswers) string {
	var sb strings.Builder
	sb.WriteString("以下のユーザープロフィール、企業分析結果、本人の自由記述をもとに、その企業への志望理由を300文字以内で日本語で作成してください。\n\n")

	sb.WriteString("【ユーザー情報】\n")
	sb.WriteString("名前: " + util.Coalesce(user.Name) + "\n")
	sb.WriteString("志望職種: " + util.Coalesce(user.DesiredJobType) + "\n")
	sb.WriteString("志望勤務地: " + util.Coalesce(user.DesiredLocation) + "\n")
	sb.WriteString("志望企業の規模: " + util.Coalesce(user.DesiredCompanySize) + "\n")
	sb.WriteString("就活軸①: " + util.Coalesce(user.CareerAxis1) + "\n")
	sb.WriteString("就活軸②: " + util.Coalesce(user.CareerAxis2) + "\n")
	sb.WriteString("自己PR: " + util.Coalesce(user.SelfPr) + "\n\n")

	sb.WriteString("【企業とのマッチ分析結果】\n")
	for _, item := range match {
		sb.WriteString(fmt.Sprintf("・%s（スコア: %d）: %s\n", item.Axis, item.Score, item.Reason))
	}

	sb.WriteString("\n【本人記入の自由回答】\n")
	sb.WriteString("・この企業に興味を持った理由: " + qa.ReasonInterest + "\n")
	sb.WriteString("・魅力を感じた製品・サービス: " + qa.AttractiveService + "\n")
	sb.WriteString("・それに関連する自身の経験: " + qa.RelatedExperience + "\n\n")

	sb.WriteString("この情報をもとに、エンジニアとしてのキャリア視点から一貫性のある志望理由を300文字程度で生成してください。\n")
	return sb.String()
}
