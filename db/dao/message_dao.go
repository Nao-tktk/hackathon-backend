package dao

import (
	"database/sql"
	"db/model"
)

type MessageDao struct {
	db *sql.DB
}

func NewMessageDao(db *sql.DB) *MessageDao {
	return &MessageDao{db: db}
}

// Create: メッセージを保存
func (dao *MessageDao) Create(msg *model.Message) error {
	// item_id を追加してINSERT
	query := "INSERT INTO messages (item_id, sender_id, receiver_id, content) VALUES (?, ?, ?, ?)"
	_, err := dao.db.Exec(query, msg.ItemID, msg.SenderID, msg.ReceiverID, msg.Content)
	return err
}

// GetConversation: 「特定の商品」についての「2人のユーザー」の会話を取得
func (dao *MessageDao) GetConversation(itemID, user1, user2 int) ([]model.Message, error) {
	// item_id が一致し、かつ (自分→相手 OR 相手→自分) のメッセージを取得
	query := `
        SELECT id, item_id, sender_id, receiver_id, content, created_at 
        FROM messages 
        WHERE item_id = ? 
          AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))
        ORDER BY created_at ASC`

	rows, err := dao.db.Query(query, itemID, user1, user2, user2, user1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.Message
	for rows.Next() {
		var m model.Message
		if err := rows.Scan(&m.ID, &m.ItemID, &m.SenderID, &m.ReceiverID, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (dao *MessageDao) GetNotifications(userID int) ([]model.Notification, error) {
	query := `
        SELECT 
            m.id, m.item_id, i.name, m.sender_id, u.name, m.content, m.created_at
        FROM messages m
        JOIN items i ON m.item_id = i.id
        JOIN users u ON m.sender_id = u.id
        WHERE m.receiver_id = ?
        ORDER BY m.created_at DESC
    `
	rows, err := dao.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []model.Notification
	for rows.Next() {
		var n model.Notification
		if err := rows.Scan(&n.ID, &n.ItemID, &n.ItemName, &n.SenderID, &n.SenderName, &n.Content, &n.CreatedAt); err != nil {
			return nil, err
		}
		notifs = append(notifs, n)
	}
	return notifs, nil
}
