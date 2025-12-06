package service

import (
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TextbookCopyService struct{}

func (s *TextbookCopyService) AddTextbook(userID string, targetFolderID string, sourceTextbookID string) error {

	return database.DB.Transaction(func(tx *gorm.DB) error {

		var sourceBook models.Textbook
		// 1. コピー元の教科書を取得
		if err := tx.Where("id = ?", sourceTextbookID).First(&sourceBook).Error; err != nil {
			return err
		}

		// 2. 新しい教科書を作成 (コピー)
		newBook := models.Textbook{
			ID:         	 uuid.New().String(),
			UserID:          userID,         // 自分のID
			FolderID:        targetFolderID, // 保存先フォルダ
			Name:            sourceBook.Name,
			Type:            sourceBook.Type,
			StudyMaterialID: sourceBook.StudyMaterialID,
		}

		if err := tx.Create(&newBook).Error; err != nil {
			return err
		}

		// 3. 元の教科書に紐づく「問題」を、「問題文(QuestionStatements)」ごと取得する
		var questions []models.Question
		// Preloadを使うことで、子テーブル(QuestionStatements)のデータも一緒に引っ張ってくる
		if err := tx.Where("textbook_id = ?", sourceBook.ID).
			Preload("QuestionStatements").
			Find(&questions).Error; err != nil {
			return err
		}

		// 4. 問題と問題文をコピー
		for _, q := range questions {
			// 新しい問題データを作成
			newQuestion := models.Question{
				ID:         uuid.New().String(),
				TextbookID: newBook.ID, // ★新しい教科書に紐付け
				Answer:     q.Answer,
				// ID: 0, CreatedAtなどはリセット
			}

			// 問題を保存 (ここで newQuestion.ID が発行される)
			if err := tx.Create(&newQuestion).Error; err != nil {
				return err
			}

			// 5. その問題に紐づく「問題文」をコピー
			if len(q.QuestionStatements) > 0 {
				var newStatements []models.QuestionStatement

				for _, stmt := range q.QuestionStatements {
					// 構造体をコピー
					newStmt := stmt

					// 重要なID部分を書き換える
					newStmt.ID = ""                      // 新規作成にする
					newStmt.QuestionID = newQuestion.ID // 新しい問題に紐付け

					newStatements = append(newStatements, newStmt)
				}

				// 問題文を一括保存
				if err := tx.Create(&newStatements).Error; err != nil {
					return err
				}
			}
		}

		return nil // コミット
	})
}
