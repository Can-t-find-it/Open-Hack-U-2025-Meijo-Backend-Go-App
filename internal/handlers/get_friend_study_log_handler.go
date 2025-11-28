package handlers

import (
	"hacku_2025_meijo/internal/dtos"
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
	var input dtos.InputFriendsStudyLog

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSONに間違いがあります"})
		return
	}

	response, err := h.service.GetFriendStudyLog(input.FriendsID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, response)
}
