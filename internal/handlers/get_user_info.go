package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "hacku_2025_meijo/internal/database"
    "hacku_2025_meijo/internal/models"
)

// --- ユーザーのステータスを返す ---
func GetUserStatus(c *gin.Context) {
    userID := c.GetUint("userID") // JWTから取得している想定

    var user models.User
    if err := database.DB.Preload("Folders.Textbooks").First(&user, userID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "ユーザーが見つかりません"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Textbookの数をカウント
    textbookCount := 0
    for _, folder := range user.Folders {
        textbookCount += len(folder.Textbooks)
    }

    // 友達の数
    var friendCount int64
    database.DB.Model(&models.Friend{}).Where("user_id = ?", userID).Count(&friendCount)

    // 継続日数（streakDays）を計算
    streakDays := calculateStreakDays(userID)

    c.JSON(http.StatusOK, gin.H{
        "id":   user.ID,
        "name": user.Name,
        "status": gin.H{
            "textbookCount": textbookCount,
            "streakDays":    streakDays,
            "friendCount":   friendCount,
        },
    })
}

// --- 学習の継続日数を計算するサンプル関数 ---
func calculateStreakDays(userID uint) int {
    var logs []models.StudyLog
    database.DB.Where("user_id = ?", userID).Order("answered_at DESC").Find(&logs)

    if len(logs) == 0 {
        return 0
    }

    streak := 0
    today := time.Now().Truncate(24 * time.Hour)

    for i, log := range logs {
        logDay := log.AnsweredAt.Truncate(24 * time.Hour)
        expectedDay := today.AddDate(0, 0, -i)
        if logDay.Equal(expectedDay) {
            streak++
        } else {
            break
        }
    }

    return streak
}
