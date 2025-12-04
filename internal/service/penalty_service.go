package service

import (
	"fmt"
	"math/rand"
	"time"

	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/models"
)


// サボり判定 
func CheckPenaltyAndNotify() error {

	threeDaysAgo := time.Now().Add(-72 * time.Hour)

	var userIDs []string

	// 勉強していないユーザーを取得
	err := database.DB.
		Table("study_logs").
		Select("user_id").
		Group("user_id").
		Having("MAX(answered_at) < ?", threeDaysAgo).
		Find(&userIDs).Error

	if err != nil {
		return err
	}

	if len(userIDs) == 0 {
		fmt.Println("サボっているユーザーはいません。")
		return nil
	}

	for _, uid := range userIDs {
		notifyFriends(uid)
	}

	return nil
}


// フレンドを取得してランダムに通知
func notifyFriends(targetUserID string) {
	var friends []models.Friend

	// 通知ONのフレンドを取得
	database.DB.Where("user_id = ? AND notify_enabled = TRUE", targetUserID).
		Find(&friends)

	if len(friends) == 0 {
		fmt.Printf("ユーザー %s に通知対象フレンドなし\n", targetUserID)
		return
	}

	selected := pickRandomFriends(friends, 2)

	for _, f := range selected {
		sendPenaltyPush(f.FriendUserID, targetUserID)
	}
}


// ランダムに選ぶ
func pickRandomFriends(friends []models.Friend, n int) []models.Friend {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(friends), func(i, j int) {
		friends[i], friends[j] = friends[j], friends[i]
	})

	if len(friends) < n {
		return friends
	}
	return friends[:n]
}

// FCM 通知
func sendPenaltyPush(toUserID string, saboUserID string) {

	message := fmt.Sprintf(
		"あなたの友達（ID: %s）が最近学習していません！励ましてあげてください！",
		saboUserID,
	)

	fmt.Printf("[通知] ToUser=%s : %s\n", toUserID, message)
}
