package repository

import "backend/internal/domain/model"

type AuthRepository interface {
    FindByEmail(email string) (*model.User, error)
    Create(user *model.User) error
}
