package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

type ChangeFriendsService struct{}

func (s *ChangeFriendsService) AddFriends(userID uint, input dtos.AddFriends) (*dtos.AddFriends, error) {
	// 1. すでに登録済みの友達IDを取得する
	var existingFriends []models.Friend
	// "user_id = ?" で自分の友達データを検索し、"friend_user_id" カラムだけを取得する
	if err := database.DB.Select("friend_user_id").Where("user_id = ?", userID).Find(&existingFriends).Error; err != nil {
		return nil, err
	}

	// 検索を高速にするために、取得したIDをマップ（辞書）に変換する
	// map[友達ID]bool という形にする
	existingMap := make(map[uint]bool)
	for _, f := range existingFriends {
		existingMap[f.FriendUserID] = true
	}

	// 2. 新しい友達リストを作成
	var newFriends []models.Friend

	for _, friendIDInt := range input.Friends {
		targetID := uint(friendIDInt)

		// 自分自身ならスキップ
		if userID == targetID {
			continue
		}

		// ★ここが追加部分：すでに友達ならスキップ
		if existingMap[targetID] {
			continue
		}

		// 重複していない場合のみリストに追加
		friend := models.Friend{
			UserID:       userID,
			FriendUserID: targetID,
		}
		newFriends = append(newFriends, friend)
	}

	// 3. まとめて保存
	if len(newFriends) > 0 {
		if err := database.DB.Create(&newFriends).Error; err != nil {
			return nil, err
		}
	}

	return &input, nil
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
	if err := database.DB.Select("friend_user_id").Where("user_id = ?", userID).Find(&existingFriends).Error; err != nil {
		return nil, err
	}

	var friendIDs []int
	for _, f := range existingFriends {
		// uint から int へ型変換して追加
		friendIDs = append(friendIDs, int(f.FriendUserID))
	}

	// 3. DTO（レスポンス用の箱）を作って返す
	response := dtos.GetFriends{
		GetFriends: friendIDs,
	}

	return &response, nil
}
