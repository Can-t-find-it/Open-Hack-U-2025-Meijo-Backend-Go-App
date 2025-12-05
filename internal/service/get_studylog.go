package service

import (
	"fmt"
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
	"time"
)

func GetUserStudyLogs(userID uint) ([]dtos.StudyLogResponse, error) {
	var logs []models.StudyLog
	if err := database.DB.Where("user_id = ?", userID).Find(&logs).Error; err != nil {
		return nil, err
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	responses := make([]dtos.StudyLogResponse, 0)

	for _, l := range logs {
		var q models.Question
		database.DB.First(&q, l.QuestionID)

		var textbook models.Textbook
		database.DB.First(&textbook, q.TextbookID)

		responses = append(responses, dtos.StudyLogResponse{
			ID:           fmt.Sprintf("id%s", l.ID),
			UserName:     user.Name,
			DateTime:     l.AnsweredAt.Format(time.RFC3339),
			TextbookName: textbook.Name,
			Accuracy:     float64(l.Score),
		})
	}

	return responses, nil
}
