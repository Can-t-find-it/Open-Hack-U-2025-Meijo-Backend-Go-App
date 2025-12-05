package handlers

import (
	"net/http"

	"hacku_2025_meijo/internal/service"

	"github.com/gin-gonic/gin"
)

type GetTextBooksHandler struct {
	service service.GetTextbooksService
}

func NewGetTextBooksHandler() *GetTextBooksHandler {
	return &GetTextBooksHandler{
		service: service.GetTextbooksService{},
	}
}

func (h *GetTextBooksHandler) GetTextbooks(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインしてください"})
		return
	}

	response, err := h.service.GetTextbooks(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *GetTextBooksHandler) GetFriendTextbookDetail(c *gin.Context) {
	// 1. URLから教科書IDを取得
	textbookID := c.Param("textbook_id")

	// 2. 既存のService関数を再利用して詳細データを取得
	// (Textbookデータは誰のものでも構造は同じなので、共通の取得処理を使えます)
	result, err := service.GetTextbookDetail(textbookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Textbook not found"})
		return
	}

	// 3. 指定された形式 { "textbook": ... } で返す
	c.JSON(http.StatusOK, gin.H{
		"textbook": gin.H{
			"id":        result.ID,
			"name":      result.Name,
			"type":      result.Type,
			"questions": result.Questions,
			// Score と Times は含めない
		},
	})
}
