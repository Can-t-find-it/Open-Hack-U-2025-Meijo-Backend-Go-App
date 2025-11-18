package models

import (
	"time"          // 時刻型を使う場合
)

type IconHistory struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint
    ImageURL  string
    CreatedAt time.Time
}
