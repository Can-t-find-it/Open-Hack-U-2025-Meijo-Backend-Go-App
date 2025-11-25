package handlers

import (
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FriendGetAllHandler struct {
	service service.GetFriendService
}

func NewFriendGetAllHandler() *FriendGetAllHandler {
	return &FriendGetAllHandler{
		service: service.GetFriendService{},
	}
}

func (h *FriendGetAllHandler) GetAllFriends(c *gin.Context) {
	var input dtos.InputFriendGetAll

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です"})
		return
	}
	response, err := h.service.GetAllFriends(input.UserID)	

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, response)
}
