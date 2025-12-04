package service

import (
	"hacku_2025_meijo/internal/database" // DB接続変数 (DB) がある場所
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

// ==========================================
// 教科書・フォルダ操作 (DB関連)
// ==========================================

// GetUserTextbooks: ユーザーのフォルダと問題集一覧を全部取得
func GetUserTextbooks(userID uint) ([]models.Folder, error) {
	var folders []models.Folder

	// Preloadを使って、Folder -> Textbooks まで一気に取得する
	// DB変数は internal/database/db.go で定義されている "DB" を使用
	result := database.DB.
		Preload("Textbooks").
		Where("user_id = ?", userID).
		Find(&folders)

	if result.Error != nil {
		return nil, result.Error
	}
	return folders, nil
}

// GetTextbookDetail: 問題集の詳細（中の問題リスト含む）を取得
func GetTextbookDetail(textbookID string) (*models.Textbook, error) {
	var textbook models.Textbook

	// Textbook -> Questions -> QuestionStatements という深い階層を取得
	result := database.DB.
		Preload("Questions.QuestionStatements").
		First(&textbook, "id = ?", textbookID)

	if result.Error != nil {
		return nil, result.Error
	}
	return &textbook, nil
}

// CreateTextbook: 新しい問題集を作成
func CreateTextbook(name string, typeStr string, folderID uint) error {
	newTextbook := models.Textbook{
		Name:     name,
		Type:     typeStr,
		FolderID: folderID,
	}
	result := database.DB.Create(&newTextbook)
	return result.Error
}

// DeleteTextbook: 問題集を削除
func DeleteTextbook(textbookID string) error {
	// Cascade設定がModelにあれば、関連するQuestionも消える
	result := database.DB.Delete(&models.Textbook{}, "id = ?", textbookID)
	return result.Error
}

// ==========================================
// AI生成問題の保存機能
// ==========================================

// AddQuestionToTextbook: 生成された問題データをDB構造に変換して保存
func AddQuestionToTextbook(textbookID uint, item dtos.ResultItem, answer string) error {

	// 1. Question（親：正解データ）を作成
	question := models.Question{
		TextbookID: textbookID,
		Answer:     answer, // 正解の文字列（例: "リンゴ"）
		// QuestionStatements（子：出題文・選択肢）をネストして作成
		QuestionStatements: []models.QuestionStatement{
			{
				Statement: item.Question,    // 問題文
				Explain:   item.Explanation, // 解説
				Choices:   item.Options,     // 選択肢配列 (JSONとして保存される)
			},
		},
	}

	// 2. DBに保存
	// 親のQuestionを保存すれば、子のStatementも自動で保存されます
	if err := database.DB.Create(&question).Error; err != nil {
		return err
	}

	return nil
}

// ----------------------------------------------------
// 未実装だった削除・追加機能 (Question / QuestionStatement)
// ----------------------------------------------------

// DeleteQuestion: 問題（親）を削除
// これを消すと、紐付いている問題文（子）も全部消えます（CASCADE設定のため）
func DeleteQuestion(questionID string) error {
	result := database.DB.Delete(&models.Question{}, "id = ?", questionID)
	return result.Error
}

// AddQuestionStatement: 既存の問題（親）に、新しい聞き方（子）を追加する
// 例: 「1+1は？」という問題(ID:10)に対して、「1に1を足すと？」という別パターンを追加
func AddQuestionStatement(questionID uint, statement string, explain string, choices []string) error {
	newStatement := models.QuestionStatement{
		QuestionID: questionID,
		Statement:  statement,
		Explain:    explain,
		Choices:    choices,
	}
	result := database.DB.Create(&newStatement)
	return result.Error
}

// DeleteQuestionStatement: 特定の問題文（子）だけを削除する
func DeleteQuestionStatement(statementID string) error {
	result := database.DB.Delete(&models.QuestionStatement{}, "id = ?", statementID)
	return result.Error
}

// GetSuggestedWord: 覚えたい単語を提案する
// 今回は仮として「DBにある問題の正解」からランダムに1つ取得して返します
func GetSuggestedWord() (string, error) {
	var question models.Question

	// ORDER BY RANDOM() でランダムに1件取得 (PostgreSQLの場合)
	// MySQLなら "RAND()" を使用
	result := database.DB.Order("RANDOM()").First(&question)

	if result.Error != nil {
		return "", result.Error
	}

	return question.Answer, nil // 正解の単語を返す
}

// UpdateTextbookStatus: 教科書のスコアと回数を更新する（学習後に呼ぶ用）
func UpdateTextbookStatus(textbookID string, newScore float64) error {
	var textbook models.Textbook
	if err := database.DB.First(&textbook, "id = ?", textbookID).Error; err != nil {
		return err
	}

	// 回数を+1
	textbook.PlayTimes += 1
	// スコア履歴に追加
	textbook.ScoreHistory = append(textbook.ScoreHistory, newScore)

	return database.DB.Save(&textbook).Error
}

// DeleteFolder: フォルダを削除する（中身も全消去）
func DeleteFolder(folderID string) error {
	// Unscoped() をつけると、論理削除ではなく物理削除（完全に消す）になります
	// 構成によっては Unscoped なしでもOKですが、今回は確実に消すために付けます
	result := database.DB.Unscoped().Delete(&models.Folder{}, "id = ?", folderID)
	return result.Error
}

// GetWordsInTextbook: 教科書に含まれる全ての正解単語を取得する
func GetWordsInTextbook(textbookID string) ([]string, error) {
	var questions []models.Question

	// 指定された教科書のQuestionを全部取ってくる
	result := database.DB.Where("textbook_id = ?", textbookID).Find(&questions)
	if result.Error != nil {
		return nil, result.Error
	}

	var words []string
	for _, q := range questions {
		words = append(words, q.Answer)
	}
	return words, nil
}
