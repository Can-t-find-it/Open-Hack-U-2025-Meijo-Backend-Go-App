package models

import (
	"time"          // 時刻型を使う場合
)

type StudyLog struct {
    ID         uint      `gorm:"primaryKey"`
    UserID     uint
    QuestionID uint
    Answered   bool
    AnsweredAt time.Time
    Score      int
}
