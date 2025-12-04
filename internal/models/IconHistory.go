package models

import (
	"time"          // 時刻型を使う場合
)

type IconHistory struct {
    ID        string      `gorm:"primaryKey"`
    UserID    string
    ImageURL  string
    CreatedAt time.Time
}
