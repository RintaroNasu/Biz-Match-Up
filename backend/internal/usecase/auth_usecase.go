package usecase

import (
	"backend/internal/domain/model"
	"backend/internal/repository"
	"backend/internal/util"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	Repo repository.AuthRepository
}

type RegisterRequest struct {
	Email              string
	Password           string
	Name               *string
	DesiredJobType     *string
	DesiredLocation    *string
	DesiredCompanySize *string
	CareerAxis1        *string
	CareerAxis2        *string
	SelfPr             *string
}

type LoginRequest struct {
	Email    string
	Password string
}

func (u *AuthUsecase) SignUp(req RegisterRequest) (*model.AuthResponse, error) {
	existingUser, err := u.Repo.FindByEmail(req.Email)
	if err == nil && existingUser != nil {
		return nil, err
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := model.User{
		Email:              req.Email,
		Password:           string(hash),
		Name:               req.Name,
		DesiredJobType:     req.DesiredJobType,
		DesiredLocation:    req.DesiredLocation,
		DesiredCompanySize: req.DesiredCompanySize,
		CareerAxis1:        req.CareerAxis1,
		CareerAxis2:        req.CareerAxis2,
		SelfPr:             req.SelfPr,
	}

	if err := u.Repo.Create(&user); err != nil {
		return nil, err
	}

	token, _ := util.GenerateToken(user.ID, user.Email)

	return &model.AuthResponse{
		Message: "サインアップ成功",
		User:    user,
		Token:   token,
	}, nil
}

func (u *AuthUsecase) SignIn(req LoginRequest) (*model.AuthResponse, error) {
	user, err := u.Repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, _ := util.GenerateToken(user.ID, user.Email)
	return &model.AuthResponse{
		Message: "ログイン成功",
		User:    *user,
		Token:   token,
	}, nil
}
