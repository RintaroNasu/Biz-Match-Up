package infrastructure

import (
	"backend/internal/domain/model"

	"gorm.io/gorm"
)

type AuthRepositoryImpl struct {
    DB *gorm.DB
}

func (r *AuthRepositoryImpl) FindByEmail(email string) (*model.User, error) {
    var user model.User
    if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *AuthRepositoryImpl) Create(user *model.User) error {
    return r.DB.Create(user).Error
}
