package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/models"
	"time"
	"github.com/google/uuid"
)

type StudyLogService struct{}

func (s *StudyLogService) RecordStudyLog(userId string, textbookId string, score float64) error {
	var textbook models.Textbook


	if err := database.DB.Where("id = ?", textbookId).First(&textbook).Error; err != nil {
		return err
	}
	
	var scores []float64 = textbook.ScoreHistory
    scores = append(scores, score) // ★ 新しいスコアを追加
    
    var total float64
    for _, s := range scores {
        total += s
    }

    var accuracy float64
    if len(scores) > 0 {
        accuracy = total / float64(len(scores)) // 新しい平均値を計算
    }

	// 教科書名などは直接取得しているので、マップ生成処理は削除してシンプルにしました
	// (database.GetTextbookIDToNameMap は不要)

	var user models.User
	
	if err := database.DB.Select("name").Where("id = ?", userId).First(&user).Error; err != nil {
		return err
	}
	name := user.Name

	// 暫定: 毎回1になるロジックのようなのでそのまま
	var todayProgress int = 1 

	insertStudyLogs := models.StudyLog{
		ID:            uuid.New().String(),
		UserID:        userId,
		TextbookID:    textbookId,
		Answered:      true,
		AnsweredAt:    time.Now(),
		Score:         score,
		Name:          name,
		TextbookName:  textbook.Name,
		Accuracy:      accuracy,
		TodayProgress: uint(todayProgress),
	}

	if err := database.DB.Create(&insertStudyLogs).Error; err != nil {
		return err
	}

	return nil
}