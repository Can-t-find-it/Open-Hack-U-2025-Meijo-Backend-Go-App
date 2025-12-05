package handlers

import (
	"fmt"
	"net/http"

	"hacku_2025_meijo/internal/models"
	"hacku_2025_meijo/internal/service"

	"github.com/gin-gonic/gin"
)

// GET /api/user/status
func GetUserStatus(c *gin.Context) {
	// トークンからIDを取得
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Serviceを呼んで情報を取得
	response, err := service.GetUserStatus(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func UserSearchHandler(c *gin.Context) {
    // 1. 認証ユーザーIDの取得 (ログインユーザー)
	currentUserID := c.GetString("userID")
	if currentUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証ユーザー情報が見つかりません"})
		return
	}

    // 2. クエリパラメータ 'name' から検索キーワードを取得
	searchName := c.Query("name")
    
    // 3. サービス層を呼び出し
    // サービス層の構造体メソッドとして定義されている場合は、h.service.SearchAllUsers を使用してください。
	users, err := service.SearchAllUsers(currentUserID, searchName) 

	fmt.Printf("SearchName: %s\n", searchName)
	fmt.Printf("err: %v\n", err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー検索に失敗しました: " + err.Error()})
		return
	}

	if len(users) == 0 {
        // ユーザーが見つからなかった場合、空のリストとメッセージを返して終了
        c.JSON(http.StatusOK, gin.H{
            "users":   []models.User{}, //usersはmodels.Userのリストを想定
        })
        return
    }

    // 4. 成功レスポンス (JSONオブジェクト {"users": [...]})
	c.JSON(http.StatusOK, gin.H{"users": users})
}