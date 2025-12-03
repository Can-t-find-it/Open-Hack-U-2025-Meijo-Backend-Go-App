package handlers

import (
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "hacku_2025_meijo/internal/database"
    "hacku_2025_meijo/internal/models"

    "github.com/sideshow/apns2"
    "github.com/sideshow/apns2/token"
)

type PenaltyHandler struct {
    APNsClient *apns2.Client
}

// --- APNs 初期化 ---
func NewPenaltyHandler() *PenaltyHandler {
    keyPath := os.Getenv("APNS_AUTH_KEY_PATH")
    keyID := os.Getenv("APNS_KEY_ID")
    teamID := os.Getenv("APNS_TEAM_ID")

    if keyPath == "" || keyID == "" || teamID == "" {
        panic("APNs用環境変数が設定されていません")
    }

    authKey, err := token.AuthKeyFromFile(keyPath)
    if err != nil {
        panic("APNs AuthKey の読み込みに失敗: " + err.Error())
    }

    tok := &token.Token{
        AuthKey: authKey,
        KeyID:   keyID,
        TeamID:  teamID,
    }

    client := apns2.NewTokenClient(tok).Development() // ← 開発中は Development()

    return &PenaltyHandler{
        APNsClient: client,
    }
}

func (h *PenaltyHandler) CheckPenalty(c *gin.Context) {
    postUserID := c.GetUint("userID") // ← ログインユーザー

    var friends []models.Friend
    database.DB.Where("user_id = ?", postUserID).Find(&friends)

    if len(friends) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "inactive_friends": []string{},
            "message":          "フレンドがいません",
        })
        return
    }

    friendIDs := []uint{}
    for _, f := range friends {
        friendIDs = append(friendIDs, f.FriendUserID)
    }

    inactiveFriends := []gin.H{}
    limitTime := time.Now().Add(-24 * time.Hour)

    for _, fid := range friendIDs {
        var log models.StudyLog
        logResult := database.DB.Where("user_id = ?", fid).
            Order("answered_at DESC").
            First(&log)

        var u models.User
        userResult := database.DB.First(&u, fid)
        if userResult.Error != nil {
            if userResult.Error == gorm.ErrRecordNotFound {
                continue
            }
            c.JSON(http.StatusInternalServerError, gin.H{"error": userResult.Error.Error()})
            return
        }

        // --- 友達が24時間勉強していない場合 ---
        if logResult.RowsAffected == 0 || log.AnsweredAt.Before(limitTime) {

            // ★ 通知先はログインしているユーザー本人
            go h.sendPenaltyNotification(postUserID, u.Name)

            inactiveFriends = append(inactiveFriends, gin.H{
                "id":   u.ID,
                "name": u.Name,
            })
        }
    }

    if len(inactiveFriends) == 0 {
        c.JSON(http.StatusOK, gin.H{
            "inactive_friends": []string{},
            "count":            0,
            "message":          "全員24時間以内に勉強しています！",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "inactive_friends": inactiveFriends,
        "count":            len(inactiveFriends),
    })
}

// --- APNs 通知 ---
func (h *PenaltyHandler) sendPenaltyNotification(userID uint, name string) {

    var device models.UserDeviceToken
    result := database.DB.Where("user_id = ?", userID).Last(&device)

    if result.Error != nil || device.Token == "" {
        fmt.Println("APNs: トークンなし、通知スキップ")
        return
    }

    payload := fmt.Sprintf(`{
        "aps": {
            "alert": {
                "title": "勉強していない友だちがいます！",
                "body": "%sさんは24時間勉強していません。"
            },
            "sound": "default"
        }
    }`, name)

    notification := &apns2.Notification{
        DeviceToken: device.Token,
        Payload:     []byte(payload),
        Topic:       os.Getenv("APNS_BUNDLE_ID"), // ← ★ MissingTopic 対策 ★ 必須 ★
    }

    res, err := h.APNsClient.Push(notification)
    if err != nil {
        fmt.Println("APNs送信エラー:", err)
        return
    }

    if !res.Sent() {
        fmt.Println("APNs失敗:", res.Reason)
        return
    }

    fmt.Println("APNs送信成功 →", userID, name)
}
