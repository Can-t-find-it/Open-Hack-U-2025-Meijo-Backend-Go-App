package models

import "gorm.io/gorm"

// 問題情報を管理するモデル
type Question struct {
	gorm.Model
	Content  string // 問題文
	Answer   string // 答え
	Category string // 分野
}
