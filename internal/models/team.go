package models

import "gorm.io/gorm"

// チーム情報を管理するモデル
type Team struct {
	gorm.Model
	Name    string // チーム名
	Members []User `gorm:"many2many:team_users;"` // 多対多リレーション
}
