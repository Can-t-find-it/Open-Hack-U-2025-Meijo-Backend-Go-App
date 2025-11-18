package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"size:200;unique;not null"`
    Password  string    `gorm:"size:255;not null"` // ハッシュ化したパスワード
    CreatedAt time.Time
    UpdatedAt time.Time
}
