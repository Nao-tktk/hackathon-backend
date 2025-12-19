package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/vertexai/genai"
)

const (
	GeminiProjectID = "term8-naoto-takaku" // あなたのプロジェクトID
	GeminiLocation  = "asia-northeast1"    // 日本リージョン (Tokyo)
	GeminiModel     = "gemini-2.5-flash"   // 高速・安価なモデル
)

type GeminiController struct{}

func NewGeminiController() *GeminiController {
	return &GeminiController{}
}

// フロントエンドから受け取るデータ
type GenerateReq struct {
	ItemName string `json:"item_name"`
}

// フロントエンドに返すデータ
type GenerateRes struct {
	Description string `json:"description"`
}

func (c *GeminiController) HandleGenerateDescription(w http.ResponseWriter, r *http.Request) {
	// CORS設定（フロントエンドからのアクセスを許可）
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. リクエストを受け取る
	var req GenerateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Geminiで文章を生成する
	description, err := generateDescription(req.ItemName)
	if err != nil {
		fmt.Printf("Gemini Error: %v\n", err)
		http.Error(w, "AI generation failed", http.StatusInternalServerError)
		return
	}

	// 3. 結果を返す
	res := GenerateRes{Description: description}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// 実際にGeminiを呼び出す関数
func generateDescription(itemName string) (string, error) {
	ctx := context.Background()

	// クライアント作成
	client, err := genai.NewClient(ctx, GeminiProjectID, GeminiLocation)
	if err != nil {
		return "", fmt.Errorf("client creation failed: %w", err)
	}
	defer client.Close()

	// モデルを選択
	model := client.GenerativeModel(GeminiModel)
	model.SetTemperature(0.7) // 創造性の度合い（程よく自由に）

	// プロンプト（命令文）の作成
	prompt := fmt.Sprintf("フリマアプリで「%s」を出品します。購買意欲をそそる魅力的な商品説明文を、200文字以内の日本語で作成してください。挨拶は不要で、いきなり本文から始めてください。", itemName)

	// 生成実行
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	// 結果の取り出し
	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if txt, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			return string(txt), nil
		}
	}

	return "説明文の生成に失敗しました。", nil
}
