package handlers

import (
	// "encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
)

// CORS ミドルウェア
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

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
