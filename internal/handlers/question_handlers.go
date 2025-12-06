package handlers

import (
	"net/http"

	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
	"hacku_2025_meijo/internal/service"

	"github.com/gin-gonic/gin"
)

// ==========================================
// 共通: CORS ミドルウェア
// ==========================================
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// ==========================================
// 1. AI生成機能
// ==========================================

// ---- GET 四択 ----
func GenerateQuestion4ChoiceHandler(c *gin.Context) {
	word := c.Query("word")

	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word is required"})
		return
	}

	result, err := service.GenerateSingle4ChoiceQuestion(word, "1問1答", []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"word":        word,
		"question":    result.Question,
		"options":     result.Options,
		"answer":      word,
		"explanation": result.Explanation,
	})
}

// ---- POST 問題生成・保存 (教科書の設定に従う版) ----
// ---- POST 問題生成・保存 (パスパラメータ版) ----
func GenerateProblemHandler(c *gin.Context) {
	// 1. URLから教科書IDを取得
	// ★修正: 数字への変換(strconv.Atoi)を削除し、文字列のまま受け取る
	textbookID := c.Param("textbook_id")

	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// パスになければBodyから取得
	if textbookID == "" {
		if body.TextbookID != "" { // ★修正: stringなので "" と比較
			textbookID = body.TextbookID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "textbookId is required"})
			return
		}
	}

	// 2. 教科書の情報を取得
	var textbook models.Textbook
	// ★修正: ID(string)をそのまま検索条件に使う
	if err := database.DB.First(&textbook, "id = ?", textbookID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Textbook not found"})
		return
	}
	pattern := string(textbook.Type)

	// 3. 単語リストの整理
	var targetAnswers []string
	if len(body.Answers) > 0 {
		targetAnswers = body.Answers
	} else if body.Answer != "" {
		targetAnswers = []string{body.Answer}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "words (or answer) is required"})
		return
	}

	// 4. 生成ロジック
	var resultItems []dtos.ResultItem
	var err error

	switch models.TextbookType(pattern) {
	case models.Type4Choice, models.TypeFillIn4Choice:
		resultItems, err = service.Generate4ChoiceWorkbookForQAndA(targetAnswers, pattern)
	case models.TypeFillIn, models.TypeInput:
		resultItems, err = service.GenerateWorkbookForQAndA(targetAnswers, pattern)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported textbook type: " + pattern})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Generation Error: " + err.Error()})
		return
	}

	// 5. DB保存
	var finalResponse []dtos.ResultItem
	for i, item := range resultItems {
		currentAnswer := ""
		if i < len(targetAnswers) {
			currentAnswer = targetAnswers[i]
		}

		// ★修正: textbookID (string) をそのまま渡す！これでエラーが消えます
		id, err := service.SaveQuestionToDB(textbookID, item, currentAnswer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB保存エラー: " + err.Error()})
			return
		}
		item.ID = id
		finalResponse = append(finalResponse, item)
	}

	c.Status(http.StatusNoContent)
}

// ---- POST 四択 複数 (個別API) ----
func Generate4ChoiceWorkbookForQAndAHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if len(body.Answers) == 0 || body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "`解答` list and `pattern` are required"})
		return
	}

	resultItems, err := service.Generate4ChoiceWorkbookForQAndA(body.Answers, body.Pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultItems)
}

// ---- POST 四択 単体 (個別API) ----
func GenerateQuestion4ChoiceAPIHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if body.Answer == "" || body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解答 and pattern are required"})
		return
	}

	result, err := service.GenerateSingle4ChoiceQuestion(body.Answer, body.Pattern, body.ExistingQuestions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, []dtos.ResultItem{result})
}

// ---- POST 一問一答・穴埋め 複数 (個別API) ----
func GenerateWorkbookForQAndAHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if len(body.Answers) == 0 || body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "`解答` list and `pattern` are required"})
		return
	}

	resultItems, err := service.GenerateWorkbookForQAndA(body.Answers, body.Pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultItems)
}

// ==========================================
// 2. 教科書・DB管理機能
// ==========================================

// GET /api/textbooks - 自分の問題集一覧取得
func GetTextbooksHandler(c *gin.Context) {
	userID := c.GetString("userID")

	// IDが空（取れなかった）場合のエラー処理
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーIDがありません（再ログインしてください）"})
		return
	}

	result, err := service.GetUserTextbooks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ★ここで JSON の形を決めています
	// { "folder": [ ... ] } という形になります
	c.JSON(http.StatusOK, gin.H{
		"folder": result,
	})
}

