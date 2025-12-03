package models

import "time"

type UserDeviceToken struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"not null;index"`
    Token     string    `gorm:"not null"`
    CreatedAt time.Time
}
