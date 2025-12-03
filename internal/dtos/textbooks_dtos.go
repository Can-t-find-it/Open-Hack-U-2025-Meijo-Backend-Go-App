package dtos

import "time"

type SoloTextbook struct {
	UserID    uint      `json:"user_id"`
	Name      string    `json:"textbook_name"`
	FolderID  uint      `json:"folder_id"`
	Type      string    `json:"type"`
	PlayTimes int       `json:"play_counts"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"update_at"`
}

type ResponseTextbooks struct {
	UserId    uint           `json:"user_id"`
	Textbooks []SoloTextbook `json:"Textbooks"`
}
