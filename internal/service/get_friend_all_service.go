package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

type GetFriendService struct{}

func (s *GetFriendService) GetAllFriends(userID int) (*dtos.ResponseFriendGetAll, error) {
	var friends []models.Friend

	result := database.DB.Where("user_id = ?", userID).Find(&friends)

	if result.Error != nil {
		return nil, result.Error
	}

	var friendsID []int

	for _, f := range friends {
		friendsID = append(friendsID, int(f.FriendUserID))
	}
	return &dtos.ResponseFriendGetAll{
		User:          userID,
		FriendsUserID: friendsID,
	}, nil
}
