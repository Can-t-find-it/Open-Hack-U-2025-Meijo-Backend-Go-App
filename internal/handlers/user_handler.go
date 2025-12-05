package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hacku_2025_meijo/internal/service"
)

// GET /api/user/status
func GetUserStatus(c *gin.Context) {
	// トークンからIDを取得
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Serviceを呼んで情報を取得
	response, err := service.GetUserStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}