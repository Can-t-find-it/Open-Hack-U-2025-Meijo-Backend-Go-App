package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

type GetFriendStudyLogService struct{}

func (s *GetFriendStudyLogService) GetFriendStudyLog(FriendsID []uint) ([]dtos.ResponseFriendStudyLog, error) {
	var allStudyLogs []models.StudyLog

	if err := database.DB.Where("user_id IN ?", FriendsID).Find(&allStudyLogs).Error; err != nil {
		return nil, err
	}

	var response []dtos.ResponseFriendStudyLog

	for _, log := range allStudyLogs {
		dto := dtos.ResponseFriendStudyLog{
			FriendID:   log.UserID,
			QuestionID: log.QuestionID,
			Answered:   log.Answered,
			AnsweredAt: log.AnsweredAt,
			Score:      log.Score,
		}
		response = append(response, dto)
	}
	return response, nil
}
