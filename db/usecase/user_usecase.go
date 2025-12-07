package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"db/model"
	"github.com/oklog/ulid/v2"
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
	Name string
	Age  int
}

func (u *UserUsecase) RegisterUser(req RegisterUserReq) (string, error) {
	if err := u.validateRegisterRequest(req); err != nil {
		return "", err
	}

	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	user := &model.User{ID: id, Name: req.Name, Age: req.Age}
	if err := u.Repo.Insert(user); err != nil {
		return "", err
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
	if req.Age < 20 || req.Age > 80 {
		return fmt.Errorf("invalid age: out of range")
	}
	return nil
}
