package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos" // DTOsパッケージをインポート
	"hacku_2025_meijo/internal/models"
)



// OpenAI APIエンドポイント
const openaiAPIURL = "https://api.openai.com/v1/chat/completions"

// APIキーを環境変数から取得
var openaiAPIKey = os.Getenv("OPENAI_API_KEY")

// --- 共通ヘルパー関数: OpenAI API呼び出し ---

// callOpenAIAPI はOpenAIのチャットAPIを呼び出すヘルパー関数
func callOpenAIAPI(prompt string) (*dtos.OpenAIChatCompletionResponse, error) {
	if openaiAPIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY 環境変数が設定されていません")
	}

	reqBody := dtos.OpenAIChatCompletionRequest{
		Model: "gpt-4o-mini",
		Messages: []dtos.Message{
			{Role: "system", Content: "あなたは教育用の問題作成アシスタントです。"},
			{Role: "user", Content: prompt},
		},
		Temperature: 0.0,
		MaxTokens:   400,
	}

	jsonReq, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("リクエストJSONのエンコードエラー: %w", err)
	}

	req, err := http.NewRequest("POST", openaiAPIURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, fmt.Errorf("リクエストの作成エラー: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+openaiAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API呼び出しエラー: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラーレスポンス (Status: %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var apiResp dtos.OpenAIChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("レスポンスJSONのデコードエラー: %w", err)
	}

	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("APIレスポンスにコンテンツが含まれていません")
	}

	return &apiResp, nil
}

// --- サービスメソッド: 単語リストからの問題集生成 (一問一答/穴埋め) ---

// GenerateWorkbookForQAndA は単語リストとパターンから一括で問題と解説を生成する
func GenerateWorkbookForQAndA(answers []string, pattern string) ([]dtos.ResultItem, error) {
	answerListStr, _ := json.Marshal(answers)

	prompt := fmt.Sprintf(`
	次の複数の単語に関する問題文，及び解説を作成してください．
	解答形式は以下のようにしてください．それ以外の文字は完全に必要ありません．
	パターンによって問題の種別を変えてください．パターンには”1問1答”,”穴埋め”の2種類が存在します．

	単語群：%s
	パターン:%s

	// ... (プロンプトのフォーマット指示の続き) ...

	`, string(answerListStr), pattern)

	apiResp, err := callOpenAIAPI(prompt)
	if err != nil {
		return nil, err
	}
	content := apiResp.Choices[0].Message.Content

	// --- パース処理 ---
	reCodeBlock := regexp.MustCompile("```\n?(.*?)\n?```")
	matches := reCodeBlock.FindStringSubmatch(content)
	qaContent := content
	if len(matches) > 1 {
		qaContent = matches[1] // コードブロックの中身を優先
	}

	returnArray := []dtos.ResultItem{}
	qaBlocks := strings.Split(qaContent, "---")

	reQuestion := regexp.MustCompile(`問題文[:：](.*)`)
	reExplanation := regexp.MustCompile(`解説[:：](.*)`)

	for _, block := range qaBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		question := ""
		explanation := ""

		if m := reQuestion.FindStringSubmatch(block); len(m) > 1 {
			question = strings.TrimSpace(m[1])
		}
		if m := reExplanation.FindStringSubmatch(block); len(m) > 1 {
			explanation = strings.TrimSpace(m[1])
		}

		if question != "" && explanation != "" {
			returnArray = append(returnArray, dtos.ResultItem{
				Question:    question,
				Explanation: explanation,
			})
		}
	}

	return returnArray, nil
}

// --- サービスメソッド: 四択問題の生成とパース ---

// generateQuestion4ChoicePrompt は四択問題のプロンプトを生成する共通ヘルパー
func generateQuestion4ChoicePrompt(answer string, pattern string, existingQuestions []string) string {
	// 既存の問題文を文字列化（プロンプトで参照させるため）
	existingQStr := strings.Join(existingQuestions, " / ")

	return fmt.Sprintf(`
	あなたは教育用の問題作成アシスタントです。
	与えられた単語を答えとする問題を1つ作成してください。

	条件:
	- 答えは必ず「%s」になること。
	- 選択肢は必ず4つ用意してください（1つは正解、3つは誤答）。
	- 出力フォーマットは厳守してください。
	- 既存の問題文とは異なる新しい問題文を生成してください。
	既存の問題文: %s

	問題形式の指定:
	- %s に従って問題を作成してください。
	- "1問1答" の場合: 通常の四択問題形式。
	- "穴埋め" の場合: 問題文の中に空欄（（ ））を入れて四択問題を作成。

	出力フォーマット（必ず守ってください）:
	問題: ...
	選択肢:
	A: ...
	B: ...
	C: ...
	D: ...
	正解: A
	解説: ...
	`, answer, existingQStr, pattern)
}

