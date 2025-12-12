package model

type Transaction struct {
	ID        int `json:"id"`
	ItemID    int `json:"item_id"`
	BuyerID   int `json:"buyer_id"`
	CreatedAt int `json:"created_at"` // 簡易的にintとしておく(表示に使わない仮定)
}
