package models

import (
	"time"
)
type TextbookType string
// 定数をその型で定義する

// Folder: 科目や資格ごとのフォルダ
// 例: "基本情報技術者試験", "TOEIC"
type Folder struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `json:"userId"`
	Name      string     `json:"name"`
	Progress  int        `json:"progress"` // 進捗率
	
	// Folderを削除したら中身のTextbookも消える設定
	Textbooks []Textbook `gorm:"foreignKey:FolderID;constraint:OnDelete:CASCADE" json:"textbooks"`
	
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// Textbook: 具体的な問題集
// 例: "第1章 基礎理論", "過去問2023秋"
type Textbook struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	FolderID        uint       `json:"folderId"`
	Name            string     `json:"name"`
	Type            TextbookType `json:"type"`            // "4択問題形式" など
	StudyMaterialID *uint      `json:"studyMaterialId"` // 元資料ID(任意)
	
	
	// ↓【追加】仕様変更に対応
	// 過去のスコア履歴 (例: [80.5, 90.0, ...])
	ScoreHistory    []float64  `gorm:"serializer:json" json:"score"`

	// 解いた回数
	PlayTimes       int        `json:"times"`
	
	// Textbookを削除したら中身のQuestionも消える設定
	Questions       []Question `gorm:"foreignKey:TextbookID;constraint:OnDelete:CASCADE" json:"questions"`

	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}