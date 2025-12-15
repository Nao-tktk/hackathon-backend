package model

import "time"

type Message struct {
	ID         int       `json:"id"`
	ItemID     int       `json:"item_id"` // 👈 追加: どの商品のチャットか
	SenderID   int       `json:"sender_id"`
	ReceiverID int       `json:"receiver_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
