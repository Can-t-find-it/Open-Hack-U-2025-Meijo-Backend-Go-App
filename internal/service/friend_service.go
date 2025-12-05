package service

import (
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
	"hacku_2025_meijo/internal/database"
)

// FriendService はフレンド関連のビジネスロジックを担うサービス構造体です。
type FriendService struct {
    // 依存関係（例：DB接続など）が必要ならここに定義
}

// NewFriendService は FriendService のコンストラクタです。
func NewFriendService() *FriendService {
    return &FriendService{}
}

// SearchFriendsByName は、ユーザーのフレンドリストから名前を検索し、一致するフレンドのDTOスライスを返します。
// ログインユーザーのフレンドであることと、名前に部分一致することを、JOINを使って一度のクエリで処理します。
func (s *FriendService) SearchFriendsByName(userID string, searchName string) ([]dtos.FriendSearchResponse, error) {
	var responses []dtos.FriendSearchResponse
	
	searchPattern := "%" + searchName + "%"

	// GORMのJOINを使って、フレンド関係とユーザー情報を結合し、一度に検索する
	if err := database.DB.Model(&models.Friend{}).
		// ログインユーザーに絞り込む
		Where("friends.user_id = ?", userID).
		// Userテーブルと結合 (フレンドユーザーID == UserID)
		Joins("INNER JOIN users ON users.id = friends.friend_user_id").
		// 結合したUserテーブルの名前を部分一致で検索
		Where("users.name LIKE ?", searchPattern).
		// 必要なカラムを選択し、DTOのスライスに直接格納する
		Select(
			"friends.id AS id", 
			"friends.friend_user_id AS friend_user_id",
			"users.name AS name",
		).
		// 取得結果をdtos.FriendSearchResponseのスライスに直接スキャン
		Find(&responses).Error; err != nil {
		return nil, err
	}
    
	return responses, nil
}