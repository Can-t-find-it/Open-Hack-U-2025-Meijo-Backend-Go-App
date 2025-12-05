package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friend struct {
	ID            string `gorm:"primaryKey;type:varchar(36)"`
	UserID        string `gorm:"type:varchar(36)"`
	FriendUserID  string `gorm:"type:varchar(36)"`
	Name          string
	NotifyEnabled bool
}

func (f *Friend) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return
}