package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/models"
	"time"
)

type StudyLogService struct{}

func (s *StudyLogService) RecordStudyLog(userId uint, textbookId uint, score float64) error {
	var textbook models.Textbook

	if err := database.DB.Where("id = ?", textbookId).First(&textbook).Error; err != nil {
		return err
	}
	var scores []float64
	scores = textbook.ScoreHistory

	var tmp float64

	for _, i := range scores {
		tmp += i
	}

	accuracy := tmp / float64(len(scores))

	var textbookIdSlice []uint

	textbookIdSlice = append(textbookIdSlice, textbookId)

	textbookMap, err := database.GetTextbookIDToNameMap(textbookIdSlice)
	if err != nil {
		return err
	}

	var user models.User
	if err := database.DB.Select("name").First(&user, userId).Error; err != nil {
		return err
	}
	name := user.Name

	var todayProgress int

	todayProgress++

	insertStudyLogs := models.StudyLog{
		UserID:        userId,
		TextbookID:    textbookId,
		Answered:      true,
		AnsweredAt:    time.Now(),
		Score:         score,
		Name:          name,
		TextbookName:  textbookMap[textbookId],
		Accuracy:      accuracy,
		TodayProgress: uint(todayProgress),
	}

	if err := database.DB.Create(&insertStudyLogs).Error; err != nil {
		return err
	}

	return nil

}
