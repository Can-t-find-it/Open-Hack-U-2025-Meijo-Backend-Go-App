package models

import "time"

type UserDeviceToken struct {
    ID        string      `gorm:"primaryKey"`
    UserID    string      `gorm:"not null;index"`
    Token     string    `gorm:"not null"`
    CreatedAt time.Time
}
