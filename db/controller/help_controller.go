package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/transport"
)

// ▼ URLに含まれる情報（教えていただいた内容を埋めてあります）
const (
	ProjectID = "term8-naoto-takaku"                  // プロジェクトID
	Location  = "global"                              // ロケーション
	EngineID  = "hackathon-manual-help_1766104642390" // エンジンID
)

// Vertex AI Agent Builder の検索エンドポイント
const apiEndpoint = "https://discoveryengine.googleapis.com/v1beta/projects/%s/locations/%s/collections/default_collection/engines/%s/servingConfigs/default_search:search"

type HelpController struct{}

func NewHelpController() *HelpController {
	return &HelpController{}
}

// リクエストの型定義（カリキュラム準拠）
type SearchRequest struct {
	Query             string            `json:"query"`
	PageSize          int               `json:"pageSize"`
	ContentSearchSpec ContentSearchSpec `json:"contentSearchSpec"`
}
type ContentSearchSpec struct {
	SummarySpec SummarySpec `json:"summarySpec"` // 要約（回答）を要求する設定
}
type SummarySpec struct {
	SummaryResultCount int  `json:"summaryResultCount"`
	IncludeCitations   bool `json:"includeCitations"`
}

// レスポンスの型定義
type SearchResponse struct {
	Summary SummaryResponse `json:"summary"`
}
type SummaryResponse struct {
	SummaryText string `json:"summaryText"` // ここにAIの回答が入る
}

// フロントエンドからのリクエスト受け取り用
type HelpReq struct {
	Query string `json:"query"`
}

// ハンドラー関数
func (c *HelpController) HandleHelp(w http.ResponseWriter, r *http.Request) {
	// CORS設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// フロントエンドから質問文を受け取る
	var req HelpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Vertex AI に問い合わせる
	answer, err := searchWithREST(req.Query)
	if err != nil {
		fmt.Printf("Vertex AI Error: %v\n", err)
		http.Error(w, "AI processing failed", http.StatusInternalServerError)
		return
	}

	// 結果を返す
	response := map[string]string{"answer": answer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// REST APIでVertex AIを叩く関数（カリキュラムの方式）
func searchWithREST(queryText string) (string, error) {
	// 1. URLの組み立て
	url := fmt.Sprintf(apiEndpoint, ProjectID, Location, EngineID)

	// 2. リクエストボディの作成
	requestBody := SearchRequest{
		Query:    queryText,
		PageSize: 5,
		ContentSearchSpec: ContentSearchSpec{
			SummarySpec: SummarySpec{
				SummaryResultCount: 5,
				IncludeCitations:   false,
			},
		},
	}
	jsonData, _ := json.Marshal(requestBody)

	// 3. HTTPクライアント作成（自動認証付き）
	ctx := context.Background()
	// カリキュラム通り、transportを使って認証ヘッダーを自動付与します
	client, _, err := transport.NewHTTPClient(ctx, option.WithScopes("https://www.googleapis.com/auth/cloud-platform"))
	if err != nil {
		return "", fmt.Errorf("client creation failed: %v", err)
	}

	// 4. リクエスト送信
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 5. レスポンス読み取り
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API Error: %s", body)
	}

	// 6. JSONパースして回答を取り出す
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return "", fmt.Errorf("JSON parse failed: %v", err)
	}

	if searchResp.Summary.SummaryText == "" {
		return "関連する情報が見つかりませんでした。", nil
	}

	return searchResp.Summary.SummaryText, nil
}
