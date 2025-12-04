package handlers

import (
	"net/http"
	"strconv"

	"hacku_2025_meijo/internal/service"

	"github.com/gin-gonic/gin"
)

type StudyLogHandler struct {
	service service.StudyLogService
}

// NewStudyLogHandler : コンストラクタ
func NewStudyLogHandler() *StudyLogHandler {
	return &StudyLogHandler{
		service: service.StudyLogService{},
	}
}

func (h *StudyLogHandler) CreateStudyLogHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインしてください"})
		return
	}

	textbookIDStr := c.Param("id")
	_, err := strconv.Atoi(textbookIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教科書IDが無効です(数字で指定してください)"})
		return
	}

	var input struct {
		Score float64 `json:"score" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です: " + err.Error()})
		return
	}

	err = h.service.RecordStudyLog(userID, textbookIDStr, input.Score)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "学習ログの記録に失敗しました: " + err.Error()})
		return
	}

	// 5. 【提供】 成功レスポンス
	c.Status(http.StatusNoContent)
}
