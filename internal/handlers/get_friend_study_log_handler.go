package handlers

import (
	"hacku_2025_meijo/internal/service"
	"net/http"

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

	userID := c.GetString("userID")

	response, err := h.service.GetFriendStudyLog(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, response)
}
