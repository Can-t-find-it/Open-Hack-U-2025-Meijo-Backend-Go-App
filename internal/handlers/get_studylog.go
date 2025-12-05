package handlers

import (
	"net/http"
	"hacku_2025_meijo/internal/service"
	"github.com/gin-gonic/gin"
)

func GetStudyLogsHandler(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JWTからユーザーIDの取得に失敗しました"})
		return
	}
	
	logs, err := service.GetLatestStudyLog(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}
