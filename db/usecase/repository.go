package usecase

import "db/model"

type UserRepository interface {
	FindByName(name string) ([]model.User, error)
	Insert(user *model.User) (int, error)
}

type ItemRepository interface {
	GetItems() ([]model.Item, error)
	Insert(item *model.Item) (int, error)
}

type TransactionRepository interface {
	// 購入処理 (エラーなしなら購入完了)
	Purchase(itemID int, buyerID int) error
}
