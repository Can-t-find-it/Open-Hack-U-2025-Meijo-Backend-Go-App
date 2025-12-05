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
	"errors"

	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos" // DTOsパッケージをインポート
	"hacku_2025_meijo/internal/models"

	"github.com/google/uuid" // UUID生成用
	"gorm.io/gorm"

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
	あなたは厳格な問題作成マシーンです。
	提供された「ターゲット単語リスト」の各単語に対して、【必ず1単語につき1問ずつ】問題を作成してください。

	ターゲット単語リスト: %s
	作成パターン: %s

	【絶対守るべきルール】
	1. 出力する問題数は、入力された単語数と完全に一致させること。
	2. 出力の順番は、入力リストの順番と一致させること。
	3. 1つの単語に対して複数の問題を作らないこと。

	【出力フォーマット】
	以下の形式（区切り文字 "（）"）のみを出力してください。挨拶は不要です。

	---
	問題文: (1つ目の単語の問題)
	解説: (1つ目の単語の解説)
	---
	問題文: (2つ目の単語の問題)
	解説: (2つ目の単語の解説)
	---
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

	ターゲット単語: %s
	作成パターン: %s
	既存の問題文: %s

	【重要：作成パターンの定義】
	- "4択問題形式" の場合: スタンダードなクイズを作成してください。（例：「～は何？」）
	- "穴埋め4択" の場合: 問題文の中に空欄（）を作り、そこに入る言葉を答えさせてください。（例：「～は（）である」）
	出力フォーマット（必ず守ってください）:

	【条件】
	- 選択肢は必ず4つ（正解1、誤答3）作成すること。
	- 出力フォーマットを厳守すること。
	
	【出力フォーマット】
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

// SaveQuestionToDB: 生成された問題を保存する（重複チェック付き）
func SaveQuestionToDB(textbookID string, item dtos.ResultItem, answer string) (string, error) {
	var question models.Question

	// "textbook_id" と "answer" が一致するものを探す
	err := database.DB.Where("textbook_id = ? AND answer = ?", textbookID, answer).First(&question).Error

	if err != nil {
		// データが見つからない場合 (ErrRecordNotFound) は新規作成
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newQuestion := models.Question{
				ID:         uuid.New().String(), // ★修正: ここで確実にUUIDを生成して入れる！
				TextbookID: textbookID,
				Answer:     answer,
			}
			if createErr := database.DB.Create(&newQuestion).Error; createErr != nil {
				return "", createErr
			}
			// 作成したデータを question 変数に入れる
			question = newQuestion
		} else {
			// それ以外のDBエラー
			return "", err
		}
	}

	// 2. 重複チェック（同じ問題文が既にないか）
	var count int64
	database.DB.Model(&models.QuestionStatement{}).
		Where("question_id = ? AND statement = ?", question.ID, item.Question).
		Count(&count)

	if count > 0 {
		return question.ID, nil
	}

	// 3. 子データ（QuestionStatement）を追加する
	newStatement := models.QuestionStatement{
		ID:         uuid.New().String(), // ★修正: ここも確実にUUIDを入れる！
		QuestionID: question.ID,
		Statement:  item.Question,
		Explain:    item.Explanation,
		Choices:    item.Options,
	}

	if err := database.DB.Create(&newStatement).Error; err != nil {
		return "", err
	}

	return question.ID, nil
}

