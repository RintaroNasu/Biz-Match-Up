package infrastructure

import (
	"backend/internal/domain/model"

	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func (r *UserRepositoryImpl) EditUserProfile(id int, user *model.User) error {
	return r.DB.Model(&model.User{}).Where("id = ?", id).Updates(user).Error
}

func (r *UserRepositoryImpl) GetUserByID(id int) (*model.User, error) {
	var user model.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
