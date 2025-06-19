package infrastructure

import (
	"backend/internal/domain/model"

	"gorm.io/gorm"
)

type CompanyReasonRepositoryImpl struct {
	DB *gorm.DB
}

func (r *CompanyReasonRepositoryImpl) Save(reason *model.Reason) error {
	return r.DB.Create(reason).Error
}

func (r *CompanyReasonRepositoryImpl) FindByUserID(userID uint) ([]model.Reason, error) {
	var reasons []model.Reason
	if err := r.DB.Where("user_id = ?", userID).Find(&reasons).Error; err != nil {
		return nil, err
	}
	return reasons, nil
}
