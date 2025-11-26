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

    // ↓【追加】ユーザーは複数の学習フォルダを持つ
	Folders       []Folder  `gorm:"foreignKey:UserID"`
    
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
