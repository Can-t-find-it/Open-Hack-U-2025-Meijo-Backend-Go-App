package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

type GetTextbooksService struct{}

func (s *GetTextbooksService) GetTextbooks(userID string) (*dtos.FinalResponseTextbooks, error) {
	user_id := userID
	var friends []models.Friend
	var textbooks []models.Textbook

	if err := database.DB.Where("user_id = ?", user_id).Find(&friends).Error; err != nil {
		return nil, err
	}

	if len(friends) == 0 {
		return &dtos.FinalResponseTextbooks{}, nil
	}

	var friendsID []string
	for _, f := range friends {
		friendsID = append(friendsID, f.FriendUserID)
	}

	if err := database.DB.Where("user_id IN ?", friendsID).Find(&textbooks).Error; err != nil {
		return nil, err
	}

	textMap := make(map[string][]dtos.SoloTextbook)
	var responseTextbooks []dtos.ResponseTextbooks

	for _, t := range textbooks {
		dto := dtos.SoloTextbook{
			TextbookId: t.ID,
			Name:       t.Name,
			PlayTimes:  t.PlayTimes,
			CreatedAt:  t.CreatedAt,
			UpdatedAt:  t.UpdatedAt,
		}
		textMap[t.UserID] = append(textMap[t.UserID], dto)
	}

	friends_id_with_name_map, err := database.GetUserNameMap(friendsID)

	if err != nil {
		return nil, err
	}

	for _, fid := range friendsID {
		book := textMap[fid]
		if book == nil {
			book = []dtos.SoloTextbook{}
		}
		textbook := dtos.ResponseTextbooks{
			FriendId:  fid,
			UserName:  friends_id_with_name_map[fid],
			Textbooks: book,
		}
		responseTextbooks = append(responseTextbooks, textbook)
	}

	response := dtos.FinalResponseTextbooks{
		Friends: responseTextbooks,
	}
	return &response, nil
}
