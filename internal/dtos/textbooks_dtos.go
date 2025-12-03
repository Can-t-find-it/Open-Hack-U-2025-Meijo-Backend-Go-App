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
