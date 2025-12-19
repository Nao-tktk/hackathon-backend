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

// ▼ URLに含まれる情報
const (
	ProjectID = "term8-naoto-takaku"                  // プロジェクトID
	Location  = "global"                              // ロケーション
	EngineID  = "hackathon-manual-help_1766104642390" // エンジンID
)

const apiEndpoint = "https://discoveryengine.googleapis.com/v1beta/projects/%s/locations/%s/collections/default_collection/engines/%s/servingConfigs/default_search:search"

type HelpController struct{}

func NewHelpController() *HelpController {
	return &HelpController{}
}

// ▼▼▼ 変更点1: 構造体に ModelPromptSpec を追加 ▼▼▼
type SearchRequest struct {
	Query             string            `json:"query"`
	PageSize          int               `json:"pageSize"`
	ContentSearchSpec ContentSearchSpec `json:"contentSearchSpec"`
}
type ContentSearchSpec struct {
	SummarySpec SummarySpec `json:"summarySpec"`
}
type SummarySpec struct {
	SummaryResultCount int             `json:"summaryResultCount"`
	IncludeCitations   bool            `json:"includeCitations"`
	ModelPromptSpec    ModelPromptSpec `json:"modelPromptSpec"` // 追加
}
type ModelPromptSpec struct {
	Preamble string `json:"preamble"` // 追加: AIへの指示文
}

// レスポンスの型定義
type SearchResponse struct {
	Summary SummaryResponse `json:"summary"`
}
type SummaryResponse struct {
	SummaryText string `json:"summaryText"`
}

type HelpReq struct {
	Query string `json:"query"`
}

func (c *HelpController) HandleHelp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req HelpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	answer, err := searchWithREST(req.Query)
	if err != nil {
		fmt.Printf("Vertex AI Error: %v\n", err)
		http.Error(w, "AI processing failed", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"answer": answer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func searchWithREST(queryText string) (string, error) {
	url := fmt.Sprintf(apiEndpoint, ProjectID, Location, EngineID)

	// ▼▼▼ 変更点2: Preamble（指示）を追加してAIを積極化させる ▼▼▼
	requestBody := SearchRequest{
		Query:    queryText,
		PageSize: 5,
		ContentSearchSpec: ContentSearchSpec{
			SummarySpec: SummarySpec{
				SummaryResultCount: 5,
				IncludeCitations:   false,
				ModelPromptSpec: ModelPromptSpec{
					// ここで「検索結果を使って日本語で答えて」と明示します
					Preamble: "あなたはフリマアプリの親切なガイドです。提供された検索結果に基づいて、ユーザーの質問に日本語で回答してください。",
				},
			},
		},
	}
	jsonData, _ := json.Marshal(requestBody)

	ctx := context.Background()
	client, _, err := transport.NewHTTPClient(ctx, option.WithScopes("https://www.googleapis.com/auth/cloud-platform"))
	if err != nil {
		return "", fmt.Errorf("client creation failed: %v", err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API Error: %s", body)
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return "", fmt.Errorf("JSON parse failed: %v", err)
	}

	if searchResp.Summary.SummaryText == "" {
		// 本当に空だった場合のフォールバック
		return "申し訳ありません、関連する情報が見つかりませんでした。", nil
	}

	return searchResp.Summary.SummaryText, nil
}
