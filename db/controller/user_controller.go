package controller

import (
	"db/usecase"
	"encoding/json"
	"net/http"
)

type UserController struct {
	Usecase *usecase.UserUsecase
}

func NewUserController(u *usecase.UserUsecase) *UserController {
	return &UserController{Usecase: u}
}

func (c *UserController) Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	switch r.Method {
	case http.MethodGet:
		name := r.URL.Query().Get("name")
		users, err := c.Usecase.SearchUser(name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var req usecase.RegisterUserReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		id, err := c.Usecase.RegisterUser(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]int{"id": id})
	}
}

func (c *UserController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	// 1. CORS設定 (フロントエンドからのアクセスを許可)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// プリフライトリクエスト(OPTIONS)への対応
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 2. POSTメソッド以外は拒否
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 3. リクエストボディ(JSON)を読み取る
	var req usecase.RegisterUserReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// 4. ユースケースを呼び出して登録処理を実行
	id, err := c.Usecase.RegisterUser(req)
	if err != nil {
		// エラー内容をそのまま返す（"invalid name" や "user already exists" など）
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 5. 成功したら登録されたIDを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (c *UserController) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// CORS設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req usecase.LoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	id, err := c.Usecase.Login(req)
	if err != nil {
		// ログイン失敗 (401 Unauthorized)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// 成功したらIDを返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (c *UserController) HandleSocialLogin(w http.ResponseWriter, r *http.Request) {
	// 1. CORS設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. リクエスト読み込み
	var req usecase.SocialLoginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// 3. Usecase呼び出し
	id, name, err := c.Usecase.SocialLogin(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. 結果返却
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   id,
		"name": name,
	})
}