// DeleteQuestionByID: モデルを指定して削除
func DeleteQuestionByID(id string) error {
	if err := database.DB.Delete(&models.Question{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// 修正版: 教科書のタイプ（4択など）に合わせて、新しい切り口の問題を追加する
func GenerateAndAddStatement(questionID string) (*models.QuestionStatement, error) {
	// 1. 親問題(Question)を取得（Textbookの情報も一緒に！）
	var parentQuestion models.Question

	// Preload("Textbook") を追加して、教科書のタイプ（4択など）を知れるようにする
	if err := database.DB.Preload("Textbook").Preload("QuestionStatements").First(&parentQuestion, "id = ?", questionID).Error; err != nil {
		return nil, err
	}

	// 2. 既存の聞き方リストを作る（AIに「これとは違う聞き方にして」と伝えるため）
	var existingTexts []string
	for _, s := range parentQuestion.QuestionStatements {
		existingTexts = append(existingTexts, s.Statement)
	}

	// 3. 教科書のタイプをパターンとして使用
	// 例: textbook.Type が "4択問題形式" なら、それがそのままAIへの指示になる
	pattern := string(parentQuestion.Textbook.Type)

	// 4. AIに生成を依頼
	// existingTexts を渡すことで、AIは「これらと被らない、違う方向性の問題」を作ろうとします
	resultItem, err := GenerateSingle4ChoiceQuestion(parentQuestion.Answer, pattern, existingTexts)
	if err != nil {
		return nil, err
	}

	// 5. 保存
	newStatement := models.QuestionStatement{
		QuestionID: questionID,
		Statement:  resultItem.Question,
		Explain:    resultItem.Explanation,
		Choices:    resultItem.Options,
	}

	if err := database.DB.Create(&newStatement).Error; err != nil {
		return nil, err
	}

	return &newStatement, nil
}

// CreateFolder: 新しいフォルダを作成する
func CreateFolder(userID string, name string) (*models.Folder, error) {
	newFolder := models.Folder{
		UserID:   userID,
		Name:     name,
		Progress: 0, // 最初は0%
	}

	if err := database.DB.Create(&newFolder).Error; err != nil {
		return nil, err
	}

	return &newFolder, nil

}

// SuggestNewWordViaAI: 既存の単語リストを元に、AIに新しい単語を提案させる
func SuggestNewWordsViaAI(currentWords []string) ([]string, error) {
	// 単語リストを文字列にする（例: "バイナリ, クラウド, サーバー"）
	wordsStr := strings.Join(currentWords, ", ")

	prompt := fmt.Sprintf(`
You are an learning curriculum creator.
A student has already learned the following terms:
Learned terms: [%s]

Please propose five new terms that the student has *not learned yet* and that are highly related, helping them take the next step in their learning.

[ABSOLUTE RESTRICTIONS]
- Do NOT output any term that appears in the list of "Learned terms".
- Do NOT repeat the same term.

[OUTPUT FORMAT]
Output only the three terms separated by commas, like: "term1, term2, term3". No explanations.

[ADDITIONAL REQUIREMENT]
Provide your output in Japanese.
	`, wordsStr)

	apiResp, err := callOpenAIAPI(prompt)
	if err != nil {
		return nil, err
	}

	rawContent := strings.TrimSpace(apiResp.Choices[0].Message.Content)
	rawContent = strings.ReplaceAll(rawContent, "、", ",")
	rawList := strings.Split(rawContent, ",")

	// 重複フィルター処理
	var resultList []string

	// 1. 既存単語をマップに登録（検索を速くするため）
	existingMap := make(map[string]bool)
	for _, w := range currentWords {
		existingMap[strings.TrimSpace(w)] = true
	}

	// 2. AIの提案をチェックして、知らない単語だけ追加する
	for _, w := range rawList {
		cleanWord := strings.TrimSpace(w)
		// 空文字でなく、かつ既存リストに含まれていない場合のみ追加
		if cleanWord != "" && !existingMap[cleanWord] {
			resultList = append(resultList, cleanWord)
		}
	}

	// 結果が0個なら空リストを返す（無理やりDBから取らない）
	if len(resultList) == 0 {
		fmt.Println("AIからの提案がすべて重複していたため、候補なし(0件)を返します。")
		return []string{}, nil
	}

	return resultList, nil
}

// GenerateAndSaveBatch: 単語リストから問題を一括生成し、指定の教科書に保存する（共通機能）
func GenerateAndSaveBatch(textbookID string, textbookType string, answers []string) ([]dtos.ResultItem, error) {
	var resultItems []dtos.ResultItem
	var err error

	// 1. AI生成
	switch models.TextbookType(textbookType) {
	case models.Type4Choice, models.TypeFillIn4Choice:
		resultItems, err = Generate4ChoiceWorkbookForQAndA(answers, textbookType)
	case models.TypeFillIn, models.TypeInput:
		resultItems, err = GenerateWorkbookForQAndA(answers, textbookType)
	default:
		resultItems, err = GenerateWorkbookForQAndA(answers, textbookType)
	}

	if err != nil {
		return nil, err
	}

	// 2. DB保存 & ID埋め込み
	var finalQuestions []dtos.ResultItem
	for i, item := range resultItems {
		currentAnswer := ""
		if i < len(answers) {
			currentAnswer = answers[i]
		}

		id, err := SaveQuestionToDB(textbookID, item, currentAnswer)
		if err != nil {
			fmt.Println("Save Error:", err)
			continue
		}

		item.ID = id
		finalQuestions = append(finalQuestions, item)
	}

	return finalQuestions, nil
}
