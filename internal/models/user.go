package models

import (
	"time"          // 時刻型を使う場合
)

type User struct {
    ID            uint           `gorm:"primaryKey"`
    Name          string
    Email         string         `gorm:"unique"`
    PasswordHash  string
    IconURL       string
    PenaltyStatus int
    LastActiveAt  time.Time
    TeamID        *uint
    Team          Team           `gorm:"foreignKey:TeamID"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
