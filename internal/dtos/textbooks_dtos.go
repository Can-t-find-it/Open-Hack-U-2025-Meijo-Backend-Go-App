package dtos

import "time"

type SoloTextbook struct {
	TextbookId uint      `json:"textbookId"`
	Name       string    `json:"name"`
	PlayTimes  int       `json:"questionCount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"update_at"`
}

type ResponseTextbooks struct {
	FriendId  uint           `json:"friendId"`
	UserName  string         `json:"userName"`
	Textbooks []SoloTextbook `json:"textbooks"`
}

type FinalResponseTextbooks struct {
	Friends []ResponseTextbooks `json:"friends"`
}

type QuestionStatements struct {
	Id                uint     `json:"id"`
	QuestionStatement string   `json:"questionStatement"`
	Choice            []string `json:"choice"`
	Explain           string   `json:"explain"`
}

type ResponseFriendTextbookInformation struct {
	Id        uint                 `json:"id"`
	Name      string               `json:"name"`
	Type      string               `json:"type"`
	Questions []QuestionStatements `json:"questions"`
}

type FinalResponseFriendTextbookInformation struct {
	Textbook []ResponseFriendTextbookInformation `json:"textbook"`
}