// parse4ChoiceOutput はモデルからの出力文字列を ResultItem 構造体にパースする
func parse4ChoiceOutput(content string) (dtos.ResultItem, error) {
	result := dtos.ResultItem{}
	
	lines := strings.Split(content, "\n")
	
	// 問題文と解説の抽出（最初の出現行を優先）
	reQuestion := regexp.MustCompile(`問題[:：]\s*(.*)`)
	reExplanation := regexp.MustCompile(`解説[:：]\s*(.*)`)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if result.Question == "" {
			if m := reQuestion.FindStringSubmatch(line); len(m) > 1 {
				result.Question = strings.TrimSpace(m[1])
			}
		}
		if result.Explanation == "" {
			if m := reExplanation.FindStringSubmatch(line); len(m) > 1 {
				result.Explanation = strings.TrimSpace(m[1])
			}
		}
	}
	
	// 選択肢の抽出 (A: ... B: ... C: ... D: ...)
	options := []string{}
	// ラベルと本文を抽出する正規表現 (マルチバイト文字対応)
	// (?:A|B|C|D|Ａ|Ｂ|Ｃ|Ｄ)[\)\.\）．:：\s]*([^A-DＡ-Ｄ\n]+)
	reOption := regexp.MustCompile(`(?:[A-DＡ-Ｄ])[)\.）．:：\s]*([^A-DＡ-Ｄ\n]+)`) 
	matches := reOption.FindAllStringSubmatch(content, -1)
	
	for _, match := range matches {
		if len(match) > 1 {
			options = append(options, strings.TrimSpace(match[1]))
		}
	}
	
	if len(options) < 4 {
		return dtos.ResultItem{}, fmt.Errorf("モデル出力から4つの選択肢を抽出できませんでした")
	}
	
	// 最大4つに限定
	result.Options = options[:4]
	
	// 正解はResultItemには含めず、問題と選択肢を返します (元のPythonコードのロジックに合わせる)
	// 正解はクライアント側で管理するか、別のフィールドで返す必要がありますが、
	// ここではResultItemの定義に合わせて問題と選択肢のみを返します。

	return result, nil
}


// GenerateSingle4ChoiceQuestion は単語一つから四択問題を生成する
func GenerateSingle4ChoiceQuestion(answer string, pattern string, existingQuestions []string) (dtos.ResultItem, error) {
	prompt := generateQuestion4ChoicePrompt(answer, pattern, existingQuestions)
	
	apiResp, err := callOpenAIAPI(prompt)
	if err != nil {
		return dtos.ResultItem{}, err
	}
	
	content := apiResp.Choices[0].Message.Content
	return parse4ChoiceOutput(content)
}

// Generate4ChoiceWorkbookForQAndA は単語リストから四択問題集を生成する
func Generate4ChoiceWorkbookForQAndA(answers []string, pattern string) ([]dtos.ResultItem, error) {
	results := []dtos.ResultItem{}
	
	// 単語ごとに問題を生成（Pythonコードのロジックを再現）
	for _, answer := range answers {
		// 既存の問題文リストは、ここでは空として扱う（単語単位で独立しているため）
		resultItem, err := GenerateSingle4ChoiceQuestion(answer, pattern, []string{})
		if err != nil {
			// エラー発生時はその問題をスキップまたはエラーを返す
			fmt.Printf("警告: 単語 '%s' の四択問題生成に失敗しました: %v\n", answer, err)
			continue
		}
		results = append(results, resultItem)
	}

	return results, nil
}

// SaveQuestionToDB: 生成された問題を、階層構造DBに保存する
func SaveQuestionToDB(textbookID uint, item dtos.ResultItem, answer string) (uint, error) {

	// 1. 親データ（Question）を作成
	question := models.Question{
		TextbookID: textbookID,
		Answer:     answer, // 正解の単語
		// 2. 子データ（QuestionStatement）を作成
		QuestionStatements: []models.QuestionStatement{
			{
				Statement: item.Question,    // 問題文
				Explain:   item.Explanation, // 解説
				Choices:   item.Options,     // 選択肢 (GORMがJSON化して保存)
			},
		},
	}

	// 3. DBに保存
	result := database.DB.Create(&question)
	return question.ID, result.Error
}

// DeleteQuestionByID: モデルを指定して削除
func DeleteQuestionByID(id string) error {
	if err := database.DB.Delete(&models.Question{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}