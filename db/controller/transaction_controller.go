package controller

import (
	"db/usecase"
	"encoding/json"
	"net/http"
)

type TransactionController struct {
	Usecase *usecase.TransactionUsecase
}

func NewTransactionController(u *usecase.TransactionUsecase) *TransactionController {
	return &TransactionController{Usecase: u}
}

func (c *TransactionController) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		var req usecase.PurchaseReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if err := c.Usecase.Purchase(req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// 成功したら空のJSONを返す（または {"message": "ok"} など）
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}
}
