package dtos

import "time"

// type InputFriendsStudyLog struct {
// 	FriendsID []uint `json:"friends_id"`
// }

type FriendStudyLog struct {
	ID            string      `json:"logId"`
	FriendID      string      `json:"friendId"`
	FriendName    string    `json:"friendName"` //新規追加
	AnsweredAt    time.Time `json:"dateTime"`
	TextbookName  string    `json:"textbookName"`  //新規追加
	Accuracy      float64   `json:"accuracy"`      //新規追加
	TodayProgress uint      `json:"todayProgress"` //新規追加

	// QuestionID uint      `json:"question_id"`
	// Answered   bool      `json:"answered"`
	// Score      int       `json:"score"`
}

type AllFriendsStudyLog struct {
	Logs []FriendStudyLog `json:"logs"`
}
