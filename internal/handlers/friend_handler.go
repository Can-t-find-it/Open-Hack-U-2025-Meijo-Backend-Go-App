package handlers

import (
	"net/http"
	"hacku_2025_meijo/internal/service"
	"github.com/gin-gonic/gin"
)

// FriendHandler はフレンド関連のハンドラをまとめる構造体です。
type FriendHandler struct {
	// service.FriendService を依存性として保持
	service *service.FriendService 
}

// NewFriendHandler は FriendHandler のコンストラクタです。
func NewFriendHandler(s *service.FriendService) *FriendHandler {
    return &FriendHandler{
        service: s,
    }
}

// SearchFriendsHandler は、フレンド名を検索するエンドポイントです。
// GET /api/friends/search?name=キーワード
func (h *FriendHandler) SearchFriendsHandler(c *gin.Context) {
	// 1. 認証ユーザーIDの取得 (ミドルウェアから)
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証ユーザー情報が見つかりません"})
		return
	}

	// 2. クエリパラメータ 'name' から検索キーワードを取得
	searchName := c.Query("name")
	if searchName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "検索キーワード (name) が指定されていません"})
		return
	}

	// 3. サービス層を呼び出し
	// h.service 経由で SearchFriendsByName を呼び出す
	logs, err := h.service.SearchFriendsByName(userID, searchName) 

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "フレンド検索に失敗しました: " + err.Error()})
		return
	}

	// 4. 成功レスポンス
	// DTOのスライスを "friends" キーの下に返します
	c.JSON(http.StatusOK, gin.H{"friends": logs})
}