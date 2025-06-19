package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/repository"
)

type UserUsecase struct {
	Repo repository.UserRepository
}

type UpdateUserProfileRequest struct {
	Name               string
	DesiredJobType     string
	DesiredLocation    string
	DesiredCompanySize string
	CareerAxis1        string
	CareerAxis2        string
	SelfPr             string
}

func (u *UserUsecase) EditUserProfile(id int, req UpdateUserProfileRequest) (*model.User, error) {
	user := &model.User{
		Name:               &req.Name,
		DesiredJobType:     &req.DesiredJobType,
		DesiredLocation:    &req.DesiredLocation,
		DesiredCompanySize: &req.DesiredCompanySize,
		CareerAxis1:        &req.CareerAxis1,
		CareerAxis2:        &req.CareerAxis2,
		SelfPr:             &req.SelfPr,
	}
	if err := u.Repo.EditUserProfile(id, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUsecase) GetUserProfile(id int) (*model.User, error) {
	return u.Repo.GetUserByID(id)
}
