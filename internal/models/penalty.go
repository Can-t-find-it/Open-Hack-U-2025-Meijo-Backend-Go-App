package models

import (
	"time"          // 時刻型を使う場合
)

type Penalty struct {
    ID         uint      `gorm:"primaryKey"`
    UserID     uint
    Type       string    // icon_lock, spam_notifications, data_erode
    ExecutedAt time.Time
    ExpiresAt  time.Time
}
