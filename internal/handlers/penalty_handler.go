package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "hacku_2025_meijo/internal/database"
    "hacku_2025_meijo/internal/models"
)

type PenaltyHandler struct{}

func NewPenaltyHandler() *PenaltyHandler {
    return &PenaltyHandler{}
}

func (h *PenaltyHandler) CheckPenalty(c *gin.Context) {

    // JWT または Context から取得するユーザーID
    postUserID := c.GetUint("userID")

    // フレンド取得
    var friends []models.Friend
    database.DB.Where("user_id = ?", postUserID).Find(&friends)

    if len(friends) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "inactive_friends": []string{},
            "message": "フレンドがいません",
        })
        return
    }

    // フレンドの ID を抽出
    friendIDs := []uint{}
    for _, f := range friends {
        friendIDs = append(friendIDs, f.FriendUserID)
    }

    // フレンドの「最後の学習時間」を取得
    inactiveFriends := []gin.H{}
    limitTime := time.Now().Add(-24 * time.Hour)

    for _, fid := range friendIDs {

        // 最新のログ1件を取得
        var log models.StudyLog
        logResult := database.DB.Where("user_id = ?", fid).
            Order("answered_at DESC").
            First(&log)

        // フレンドの情報を取得
        var u models.User
        userResult := database.DB.First(&u, fid)

        if userResult.Error != nil {
            if userResult.Error == gorm.ErrRecordNotFound {
                // 存在しないユーザーはスキップ
                continue
            } else {
                // DBエラーの場合は返す
                c.JSON(http.StatusInternalServerError, gin.H{"error": userResult.Error.Error()})
                return
            }
        }

        // ログが無い → サボり扱い
        if logResult.RowsAffected == 0 || log.AnsweredAt.Before(limitTime) {
            inactiveFriends = append(inactiveFriends, gin.H{
                "id":   u.ID,
                "name": u.Name,
            })
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "inactive_friends": inactiveFriends,
        "count":            len(inactiveFriends),
    })
}
