package database

import (
	"hacku_2025_meijo/internal/models"
)

func GetTextbookNameMap(questionIDs []string) (map[string]string, error) {
	// 結果用のマップを初期化
	nameMap := make(map[string]string)

	// IDリストが空なら、空のマップを返して終了 (無駄なDBアクセスを防ぐ)
	if len(questionIDs) == 0 {
		return nameMap, nil
	}

	// データを取得するためのスライス
	var questions []models.Question

	// DBから検索
	// Preload("Textbook") : 問題に紐付いている教科書データも一緒に持ってくる
	// Select("id", "textbook_id") : 必要なカラムだけ指定して高速化（Textbook側は全カラム取得される）
	if err := DB.Preload("Textbook").
		Where("id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		return nil, err
	}

	// 取得したリストをマップに変換
	for _, q := range questions {
		// q.Textbook.Name で教科書名が取れる
		nameMap[q.ID] = q.Textbook.Name
	}

	return nameMap, nil
}

// GetTextbookIDToNameMap : 教科書IDのリストを受け取り、{TextbookID: TextbookName} のマップを返す関数
func GetTextbookIDToNameMap(textbookIDs []string) (map[string]string, error) {
	// 結果用のマップを初期化
	nameMap := make(map[string]string)

	// IDリストが空なら、空のマップを返して終了 (DBアクセスしない)
	if len(textbookIDs) == 0 {
		return nameMap, nil
	}

	// 教科書情報を取得するためのスライス
	var textbooks []models.Textbook

	// DBから検索
	// Select("id", "name") で、必要なカラムだけ取ってくる（高速化）
	if err := DB.Select("id", "name").Where("id IN ?", textbookIDs).Find(&textbooks).Error; err != nil {
		return nil, err
	}

	// 取得したリストをマップに変換
	for _, t := range textbooks {
		nameMap[t.ID] = t.Name
	}

	return nameMap, nil
}

// GetTextbookOwnerMap : 教科書IDのリストを受け取り、{TextbookID: UserID} のマップを返す関数
// これを使えば、「この教科書は誰のもの？」を一発で取得できます。
func GetTextbookOwnerMap(textbookIDs []string) (map[string]string, error) {
	// 結果用のマップを初期化
	ownerMap := make(map[string]string)

	// IDリストが空なら、空のマップを返して終了 (無駄なDBアクセスを防ぐ)
	if len(textbookIDs) == 0 {
		return ownerMap, nil
	}

	// データを取得するためのスライス (Textbookモデル)
	var textbooks []models.Textbook

	// DBから検索
	// Select("id", "user_id") で必要なカラムだけ取ってくる (高速化)
	if err := DB.Select("id", "user_id").Where("id IN ?", textbookIDs).Find(&textbooks).Error; err != nil {
		return nil, err
	}

	// 取得したリストをマップに変換
	for _, t := range textbooks {
		ownerMap[t.ID] = t.UserID
	}

	return ownerMap, nil
}
