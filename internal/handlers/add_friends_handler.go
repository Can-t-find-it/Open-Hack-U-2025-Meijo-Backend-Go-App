package handlers

import (
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
	"net/http"
	"strconv"

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

	idStr := c.Param("id")

	friendID, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDは数字で指定してください"})
		return
	}

	var FriendID_slice []int

	FriendID_slice = append(FriendID_slice, friendID)

	var Finally_FriendID dtos.AddFriends

	Finally_FriendID = dtos.AddFriends{
		Friends: FriendID_slice,
	}

	_, err = h.service.AddFriends(userID, Finally_FriendID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, response)
	c.Status(http.StatusNoContent)
}

// 複数人削除 (DELETE /api/friends)
func (h *ChangeFriendsHandler) DeleteFriendsBatch(c *gin.Context) {
	// 1. 自分のID (トークン)
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	idStr := c.Param("id")

	deleteId, err := strconv.Atoi(idStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "IDは数字で指定してください"})
		return
	}

	var deleteIdSlice []int

	deleteIdSlice = append(deleteIdSlice, deleteId)

	// 3. Serviceを一括削除モードで呼ぶ
	error := h.service.DeleteFriendsBatch(userID, deleteIdSlice)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	// c.JSON(http.StatusOK, gin.H{"message": "指定されたフレンドを削除しました"})
	c.Status(http.StatusNoContent)
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
