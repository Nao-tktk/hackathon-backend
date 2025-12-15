package controller

import (
	"db/model"
	"db/usecase"
	"encoding/json"
	"net/http"
	"strconv"
)

type MessageController struct {
	Usecase *usecase.MessageUsecase
}

func NewMessageController(u *usecase.MessageUsecase) *MessageController {
	return &MessageController{Usecase: u}
}

func (c *MessageController) HandleMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// POST: メッセージ送信
	if r.Method == http.MethodPost {
		var req usecase.SendMessageReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if err := c.Usecase.SendMessage(req); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
		return
	}

	// GET: 履歴取得 (/messages?item_id=10&user_id=1&partner_id=2)
	if r.Method == http.MethodGet {
		q := r.URL.Query()
		itemID, _ := strconv.Atoi(q.Get("item_id")) // 👈 追加
		userID, _ := strconv.Atoi(q.Get("user_id"))
		partnerID, _ := strconv.Atoi(q.Get("partner_id"))

		if itemID == 0 || userID == 0 || partnerID == 0 {
			http.Error(w, "item_id, user_id, and partner_id are required", http.StatusBadRequest)
			return
		}

		msgs, err := c.Usecase.GetHistory(itemID, userID, partnerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if msgs == nil {
			msgs = []model.Message{}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msgs)
		return
	}
}
func (c *MessageController) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// クエリパラメータ ?user_id=1 を取得
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	notifs, err := c.Usecase.GetNotifications(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if notifs == nil {
		notifs = []model.Notification{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifs)
}
