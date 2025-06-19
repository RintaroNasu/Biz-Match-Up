package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/repository"
	"context"
)

type CompanyScrapeUsecase struct {
	Repo repository.CompanyScrapeRepository
}

func (u *CompanyScrapeUsecase) ScrapeCompanyInfo(ctx context.Context, url string, user *model.User) ([]model.MatchItem, error) {
	return u.Repo.ScrapeCompany(ctx, url, user)
}
