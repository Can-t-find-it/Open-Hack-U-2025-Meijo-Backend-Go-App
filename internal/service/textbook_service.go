package service

import (
	"fmt"
	"hacku_2025_meijo/internal/database" // DB接続変数 (DB) がある場所
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

// ==========================================
// 教科書・フォルダ操作 (DB関連)
// ==========================================

// GetUserTextbooks: ユーザーのフォルダと、その中の問題集一覧を全部取得
func GetUserTextbooks(userID string) ([]dtos.FolderResponse, error) {
	var folders []models.Folder

	// 1. DBから全データを取得 (ここは変わらず)
	result := database.DB.
		Preload("Textbooks").
		Where("user_id = ?", userID).
		Find(&folders)

	if result.Error != nil {
		return nil, result.Error
	}

	// 2. 必要なデータだけを「詰め替え」作業
	response := []dtos.FolderResponse{}

	for _, f := range folders {
		// 教科書リストの詰め替え
		var currentTextbooks []dtos.TextbookResponse
		if f.Textbooks == nil {
			currentTextbooks = []dtos.TextbookResponse{}
		}
		for _, t := range f.Textbooks {
			currentTextbooks = append(currentTextbooks, dtos.TextbookResponse{
				ID:   t.ID,
				Name: t.Name,
				Type: string(t.Type),
			})
		}

		// フォルダの詰め替え
		response = append(response, dtos.FolderResponse{
			ID:        f.ID,
			Name:      f.Name,
			Progress:  f.Progress,
			Textbooks: currentTextbooks,
		})
	}

	return response, nil
}

// GetTextbookDetail: 問題集の詳細（中の問題リスト含む）を取得
func GetTextbookDetail(textbookID string) (*dtos.TextbookDetailResponse, error) {
	var textbook models.Textbook

	// Textbook -> Questions -> QuestionStatements という深い階層を取得
	result := database.DB.
		Preload("Questions.QuestionStatements").
		First(&textbook, "id = ?", textbookID)

	if result.Error != nil {
		return nil, result.Error
	}
	
	var questionsResp []dtos.QuestionResponse
	for _, q := range textbook.Questions {
		var statementsResp []dtos.QuestionStatementResponse
		for _, s := range q.QuestionStatements {
			statementsResp = append(statementsResp, dtos.QuestionStatementResponse{
				ID:                s.ID,
				QuestionStatement: s.Statement,
				Choices:           s.Choices,
				Explain:          s.Explain,
			})
		}
		questionsResp = append(questionsResp, dtos.QuestionResponse{
			ID:                 q.ID,
			Answer:             q.Answer,
			QuestionStatements: statementsResp,
		})
	}

	response := &dtos.TextbookDetailResponse{
		ID:        textbook.ID,
		Name:      textbook.Name,
		Type:      string(textbook.Type),
		Questions: questionsResp,
		Score:     textbook.ScoreHistory,
		Times:     textbook.PlayTimes,
	}

	return response, nil
}

// CreateTextbook: 新しい問題集を作成
// 引数の typeStr は、ユーザーからの入力なので string のままでOK
func CreateTextbook(name string, typeStr string, folderID string) error {
	
	// 1. 入力された文字列を、専用の型にキャスト（変換）してみる
	inputType := models.TextbookType(typeStr)

	// 2. 許可リスト（スイッチ文を使うとスッキリします）
	switch inputType {
	case models.Type4Choice, models.TypeFillIn, models.TypeFillIn4Choice, models.TypeInput:
		// OK！何もしない
	default:
		// NG！
		return fmt.Errorf("無効なタイプです: %s", typeStr)
	
	}

	// 3. 保存処理
	newTextbook := models.Textbook{
		Name:     name,
		Type:     inputType, // ここで型付きのデータを渡す
		FolderID: folderID,
	}
	result := database.DB.Create(&newTextbook)
	return result.Error
}

// ==========================================
// AI生成問題の保存機能
// ==========================================

// AddQuestionToTextbook: 生成された問題データをDB構造に変換して保存
func AddQuestionToTextbook(textbookID string, item dtos.ResultItem, answer string) error {
	// ★修正: 自分で保存処理を書かず、重複チェック機能付きの SaveQuestionToDB を呼ぶだけにする
	_, err := SaveQuestionToDB(textbookID, item, answer)
	return err
}

// ----------------------------------------------------
// 未実装だった削除・追加機能 (Question / QuestionStatement)
// ----------------------------------------------------

// DeleteQuestions: 問題を削除する
func DeleteQuestion(questionID string) error {
	// 指定されたIDリストの問題をまとめて削除
	// OnDelete:CASCADE設定により、紐付くQuestionStatementも自動で消えます
	result := database.DB.Unscoped().Where("id = ?", questionID).Delete(&models.Question{})
	return result.Error
}

// AddQuestionStatement: 既存の問題（親）に、新しい聞き方（子）を追加する
// 例: 「1+1は？」という問題(ID:10)に対して、「1に1を足すと？」という別パターンを追加
func AddQuestionStatement(questionID string, statement string, explain string, choices []string) error {
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

// GetRandomSuggestedWords: ランダムに指定個数の単語を取得する
func GetRandomSuggestedWords(limit int) ([]string, error) {
	var questions []models.Question
	
	// Limit(limit) で個数を指定
	result := database.DB.Order("RANDOM()").Limit(limit).Find(&questions)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	var words []string
	for _, q := range questions {
		words = append(words, q.Answer)
	}
	return words, nil
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
func DeleteFoldersBatch(folderIDs []string) error {
	// Unscoped() をつけると、論理削除ではなく物理削除（完全に消す）になります
	// 構成によっては Unscoped なしでもOKですが、今回は確実に消すために付けます
	result := database.DB.Unscoped().Delete(&models.Folder{}, "id IN ?", folderIDs)
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

func DeleteTextbook(textbookIDs string) error {
	// 指定されたIDのTextbookを削除
	// OnDelete:CASCADE設定があるため、中身の問題なども連鎖して削除されます
	if err := database.DB.Delete(&models.Textbook{}, "id = ?", textbookIDs).Error; err != nil {
		return err
	}
	return nil
}