package dtos

import "time"

type SoloTextbook struct {
	TextbookId string    `json:"textbookId"`
	Name       string    `json:"name"`
	QuestionCount          int       `json:"questionCount"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"update_at"`
}

type ResponseTextbooks struct {
	FriendId  string         `json:"friendId"`
	UserName  string         `json:"userName"`
	Textbooks []SoloTextbook `json:"textbooks"`
}

type FinalResponseTextbooks struct {
	Friends []ResponseTextbooks `json:"friends"`
}


