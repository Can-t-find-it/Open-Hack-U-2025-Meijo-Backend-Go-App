package models

import (
	"time" // 時刻型を使う場合
)

type StudyLog struct {
	ID            string `gorm:"primaryKey"`
	UserID        string
	QuestionID    string
	Answered      bool
	AnsweredAt    time.Time
	Score         int
	FriendName    string
	TextbookName  string
	Accuracy      float64
	TodayProgress uint
}
