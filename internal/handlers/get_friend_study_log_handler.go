package handlers

import (
	"net/http"
	"hacku_2025_meijo/internal/service"
	"github.com/gin-gonic/gin"
)

type FriendStudyLogHandler struct {
	service service.GetFriendStudyLogService
}

func NewFriendStudyLogHandler() *FriendStudyLogHandler {
	return &FriendStudyLogHandler{
		service: service.GetFriendStudyLogService{},
	}
}

func (h *FriendStudyLogHandler) GetFriendLog(c *gin.Context) {
	// 1. ユーザーIDを取得 (string)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// 2. Serviceを呼ぶ
	response, err := h.service.GetFriendStudyLog(userID)
	if err != nil {
		// ★修正: err だけだとJSONが空になることがあるので、err.Error() にする
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
