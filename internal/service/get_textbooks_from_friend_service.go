package service

import (
	"fmt"
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
)

// GetTextbooksService 構造体定義
type GetTextbooksService struct{}

// 質問数集計クエリの結果を格納するための内部構造体
type QuestionCountResult struct {
	TextbookID string
	Count      int
}

// getQuestionCounts は、教科書IDのリストに基づいて、各教科書の質問数をデータベースから一度に取得します。
func (s *GetTextbooksService) getQuestionCounts(textbookIDs []string) (map[string]int, error) {
	if len(textbookIDs) == 0 {
		return map[string]int{}, nil
	}

	var results []QuestionCountResult
	
    // SELECT textbook_id, count(id) FROM questions WHERE textbook_id IN (?) GROUP BY textbook_id
	err := database.DB.Model(&models.Question{}).
		Select("textbook_id, count(id) as count").
		Where("textbook_id IN (?)", textbookIDs).
		Group("textbook_id").
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[string]int)
	for _, r := range results {
		countMap[r.TextbookID] = r.Count
	}
	return countMap, nil
}


// GetTextbooks はフレンドの教科書を取得し、フレンドごとにグループ化して返します。
func (s *GetTextbooksService) GetTextbooks(userID string) (*dtos.FinalResponseTextbooks, error) {
	user_id := userID
	fmt.Printf("GetTextbooksService called with userID: %s\n", user_id)

	var friends []models.Friend
	var textbooks []models.Textbook

	// 1. フレンド関係の取得
	if err := database.DB.Where("user_id = ?", user_id).Find(&friends).Error; err != nil {
		return nil, err
	}

	// 2. フレンドがいない場合は空のレスポンスを返す (早期リターン)
	if len(friends) == 0 {
		return &dtos.FinalResponseTextbooks{
			Friends: []dtos.ResponseTextbooks{},
		}, nil
	}

	var friendsID []string
	for _, f := range friends {
		friendsID = append(friendsID, f.FriendUserID)
		fmt.Printf("Found friend with UserID: %s\n", f.FriendUserID)
	}

    // 3. 全フレンドの教科書を一度に取得
	if err := database.DB.Where("user_id IN ?", friendsID).Find(&textbooks).Error; err != nil {
		fmt.Printf("Error fetching textbooks for friends: %v\n", err)
		return nil, err
	}
	fmt.Printf("Number of textbooks found: %d\n", len(textbooks))


	// 教科書が空の場合でも、以降の処理で空の教科書リストを持つフレンドリストを返します

	// 4. 質問数を効率的に取得 (TextbookID -> QuestionCount のマップ)
	textbookIDs := make([]string, 0, len(textbooks))
	for _, t := range textbooks {
		textbookIDs = append(textbookIDs, t.ID)
	}
	
	questionCounts, err := s.getQuestionCounts(textbookIDs)
	if err != nil {
		fmt.Printf("Error fetching question counts: %v\n", err)
		return nil, err
	}

	// 5. 教科書を UserID (フレンドID) ごとにグループ化
	textMap := make(map[string][]dtos.SoloTextbook)
	for _, t := range textbooks {
		dto := dtos.SoloTextbook{
			TextbookId:    t.ID,
			Name:          t.Name,
			QuestionCount: questionCounts[t.ID], // 質問数を設定
			CreatedAt:     t.CreatedAt,
			UpdatedAt:     t.UpdatedAt,
		}
		textMap[t.UserID] = append(textMap[t.UserID], dto)
	}

    // 6. フレンドのユーザー名を取得
	friends_id_with_name_map, err := database.GetUserNameMap(friendsID)
	if err != nil {
		return nil, err
	}

    // 7. 最終的なレスポンス構造にマッピング
	responseTextbooks := make([]dtos.ResponseTextbooks, 0, len(friendsID))
	for _, fid := range friendsID {
		book := textMap[fid]
		if book == nil {
			book = []dtos.SoloTextbook{} // 教科書がない場合は空リスト
		}
		
		textbook := dtos.ResponseTextbooks{
			FriendId:  fid,
			UserName:  friends_id_with_name_map[fid],
			Textbooks: book,
		}
		responseTextbooks = append(responseTextbooks, textbook)
	}

	response := dtos.FinalResponseTextbooks{
		Friends: responseTextbooks,
	}
	return &response, nil
}