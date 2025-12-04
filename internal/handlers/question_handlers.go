package handlers

import (
	"net/http"
	"strconv"

	"hacku_2025_meijo/internal/dtos"
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

// ---- POST 問題生成・保存 (統合ハンドラ) ----
func GenerateProblemHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "`pattern` is required"})
		return
	}

	// 教科書IDがないと保存できないのでエラーにする
	if body.TextbookID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "textbookId is required"})
		return
	}

	var resultItems []dtos.ResultItem
	var err error

	// -----------------------
	// ① 生成ロジック
	// -----------------------
	if body.Pattern == "四択" {
		if body.Answer != "" {
			var item dtos.ResultItem
			item, err = service.GenerateSingle4ChoiceQuestion(body.Answer, body.Pattern, body.ExistingQuestions)
			resultItems = append(resultItems, item)
		} else if len(body.Answers) > 0 {
			resultItems, err = service.Generate4ChoiceWorkbookForQAndA(body.Answers, body.Pattern)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "`answer` or `answers` is required for 4choice"})
			return
		}
	} else {
		// ② 一問一答・穴埋め
		if body.Answer != "" {
			resultItems, err = service.GenerateWorkbookForQAndA([]string{body.Answer}, body.Pattern)
		} else if len(body.Answers) > 0 {
			resultItems, err = service.GenerateWorkbookForQAndA(body.Answers, body.Pattern)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "`answer` or `answers` is required"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(resultItems) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "問題生成に失敗しました"})
		return
	}

	// -----------------------
	// ③ DB保存 (りょうさんの階層構造へ)
	// -----------------------
	ids := []int{}
	if body.TextbookID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "textbookId is required for saving"})
		return
	}

	for i, item := range resultItems {
		// 正解の単語を特定（単体生成なら body.Answer、複数なら配列から）
		currentAnswer := ""
		if body.Answer != "" {
			currentAnswer = body.Answer
		} else if i < len(body.Answers) {
			currentAnswer = body.Answers[i]
		}

		// ★修正: 教科書IDと正解を渡して保存
		id, err := service.SaveQuestionToDB(body.TextbookID, item, currentAnswer)
		if err != nil {
			// エラー内容（err.Error()）を表示するようにしておくと原因がわかりやすいです
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB保存に失敗: " + err.Error()})
			return
		}
		ids = append(ids, int(id))
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "問題生成＆保存成功",
		"ids":       ids,
		"questions": resultItems,
	})
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

// GET /api/textbooks - 教科書一覧取得
func GetTextbooksHandler(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーIDがありません"})
		return
	}
	userID := userIDValue.(uint)

	result, err := service.GetUserTextbooks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"folder": result})
}

// POST /api/textbooks - 教科書作成
func CreateTextbookHandler(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		FolderID uint   `json:"folderId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateTextbook(req.Name, req.Type, req.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Textbook created"})
}

// GET /api/textbook/:id - 教科書詳細取得
func GetTextbookDetailHandler(c *gin.Context) {
	textbookID := c.Param("id")
	result, err := service.GetTextbookDetail(textbookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Textbook not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DELETE /api/textbooks/:id - 教科書削除
func DeleteTextbookHandler(c *gin.Context) {
	textbookID := c.Param("id")
	err := service.DeleteTextbook(textbookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Textbook deleted"})
}

// POST /api/textbook_result - 学習結果保存
func UpdateTextbookResultHandler(c *gin.Context) {
	var req struct {
		TextbookID uint    `json:"textbookId"`
		Score      float64 `json:"score"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	idStr := strconv.Itoa(int(req.TextbookID))
	err := service.UpdateTextbookStatus(idStr, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Score updated successfully"})
}

// POST /api/question - 生成した問題をDBに保存 (単体)
func AddQuestionHandler(c *gin.Context) {
	var req struct {
		TextbookID uint            `json:"textbookId"`
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

	c.JSON(http.StatusOK, gin.H{"message": "Question added successfully"})
}

// ==========================================
// 3. 詳細な編集機能 (削除・追加)
// ==========================================

// DELETE /api/question/:id - 問題（親）ごとの削除
func DeleteQuestionHandler(c *gin.Context) {
	questionID := c.Param("id")
	err := service.DeleteQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}

// POST /api/questionstatements - 既存の問題に新しい聞き方を追加
func AddQuestionStatementHandler(c *gin.Context) {
	var req struct {
		QuestionID uint     `json:"questionId"`
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
	c.JSON(http.StatusOK, gin.H{"message": "Statement added successfully"})
}

// DELETE /api/questionstatements/:id - 問題文（子）単体の削除
func DeleteQuestionStatementHandler(c *gin.Context) {
	statementID := c.Param("id")
	err := service.DeleteQuestionStatement(statementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Statement deleted"})
}

// GET /api/word - 覚えたい単語提案 (AI版)
func SuggestWordHandler(c *gin.Context) {
	// クエリパラメータから教科書IDを取得 (?textbook_id=5)
	textbookID := c.Query("textbook_id")

	var suggestedWord string
	var err error

	if textbookID != "" {
		// ★パターンA: 教科書IDがある場合 → その中身を分析してAIが提案

		// 1. 教科書の中の単語を取得
		currentWords, err := service.GetWordsInTextbook(textbookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch textbook words"})
			return
		}

		if len(currentWords) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Textbook is empty"})
			return
		}

		// 2. AIに提案させる
		suggestedWord, err = service.SuggestNewWordViaAI(currentWords)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "AI suggestion failed"})
			return
		}

	} else {
		// ★パターンB: 指定がない場合 → 既存のランダム取得（以前のロジック）
		suggestedWord, err = service.GetSuggestedWord()
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "No words found in DB"})
			return
		}
	}

	// 返す
	c.JSON(http.StatusOK, gin.H{
		"word": suggestedWord,
	})
}

// POST /api/generate_statement
func GenerateAndAddStatementHandler(c *gin.Context) {
	var req struct {
		QuestionID uint `json:"questionId"`
		// Pattern string `json:"pattern"`  <-- これはもう不要！削除！
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// 引数が QuestionID だけになりました
	newStmt, err := service.GenerateAndAddStatement(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Generation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "New variation generated successfully",
		"statement": newStmt,
	})
}

// POST /api/folders - フォルダ作成
func CreateFolderHandler(c *gin.Context) {
	// Middlewareでセットされた userID を取得
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーIDがありません"})
		return
	}
	userID := userIDValue.(uint)

	var req struct {
		// Name string `json:"name"`
		Name string `json:"folderName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	_, err := service.CreateFolder(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusCreated, gin.H{
	// 	"folder": gin.H{
	// 		"id":     folder.ID,
	// 		"name":   folder.Name,
	// 		"userId": folder.UserID,
	// 	},
	// })
	c.Status(http.StatusNoContent) //そーしろーの要望に合わせましたby.上野
}

// DELETE /api/folders/:id - フォルダ削除
func DeleteFolderHandler(c *gin.Context) {
	// JSON の受け皿
	var req struct {
		FolderIDs []string `json:"folderIds"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if len(req.FolderIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folderIds is required"})
		return
	}

	err := service.DeleteFolders(req.FolderIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{
	// 	"folderIds": req.FolderIDs,
	// })
	c.Status(http.StatusNoContent) // そーしろーの要望に合わせましたby.上野
}
