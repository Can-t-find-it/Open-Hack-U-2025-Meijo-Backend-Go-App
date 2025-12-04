package models

import (
	"time"          // 時刻型を使う場合
)

type StudyMaterial struct {
    ID           string      `gorm:"primaryKey"`
    UserID       string
    FilePath     string
    ExtractedText string
    Summary      string
    CreatedAt    time.Time
}
