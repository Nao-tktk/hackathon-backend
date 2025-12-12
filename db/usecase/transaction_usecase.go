package usecase

import "fmt"

type TransactionUsecase struct {
	Repo TransactionRepository
}

func NewTransactionUsecase(repo TransactionRepository) *TransactionUsecase {
	return &TransactionUsecase{Repo: repo}
}

type PurchaseReq struct {
	ItemID  int `json:"item_id"`
	BuyerID int `json:"buyer_id"`
}

func (u *TransactionUsecase) Purchase(req PurchaseReq) error {
	if req.ItemID == 0 || req.BuyerID == 0 {
		return fmt.Errorf("invalid request")
	}
	return u.Repo.Purchase(req.ItemID, req.BuyerID)
}
