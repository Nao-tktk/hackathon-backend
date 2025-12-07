package usecase

import "db/model"

type UserRepository interface {
	FindByName(name string) ([]model.User, error)
	Insert(user *model.User) error
}
