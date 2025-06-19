package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/repository"
	"context"
)

type GenerateReasonUsecase struct {
	Repo repository.GenerateReasonRepository
}

func (u *GenerateReasonUsecase) GenerateReason(ctx context.Context, user *model.User, match []model.MatchItem, qa model.QuestionAnswers) (string, error) {
	return u.Repo.GenerateReasonFromAI(ctx, user, match, qa)
}
