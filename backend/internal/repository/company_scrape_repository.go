package repository

import (
	"backend/internal/domain/model"

	"context"
)

type CompanyScrapeRepository interface {
	ScrapeCompany(ctx context.Context, url string, user *model.User) ([]model.MatchItem, error)
}
