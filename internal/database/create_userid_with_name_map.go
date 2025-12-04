package database

import (
	"hacku_2025_meijo/internal/models"
)

// GetUserNameMap : ユーザーIDのリストを受け取り、{ID: Name} のマップを返す関数
func GetUserNameMap(userIDs []string) (map[string]string, error) {
	// 結果用のマップを初期化
	nameMap := make(map[string]string)

	// IDリストが空なら、空のマップを返して終了 (DBアクセスしない)
	if len(userIDs) == 0 {
		return nameMap, nil
	}

	// ユーザー情報を取得するためのスライス
	var users []models.User

	// DBから検索
	// Select("id", "name") で、必要なカラムだけ取ってくる（高速化）
	if err := DB.Select("id", "name").Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return nil, err
	}

	// 取得したリストをマップに変換
	for _, u := range users {
		nameMap[u.ID] = u.Name
	}

	return nameMap, nil
}