package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

type ChangeFriendsService struct{}

func (s *ChangeFriendsService) AddFriends(userID uint, input dtos.AddFriends) (*dtos.ResponseAddFriends, error) {
	// 1. すでに登録済みの友達IDを取得する
	var existingFriends []models.Friend
	// "user_id = ?" で自分の友達データを検索し、"friend_user_id" カラムだけを取得する
	if err := database.DB.Where("user_id = ?", userID).Find(&existingFriends).Error; err != nil {
		return nil, err
	}

	// 検索を高速にするために、取得したIDをマップ（辞書）に変換する
	// map[友達ID]bool という形にする
	existingMap := make(map[uint]bool)
	for _, f := range existingFriends {
		existingMap[f.FriendUserID] = true
	}

	var friendsName []models.User

	if err := database.DB.Where("id IN ?", input.Friends).Find(&friendsName).Error; err != nil {
		return nil, err
	}

	userNameMap := make(map[uint]string)
	for _, i := range friendsName {
		userNameMap[i.ID] = i.Name
	}

	// 2. 新しい友達リストを作成
	var newFriends []models.Friend

	for _, friendIDInt := range input.Friends {
		targetID := uint(friendIDInt)

		// 自分自身ならスキップ
		if userID == targetID {
			continue
		}

		if existingMap[targetID] {
			continue
		}

		friendsName := userNameMap[targetID]

		// 重複していない場合のみリストに追加
		friend := models.Friend{
			UserID:       userID,
			FriendUserID: targetID,
			Name:         friendsName,
		}
		newFriends = append(newFriends, friend)
	}

	// 3. まとめて保存
	if len(newFriends) > 0 {
		if err := database.DB.Create(&newFriends).Error; err != nil {
			return nil, err
		}
	}
	var friend []dtos.SoloAddFriend

	for _, f := range newFriends {
		friend = append(friend, dtos.SoloAddFriend{
			Friend: int(f.FriendUserID),
			Name:   f.Name,
		})

	}

	return &dtos.ResponseAddFriends{
		Friends: friend,
	}, nil
}

// DeleteFriendsBatch : 指定された複数の友達を一括削除する
func (s *ChangeFriendsService) DeleteFriendsBatch(userID uint, targetFriendIDs []int) error {
	if len(targetFriendIDs) == 0 {
		return nil
	}
	result := database.DB.Where("user_id = ? AND friend_user_id IN ?", userID, targetFriendIDs).Delete(&models.Friend{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *ChangeFriendsService) GetFriends(userID uint) (*dtos.GetFriends, error) {
	// 1. すでに登録済みの友達IDを取得する
	var existingFriends []models.Friend
	// "user_id = ?" で自分の友達データを検索し、"friend_user_id" カラムだけを取得する
	if err := database.DB.Where("user_id = ?", userID).Find(&existingFriends).Error; err != nil {
		return nil, err
	}

	var friend []dtos.SoloGetFriend
	for _, f := range existingFriends {
		friend = append(friend, dtos.SoloGetFriend{
			FriendsID: int(f.FriendUserID),
			Name:      f.Name,
		})
	}

	// 3. DTO（レスポンス用の箱）を作って返す
	return &dtos.GetFriends{
		Friends: friend,
	}, nil
}
