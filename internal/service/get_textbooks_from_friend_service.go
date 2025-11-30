package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

type GetTextbooksService struct{}

func (s *GetTextbooksService) GetTextbooks(userID uint) ([]dtos.ResponseTextbooks, error) {
	user_id := userID
	var friends []models.Friend
	var textbooks []models.Textbook

	if err := database.DB.Where("user_id = ?", user_id).Find(&friends).Error; err != nil {
		return nil, err
	}

	if len(friends) == 0 {
		return []dtos.ResponseTextbooks{}, nil
	}

	var friendsID []uint
	for _, f := range friends {
		friendsID = append(friendsID, f.FriendUserID)
	}

	if err := database.DB.Where("user_id IN ?", friendsID).Find(&textbooks).Error; err != nil {
		return nil, err
	}

	textMap := make(map[uint][]dtos.SoloTextbook)
	var responseTextbooks []dtos.ResponseTextbooks

	for _, t := range textbooks {
		dto := dtos.SoloTextbook{
			UserID:    t.UserID,
			Name:      t.Name,
			FolderID:  t.FolderID,
			Type:      string(t.Type),
			PlayTimes: t.PlayTimes,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		}
		textMap[t.UserID] = append(textMap[t.UserID], dto)
	}

	for _, fid := range friendsID {
		book := textMap[fid]
		if book == nil {
			book = []dtos.SoloTextbook{}
		}
		textbook := dtos.ResponseTextbooks{
			UserId:    fid,
			Textbooks: book,
		}
		responseTextbooks = append(responseTextbooks, textbook)
	}
	return responseTextbooks, nil
}