// POST /api/textbooks - 教科書作成
func CreateTextbookHandler(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		FolderID string `json:"folderId"`
	}
	userID := c.GetString("userID")

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateTextbook(userID, req.Name, req.Type, req.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// GET /api/textbook/:id - 教科書詳細取得
func GetTextbookDetailHandler(c *gin.Context) {
	textbookID := c.Param("id")
	result, err := service.GetTextbookDetail(textbookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Textbook not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"textbook": result,
	})
}

// DELETE /api/textbooks/:id - 教科書削除
func DeleteTextbookHandler(c *gin.Context) {
	textbookID := c.Param("id")

	err := service.DeleteTextbook(textbookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// POST /api/textbook_result/:id - 学習結果保存
func UpdateTextbookResultHandler(c *gin.Context) {
	textbookID := c.Param("id")

	var req struct {
		Score float64 `json:"score"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	err := service.UpdateTextbookStatus(textbookID, req.Score, c.GetString("userID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// POST /api/question - 生成した問題をDBに保存 (単体)
func AddQuestionHandler(c *gin.Context) {
	var req struct {
		TextbookID string          `json:"textbookId"`
		Answer     string          `json:"answer"`
		ResultItem dtos.ResultItem `json:"resultItem"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	err := service.AddQuestionToTextbook(req.TextbookID, req.ResultItem, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ==========================================
// 3. 詳細な編集機能 (削除・追加)
// ==========================================

// DELETE /api/question/:id - 問題（親）ごとの削除
func DeleteQuestionHandler(c *gin.Context) {

	questionID := c.Param("id")

	// Serviceを呼ぶ
	err := service.DeleteQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// POST /api/questionstatements - 既存の問題に新しい聞き方を追加
func AddQuestionStatementHandler(c *gin.Context) {
	var req struct {
		QuestionID string   `json:"questionId"`
		Statement  string   `json:"questionStatement"`
		Explain    string   `json:"explain"`
		Choices    []string `json:"choices"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	err := service.AddQuestionStatement(req.QuestionID, req.Statement, req.Explain, req.Choices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /api/questionstatements - 問題文一括削除
func DeleteQuestionStatementHandler(c *gin.Context) {
	statementID := c.Param("id")

	// Serviceを呼ぶ
	err := service.DeleteQuestionStatement(statementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// GET /api/word - 覚えたい単語提案 (AI版)
func SuggestWordHandler(c *gin.Context) {

	textbookID := c.Param("id")

	if textbookID == "" {
		textbookID = c.Query("id")
	}

	var suggestedWords []string // 配列に変更
	var err error

	if textbookID != "" {
		// ★AI提案ルート
		currentWords, err := service.GetWordsInTextbook(textbookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch textbook words"})
			return
		}

		if len(currentWords) == 0 {
			// 教科書が空なら、AI分析できないのでランダムに切り替え
			suggestedWords, err = service.GetRandomSuggestedWords(3)
		} else {
			// AIに提案させる
			suggestedWords, err = service.SuggestNewWordsViaAI(currentWords)
		}

	} else {
		// ★完全ランダムルート（教科書指定なし）
		suggestedWords, err = service.GetRandomSuggestedWords(3)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Suggestion failed: " + err.Error()})
		return
	}

	// 配列で返す { "suggestWord": ["A", "B", "C"] }
	c.JSON(http.StatusOK, gin.H{
		"suggestWord": suggestedWords,
	})
}

// POST /api/generate_statement
func GenerateAndAddStatementHandler(c *gin.Context) {
	// 1. URLからIDを取得 (パスパラメータ)
	questionID := c.Param("question_id")

	var req struct {
		Type string `json:"type"` // 例: "穴埋め入力", "4択問題形式"
	}

	// 2. 生成 & 保存
	// ★修正: 戻り値を使わないので newStmt ではなく _ で受け取る！
	_, err := service.GenerateAndAddStatement(questionID, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Generation failed: " + err.Error()})
		return
	}

	// 3. ステータスコード 204 (No Content) を返す
	c.Status(http.StatusNoContent)
}

// POST /api/folders - フォルダ作成
func CreateFolderHandler(c *gin.Context) {
	var req struct {
		FolderName string `json:"folderName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// ユーザーIDは uint で取得される(ミドルウェア次第だが)
	userID := c.GetString("userID")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーIDがありません"})
		return
	}
	// 必要に応じて string に変換
	// userID := strconv.Itoa(int(uidUint))
	// ※ service.CreateFolder が uint を求めているならそのままでOK

	_, err := service.CreateFolder(userID, req.FolderName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// DELETE /api/folders/:id - フォルダ削除
func DeleteFoldersBatchHandler(c *gin.Context) {
	var req struct {
		FolderIds []string `json:"folderIds"` // 配列で受け取る
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if len(req.FolderIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	// Serviceを呼ぶ
	err := service.DeleteFoldersBatch(req.FolderIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
