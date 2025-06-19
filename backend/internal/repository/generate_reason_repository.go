package repository

import (
	"backend/internal/domain/model"
	"context"
)

type GenerateReasonRepository interface {
	GenerateReasonFromAI(ctx context.Context, user *model.User, match []model.MatchItem, qa model.QuestionAnswers) (string, error)
}
