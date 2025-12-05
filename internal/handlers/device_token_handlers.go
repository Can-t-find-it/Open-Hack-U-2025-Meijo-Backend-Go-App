package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "hacku_2025_meijo/internal/database"
    "hacku_2025_meijo/internal/models"
)

type DeviceHandler struct{}

func NewDeviceHandler() *DeviceHandler {
    return &DeviceHandler{}
}

func (h *DeviceHandler) SaveDeviceToken(c *gin.Context) {
    var req struct {
        Token string `json:"token"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // JWT から userID を取得
    userID := c.GetString("userID") // ミドルウェアでセットされている前提
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "ユーザー認証に失敗しました"})
        return
    }

    // 既存トークンを上書きする
    database.DB.Where("user_id = ?", userID).Delete(&models.UserDeviceToken{})
    database.DB.Create(&models.UserDeviceToken{
        UserID: userID,
        Token:  req.Token,
    })

    c.Status(http.StatusNoContent)
}

