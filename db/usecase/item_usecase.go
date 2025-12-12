package usecase

import "db/model"

type ItemUsecase struct {
	Repo ItemRepository
}

func NewItemUsecase(repo ItemRepository) *ItemUsecase {
	return &ItemUsecase{Repo: repo}
}

func (u *ItemUsecase) GetItems() ([]model.Item, error) {
	return u.Repo.GetItems()
}

// 出品時のリクエストパラメータ
type CreateItemReq struct {
	SellerID    int    `json:"seller_id"`
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	ImageName   string `json:"image_name"`
}

func (u *ItemUsecase) CreateItem(req CreateItemReq) (int, error) {
	// ここで金額チェック(マイナスじゃないか等)を入れると良いです
	item := &model.Item{
		SellerID:    req.SellerID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		ImageName:   req.ImageName,
	}
	return u.Repo.Insert(item)
}
