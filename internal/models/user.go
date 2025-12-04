package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	// IDを string に変更し、サイズを36文字(UUIDの長さ)に指定
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name         string    `json:"name"`
	Email        string    `gorm:"unique" json:"email"`
	PasswordHash string    `json:"-"` // JSONには出力しない
	IconURL      string    `json:"icon_url"`
	
	// 外部キーの型も string に合わせる
	TeamID       *string   `json:"team_id"` 
	Team         Team      `gorm:"foreignKey:TeamID"`

	Folders      []Folder  `gorm:"foreignKey:UserID"`

	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// BeforeCreate: 保存する前に自動でUUIDを生成する
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}