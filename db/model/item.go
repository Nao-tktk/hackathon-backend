package model

type Item struct {
	ID           int    `json:"id"`
	SellerID     int    `json:"seller_id"`
	CategoryID   int    `json:"category_id"`
	CategoryName string `json:"category_name"` // JOINして取得したカテゴリ名
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Description  string `json:"description"`
	ImageName    string `json:"image_name"`
	Status       string `json:"status"` // "ON_SALE" など
}
