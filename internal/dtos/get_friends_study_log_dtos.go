package dtos

import "time"

type InputFriendsStudyLog struct {
	FriendsID []uint `json:"friends_id"`
}

type ResponseFriendStudyLog struct {
	FriendID   uint      `json:"friend_id"`
	QuestionID uint      `json:"question_id"`
	Answered   bool      `json:"answered"`
	AnsweredAt time.Time `json:"answered_at"`
	Score      int       `json:"score"`
}
