package handlers

import (
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ChangeFriendsHandler struct {
	service service.ChangeFriendsService
}

func NewAddFriendsHandler() *ChangeFriendsHandler {
	return &ChangeFriendsHandler{
		service: service.ChangeFriendsService{},
	}
}

func (h *ChangeFriendsHandler) AddFriends(c *gin.Context) {
	userID := c.GetUint("userID")

	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "再度ログインしてください"})
		return
	}

	var input dtos.AddFriends

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSONが間違っています" + err.Error()})
		return
	}
	response, err := h.service.AddFriends(userID, input)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

// 複数人削除 (DELETE /api/friends)
func (h *ChangeFriendsHandler) DeleteFriendsBatch(c *gin.Context) {
	// 1. 自分のID (トークン)
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input dtos.DeleteFriends
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON形式が不正です"})
		return
	}

	// 3. Serviceを一括削除モードで呼ぶ
	err := h.service.DeleteFriendsBatch(userID, input.DeleteFriends)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "指定されたフレンドを削除しました"})
}

// フレンド一括取得
func (h *ChangeFriendsHandler) GetFriends(c *gin.Context) {
	userID := c.GetUint("userID")

	response, err := h.service.GetFriends(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, response)
}
