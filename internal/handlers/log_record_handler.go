package handlers

import (
	"fmt"
	"net/http"

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
	// 1. ユーザーIDの取得 (ミドルウェアから)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインしてください"})
		return
	}

	// 2. パスパラメーター（教科書ID）の取得
	textbookIDStr := c.Param("id")
	
	//  教科書IDが空でないかチェック
	if textbookIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "教科書IDが指定されていません"})
		return
	}
	
	// 3. リクエストボディのバインド
	var input struct {
		Score float64 `json:"score" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です: " + err.Error()})
		return
	}
	
	// 4. サービスロジックの呼び出し
	fmt.Printf("UserID: %s, TextbookID: %s, Score: %f\n", userID, textbookIDStr, input.Score)

	err := h.service.RecordStudyLog(userID, textbookIDStr, input.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "スタディログの記録に失敗しました: " + err.Error()})
		return
	}

	// 5. 成功レスポンス
	c.Status(http.StatusNoContent)
}