package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

type GetFriendStudyLogService struct{}

func (s *GetFriendStudyLogService) GetFriendStudyLog(userId string) (*dtos.AllFriendsStudyLog, error) {
	var allStudyLogs []models.StudyLog

	var myFriends []models.Friend

	var myFriendsID []string
	
	if err := database.DB.Where("user_id = ?", userId).Find(&myFriends).Error; err != nil {
		return nil, err
	}
	
	for _ ,f := range myFriends {
		myFriendsID = append(myFriendsID, f.FriendUserID)
	}

	if err := database.DB.Where("user_id IN ?", myFriendsID).Find(&allStudyLogs).Error; err != nil {
		return nil, err
	}

	var studyLogs []dtos.FriendStudyLog

	for _, log := range allStudyLogs {
		dto := dtos.FriendStudyLog{
			ID: log.ID,
			FriendID:   log.UserID,
			FriendName: log.Name,
			AnsweredAt: log.AnsweredAt,
			TextbookName: log.TextbookName,
			Accuracy: log.Accuracy,
			TodayProgress: log.TodayProgress,
		}
		studyLogs = append(studyLogs, dto)
	}

	response := dtos.AllFriendsStudyLog{
		Logs: studyLogs,
	}

	return &response, nil
}
