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
        UserID uint   `json:"user_id"`
        Token  string `json:"token"`
    }

    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 既存トークンを上書きする
    database.DB.Where("user_id = ?", req.UserID).Delete(&models.UserDeviceToken{})
    database.DB.Create(&models.UserDeviceToken{
        UserID: req.UserID,
        Token:  req.Token,
    })

    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
