package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
	"time"
)

// GetUserStatus: ユーザーの詳細なステータスを取得する
func GetUserStatus(userID string) (*dtos.UserStatusResponse, error) {
	// 1. ユーザー基本情報の取得
	var user models.User
	// "id = ?" で安全に検索
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 2. 教科書数のカウント
	var textbookCount int64
	database.DB.Model(&models.Textbook{}).Where("user_id = ?", userID).Count(&textbookCount)

	// 3. フレンド数のカウント
	var friendCount int64
	database.DB.Model(&models.Friend{}).Where("user_id = ?", userID).Count(&friendCount)

	// 4. 連続学習日数 (StreakDays) の計算
	streakDays := calculateStreakDays(userID)

	// 5. レスポンスの作成
	return &dtos.UserStatusResponse{
		ID:   user.ID,
		Name: user.Name,
		Status: dtos.UserStatus{
			TextbookCount: textbookCount,
			StreakDays:    streakDays,
			FriendCount:   friendCount,
		},
	}, nil
}

// calculateStreakDays: 勉強ログから連続日数を計算する
func calculateStreakDays(userID string) int {
	var logs []models.StudyLog
	
	// 日付順（新しい順）に取得
	database.DB.Select("answered_at").
		Where("user_id = ?", userID).
		Order("answered_at DESC").
		Find(&logs)

	if len(logs) == 0 {
		return 0
	}

	// 重複を除いた「勉強した日付」のリストを作る
	uniqueDates := make(map[string]bool)
	var dateList []time.Time
	
	for _, log := range logs {
		dateStr := log.AnsweredAt.Local().Format("2006-01-02")
		
		if !uniqueDates[dateStr] {
			uniqueDates[dateStr] = true
			y, m, d := log.AnsweredAt.Local().Date()
			dateList = append(dateList, time.Date(y, m, d, 0, 0, 0, 0, time.Local))
		}
	}

	// 連続チェック
	streak := 0
	now := time.Now().Local()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	
	checkDate := today

	// 最新の学習日が「昨日」よりも前なら途切れている
	if len(dateList) > 0 {
		lastStudy := dateList[0]
		daysDiff := int(today.Sub(lastStudy).Hours() / 24)
		
		if daysDiff > 1 {
			return 0
		}
		if daysDiff == 1 {
			checkDate = checkDate.AddDate(0, 0, -1)
		}
	}

	for _, d := range dateList {
		if d.Equal(checkDate) {
			streak++
			checkDate = checkDate.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	return streak
}