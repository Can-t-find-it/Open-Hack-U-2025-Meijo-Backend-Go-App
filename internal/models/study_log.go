package models

import (
	"time" // 時刻型を使う場合
)

type StudyLog struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint
	TextbookID    uint
	QuestionID    uint
	Answered      bool
	AnsweredAt    time.Time
	Score         float64
	Name          string
	TextbookName  string
	Accuracy      float64
	TodayProgress uint
}
