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
