package models

import "gorm.io/gorm"

// ペナルティ情報を管理するモデル
type Penalty struct {
	gorm.Model
	UserID      uint   // 関連するユーザー
	Description string // ペナルティ内容
	Active      bool   // 発動中かどうか
}
