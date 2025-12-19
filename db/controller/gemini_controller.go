package controller

import (
	"context"
	"encoding/base64" // 👈 画像デコード用に必須
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/vertexai/genai"
)

const (
	GeminiProjectID = "term8-naoto-takaku"
	GeminiLocation  = "asia-northeast1"
	GeminiModel     = "gemini-2.5-flash"
)

type GeminiController struct{}

func NewGeminiController() *GeminiController {
	return &GeminiController{}
}

// フロントエンドから受け取るデータ
type GenerateReq struct {
	ItemName  string `json:"item_name"`
	ItemImage string `json:"item_image"`
}

// フロントエンドに返すデータ
type GenerateRes struct {
	Description string `json:"description"`
}

func (c *GeminiController) HandleGenerateDescription(w http.ResponseWriter, r *http.Request) {
	// CORS設定
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

	// 2. Geminiで文章を生成する（画像も渡す！）
	// ▼▼▼ ここを修正しました（引数を2つ渡す） ▼▼▼
	description, err := generateDescription(req.ItemName, req.ItemImage)

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
func generateDescription(itemName, itemImage string) (string, error) {
	ctx := context.Background()

	// クライアント作成
	client, err := genai.NewClient(ctx, GeminiProjectID, GeminiLocation)
	if err != nil {
		return "", fmt.Errorf("client creation failed: %w", err)
	}
	defer client.Close()

	// モデルを選択
	model := client.GenerativeModel(GeminiModel)
	model.SetTemperature(0.7)

	// ▼▼▼ AIへの入力データを作る（テキスト＋画像） ▼▼▼
	var inputs []genai.Part

	// 1. まずはテキスト（プロンプト）を入れる
	prompt := fmt.Sprintf("フリマアプリで「%s」を出品します。購買意欲をそそる魅力的な商品説明文を、200文字以内の日本語で作成してください。挨拶は不要で、いきなり本文から始めてください。", itemName)
	inputs = append(inputs, genai.Text(prompt))

	// 2. 画像がある場合は、デコードして追加する
	if itemImage != "" {
		// "data:image/jpeg;base64,......" から "......" の部分だけを取り出す
		parts := strings.Split(itemImage, ",")
		if len(parts) == 2 {
			// Base64文字列をバイト列に変換
			decodedData, err := base64.StdEncoding.DecodeString(parts[1])
			if err == nil {
				// 成功したら画像データとしてリストに追加
				// ※拡張子は便宜上 jpeg にしていますが、pngでもGeminiは読んでくれます
				inputs = append(inputs, genai.ImageData("jpeg", decodedData))

				// 画像用の指示も追加しておく
				inputs = append(inputs, genai.Text("\nまた、添付した画像の特徴（色、状態、付属品など）も文章に反映してください。"))
			} else {
				fmt.Printf("Base64 Decode Error: %v\n", err)
			}
		}
	}

	// 生成実行（inputs... でまとめて渡す）
	resp, err := model.GenerateContent(ctx, inputs...)
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
