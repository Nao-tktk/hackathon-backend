package usecase

import (
	"db/dao"
	"db/model"
	"fmt"
)

type MessageUsecase struct {
	Dao *dao.MessageDao
}

func NewMessageUsecase(d *dao.MessageDao) *MessageUsecase {
	return &MessageUsecase{Dao: d}
}

type SendMessageReq struct {
	ItemID     int    `json:"item_id"` // 👈 追加
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
}

// SendMessage: メッセージ送信
func (u *MessageUsecase) SendMessage(req SendMessageReq) error {
	if req.Content == "" {
		return fmt.Errorf("message content is empty")
	}
	msg := &model.Message{
		ItemID:     req.ItemID, // 👈 追加
		SenderID:   req.SenderID,
		ReceiverID: req.ReceiverID,
		Content:    req.Content,
	}
	return u.Dao.Create(msg)
}

// GetHistory: 履歴取得 (引数に itemID を追加)
func (u *MessageUsecase) GetHistory(itemID, user1, user2 int) ([]model.Message, error) {
	return u.Dao.GetConversation(itemID, user1, user2)
}
