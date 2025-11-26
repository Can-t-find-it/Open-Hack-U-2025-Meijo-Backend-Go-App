package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
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

// ---- POST 一問一答・穴埋め 単体 ----
func GenerateProblemHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if body.Answer == "" || body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解答 and pattern are required"})
		return
	}

	resultItems, err := service.GenerateWorkbookForQAndA([]string{body.Answer}, body.Pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(resultItems) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "問題生成に失敗しました"})
		return
	}

	c.JSON(http.StatusOK, resultItems[0])
}

// ---- POST 一問一答・穴埋め 複数 ----
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

// ---- POST 四択 複数 ----
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

// ---- POST 四択（単体）----
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

// ==========================================
// 2. 教科書・DB管理機能
// ==========================================

// GET /api/textbooks - 教科書一覧取得
func GetTextbooksHandler(c *gin.Context) {
	// テスト用にユーザーIDを1に固定
	userID := uint(1)

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

// POST /api/question - 生成した問題をDBに保存
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

// GET /api/word - 覚えたい単語提案
func SuggestWordHandler(c *gin.Context) {
	// Service層から単語を取得
	word, err := service.GetSuggestedWord()
	
	if err != nil {
		// データがない場合などのエラー処理
		c.JSON(http.StatusNotFound, gin.H{"error": "No words found"})
		return
	}

	// シンプルに単語を返す、またはJSONで返す
	// 仕様書に形式が詳しく書いていないですが、とりあえずJSONで返します
	c.JSON(http.StatusOK, gin.H{
		"word": word,
	})
}