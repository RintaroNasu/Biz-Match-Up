package repository

import "backend/internal/domain/model"

type UserRepository interface {
	EditUserProfile(id int, user *model.User) error
	GetUserByID(id int) (*model.User, error)
}
