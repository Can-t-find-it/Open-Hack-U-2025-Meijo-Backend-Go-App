package handlers

import (
	"net/http"
	"hacku_2025_meijo/internal/service"
	"github.com/gin-gonic/gin"
)

func GetStudyLogsHandler(c *gin.Context) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザーIDがありません"})
		return
	}
	userID := userIDValue.(uint)

	logs, err := service.GetUserStudyLogs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
