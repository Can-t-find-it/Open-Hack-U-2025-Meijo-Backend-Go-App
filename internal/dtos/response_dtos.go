package dtos

// 教科書一覧のレスポンス用 (中身を限定する)
type FolderResponse struct {
	ID        uint               `json:"id"`
	Name      string             `json:"name"`
	Progress  int                `json:"progress"`
	Textbooks []TextbookResponse `json:"textbooks"`
}

type TextbookResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type TextbookDetailResponse struct {
	ID        uint                     `json:"id"`
	Name      string                   `json:"name"`
	Type      string                   `json:"type"`
	Questions []QuestionResponse `json:"questions"`
	Score     []float64                `json:"score"`
	Times     int                      `json:"times"`
}

type QuestionResponse struct {
	ID                 uint                     `json:"id"`
	QuestionStatements []QuestionStatementResponse `json:"questionstatements"`
	Answer             string                   `json:"answer"`
	
}

type QuestionStatementResponse struct {
	ID               uint     `json:"id"`
	QuestionStatement string   `json:"questionStatement"`
	Choices          []string `json:"choices"`
	Explain          string   `json:"explain"`
}