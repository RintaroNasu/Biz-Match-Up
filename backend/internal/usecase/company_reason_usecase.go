package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/repository"
)

type CompanyReasonUsecase struct {
	Repo repository.ReasonRepository
}

type CompanyReasonRequest struct {
	Content     string
	UserID      uint
	CompanyName string
	CompanyUrl  string
}

func (u *CompanyReasonUsecase) SaveCompanyReason(req CompanyReasonRequest) error {
	reason := model.Reason{
		Content:     req.Content,
		UserID:      req.UserID,
		CompanyName: req.CompanyName,
		CompanyUrl:  req.CompanyUrl,
	}
	return u.Repo.Save(&reason)
}

func (u *CompanyReasonUsecase) GetCompanyReasons(userID uint) ([]model.Reason, error) {
	reasons, err := u.Repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return reasons, nil
}
