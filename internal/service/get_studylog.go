package service

import (
	"fmt"
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
	"time"
	"gorm.io/gorm" // GORMのエラーチェックのために追加
)

// GetLatestStudyLog は、指定されたユーザーの最新の勉強ログ1件のみを取得します。
func GetLatestStudyLog(userID string) (*dtos.StudyLogResponse, error) {
	var log models.StudyLog
    
	// 1. 最新のログ1件を取得 (answered_atの降順でソートし、1件に限定)
	result := database.DB.Where("user_id = ?", userID).
		Order("answered_at DESC").
		Limit(1).
		Find(&log) // Findを使って結果を単一の構造体にバインド

	if result.Error != nil {
		return nil, result.Error
	}
    
    // ログが見つからなかった場合の処理
    if result.RowsAffected == 0 {
        // レコードが見つからない場合はエラーではなく、nilを返すのが一般的
        return nil, gorm.ErrRecordNotFound
    }


	// 2. ユーザー情報を取得
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	// 3. 関連情報を取得 (Question と Textbook)
	var q models.Question
    // First(&q, log.QuestionID) は主キー検索のため Find(&q, ...) でも可
	database.DB.First(&q, log.QuestionID) 

	var textbook models.Textbook
	database.DB.First(&textbook, q.TextbookID)

	// 4. レスポンスDTOを作成
	response := dtos.StudyLogResponse{
		ID: 		  fmt.Sprintf("id%s", log.ID),
		UserName: 	  user.Name,
		DateTime: 	  log.AnsweredAt.Format(time.RFC3339),
		TextbookName: textbook.Name,
		Accuracy: 	  float64(log.Score),
	}

	return &response, nil
}