package repository

import "backend/internal/domain/model"

type ReasonRepository interface {
	Save(reason *model.Reason) error
	FindByUserID(userID uint) ([]model.Reason, error)
}
