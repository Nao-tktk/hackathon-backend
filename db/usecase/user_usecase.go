package usecase

import (
	"fmt"

	"db/model"
)

type UserUsecase struct {
	Repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{Repo: repo}
}

func (u *UserUsecase) SearchUser(name string) ([]model.User, error) {
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}
	return u.Repo.FindByName(name)
}

type RegisterUserReq struct {
	Name     string
	Password string
}

func (u *UserUsecase) RegisterUser(req RegisterUserReq) (int, error) {
	if err := u.validateRegisterRequest(req); err != nil {
		return 0, err
	}

	user := &model.User{
		Name:     req.Name,
		Password: req.Password,
	}

	id, err := u.Repo.Insert(user)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *UserUsecase) validateRegisterRequest(req RegisterUserReq) error {
	if req.Name == "" {
		return fmt.Errorf("invalid name: empty")
	}
	if len(req.Name) > 50 {
		return fmt.Errorf("invalid name: too long")
	}
	if len(req.Password) < 4 {
		return fmt.Errorf("invalid password: password too short; password must have at least 4 characters")
	}
	return nil
}
