package models

import "gorm.io/gorm"

// 学習記録を管理するモデル
type StudyLog struct {
	gorm.Model
	UserID    uint   // 関連ユーザー
	QuestionID uint  // 関連する問題
	Result    bool   // 正解かどうか
	Note      string // メモ
}
