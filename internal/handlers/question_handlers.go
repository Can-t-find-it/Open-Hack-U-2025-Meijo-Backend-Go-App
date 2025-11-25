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

// ---- POST 問題生成 ----
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

	var resultItems []dtos.ResultItem
	var err error

	// -----------------------
	// ① 四択
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
	} else { // -----------------------
		// ② 一問一答・穴埋め
		// -----------------------
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

	// --- DB保存 ---
	ids := []int{}
	for _, item := range resultItems {
		id, err := service.SaveQuestionToDB(item, body.Pattern)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB保存に失敗"})
			return
		}
		ids = append(ids, int(id))
	}

	// --- レスポンス ---
	c.JSON(http.StatusOK, gin.H{
		"message":   "問題生成＆保存成功",
		"ids":       ids,
		"questions": resultItems,
	})
}

// DELETE /api/question/:id
func DeleteQuestionHandler(c *gin.Context) {
    id := c.Param("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
        return
    }

    err := service.DeleteQuestionByID(id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "問題を削除しました"})
}

