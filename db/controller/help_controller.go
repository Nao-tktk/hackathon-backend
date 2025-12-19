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

// ▼ ここをご自身のIDに書き換えてください
const (
	ProjectID = "term8-naoto-takaku"                  // プロジェクトID
	Location  = "global"                              // ロケーション
	EngineID  = "hackathon-manual-help_1766104642390" // エンジンID
)

// ※注意: カリキュラムでは default_config ですが、最近の汎用検索アプリは default_search の場合が多いです。
// もし 404 エラーが出る場合は、末尾の default_search を default_config に戻してみてください。
const apiEndpoint = "https://discoveryengine.googleapis.com/v1beta/projects/%s/locations/%s/collections/default_collection/engines/%s/servingConfigs/default_search:search"

type HelpController struct{}

func NewHelpController() *HelpController {
	return &HelpController{}
}

// ▼▼▼ ここから下はカリキュラムの構造体定義 (そのまま) ▼▼▼

type SearchRequest struct {
	Query               string              `json:"query"`
	PageSize            int                 `json:"pageSize"`
	ContentSearchSpec   ContentSearchSpec   `json:"contentSearchSpec"`
	QueryExpansionSpec  QueryExpansionSpec  `json:"queryExpansionSpec"`
	SpellCorrectionSpec SpellCorrectionSpec `json:"spellCorrectionSpec"`
}

type ContentSearchSpec struct {
	SnippetSpec SnippetSpec `json:"snippetSpec"`
	SummarySpec SummarySpec `json:"summarySpec"`
}

type SnippetSpec struct {
	ReturnSnippet bool `json:"returnSnippet"`
}

type SummarySpec struct {
	SummaryResultCount           int             `json:"summaryResultCount"`
	IncludeCitations             bool            `json:"includeCitations"`
	IgnoreAdversarialQuery       bool            `json:"ignoreAdversarialQuery"`
	IgnoreNonSummarySeekingQuery bool            `json:"ignoreNonSummarySeekingQuery"`
	ModelPromptSpec              ModelPromptSpec `json:"modelPromptSpec"`
	ModelSpec                    ModelSpec       `json:"modelSpec"`
}

type ModelPromptSpec struct {
	Preamble string `json:"preamble"`
}

type ModelSpec struct {
	Version string `json:"version"`
}

type QueryExpansionSpec struct {
	Condition string `json:"condition"`
}

type SpellCorrectionSpec struct {
	Mode string `json:"mode"`
}

// ▲▲▲ カリキュラムの定義ここまで ▲▲▲

// ▼▼▼ 追加: レスポンスを受け取るための構造体 (これが無いと回答を取り出せないので追加) ▼▼▼
type SearchResponse struct {
	Summary SummaryResponse `json:"summary"`
}
type SummaryResponse struct {
	SummaryText string `json:"summaryText"`
}

// フロントエンドからの入力を受け取る用
type HelpReq struct {
	Query string `json:"query"`
}

// ハンドラー関数（Webサーバー用）
func (c *HelpController) HandleHelp(w http.ResponseWriter, r *http.Request) {
	// CORS設定
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. フロントエンドから質問を受け取る
	var req HelpReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. カリキュラムのロジックで検索実行
	answer, err := searchSample(ProjectID, Location, EngineID, req.Query)
	if err != nil {
		fmt.Printf("Vertex AI Error: %v\n", err)
		http.Error(w, "AI processing failed", http.StatusInternalServerError)
		return
	}

	// 3. 結果をフロントエンドに返す
	response := map[string]string{"answer": answer}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// カリキュラムの searchSample を少し改造（文字列を返すように変更）
func searchSample(projectID, location, engineID, searchQuery string) (string, error) {
	url := fmt.Sprintf(apiEndpoint, projectID, location, engineID)

	requestBody := SearchRequest{
		Query:    searchQuery,
		PageSize: 10,
		ContentSearchSpec: ContentSearchSpec{
			SnippetSpec: SnippetSpec{
				ReturnSnippet: true,
			},
			SummarySpec: SummarySpec{
				SummaryResultCount:           5,
				IncludeCitations:             false,
				IgnoreAdversarialQuery:       true,
				IgnoreNonSummarySeekingQuery: true,
				ModelPromptSpec: ModelPromptSpec{
					// ★ここだけ変更: アプリ用の指示に変えています
					Preamble: "あなたはフリマアプリのガイドです。検索結果に基づいて、ユーザーの質問に日本語で回答してください。回答の確度が低かったとしても、なるべく「関連する情報が見つかりませんでした」という回答はしないでください。特に「手数料」に関する質問には、検索結果にかかわらず必ず「販売価格の10%です」と回答してください。手数料にかかわらない質問に対しては手数料の情報を回答に含めないでください。",
				},
				ModelSpec: ModelSpec{
					Version: "stable",
				},
			},
		},
		QueryExpansionSpec: QueryExpansionSpec{
			Condition: "AUTO",
		},
		SpellCorrectionSpec: SpellCorrectionSpec{
			Mode: "AUTO",
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	ctx := context.Background()
	// カリキュラム通り transport を使用
	client, _, err := transport.NewHTTPClient(ctx, option.WithScopes("https://www.googleapis.com/auth/cloud-platform"))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP client: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response: %v\n%s", resp.StatusCode, body)
	}

	// ★追加部分: JSONから回答だけを取り出す処理
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return "", fmt.Errorf("JSON parse failed: %v", err)
	}

	if searchResp.Summary.SummaryText == "" {
		return "申し訳ありません、関連する情報が見つかりませんでした。", nil
	}

	return searchResp.Summary.SummaryText, nil
}
