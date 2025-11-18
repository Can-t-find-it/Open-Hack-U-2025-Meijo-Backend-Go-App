package models

import (
	"time"          // 時刻型を使う場合
)

type Team struct {
    ID        uint   `gorm:"primaryKey"`
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
    Users     []User
}
