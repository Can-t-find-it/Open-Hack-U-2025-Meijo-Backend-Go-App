package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
)

type ChangeFriendsHandler struct {
	service service.ChangeFriendsService
}

func NewAddFriendsHandler() *ChangeFriendsHandler {
	return &ChangeFriendsHandler{
		service: service.ChangeFriendsService{},
	}
}

// AddFriends: フレンド追加 (JSONでリストを受け取る)
func (h *ChangeFriendsHandler) AddFriends(c *gin.Context) {
	// 1. 自分のID取得 (Middlewareでセットされた値)
	// modelsをstringにしたので、Middlewareもstringでセットしているはず
	userID := c.GetString("userID") 
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "再度ログインしてください"})
		return
	}

	friendID := c.Param("id")
	if friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "friend_id is required"})
		return
	}

	// 3. Service呼び出し
	input := dtos.AddFriends{
		Friends: []string{friendID},
	}

	_, err := h.service.AddFriends(userID, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteFriendsBatch: 複数人削除 (JSONでリストを受け取る)
func (h *ChangeFriendsHandler) DeleteFriendsBatch(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}


	friendID := c.Param("id")
	if friendID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "friend_id is required"})
		return
	}

	// 3. Serviceを呼ぶ
	// (既存のDeleteFriendsBatchサービスはリストを受け取るので、1つだけのリストを作って渡す)
	err := h.service.DeleteFriendsBatch(userID, []string{friendID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}


	c.Status(http.StatusNoContent)
}

// GetFriends: フレンド一括取得
func (h *ChangeFriendsHandler) GetFriends(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	response, err := h.service.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
