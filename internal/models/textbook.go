package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TextbookType string

// 定数をその型で定義する
const (
	Type4Choice       TextbookType = "4択問題形式"
	TypeFillIn        TextbookType = "穴埋め解答入力形式"
	TypeFillIn4Choice TextbookType = "穴埋め4択問題形式"
	TypeInput         TextbookType = "解答入力形式"
)

// Folder: 科目や資格ごとのフォルダ
// 例: "基本情報技術者試験", "TOEIC"
type Folder struct {
	ID       string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID   string `json:"userId" gorm:"type:varchar(36)"`
	Name     string `json:"name"`
	Progress int    `json:"progress"` // 進捗率

	// Folderを削除したら中身のTextbookも消える設定
	Textbooks []Textbook `gorm:"foreignKey:FolderID;constraint:OnDelete:CASCADE" json:"textbooks"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// BeforeCreate: 保存する前に自動でUUIDを生成する
func (f *Folder) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}
	return
}

// Textbook: 具体的な問題集
// 例: "第1章 基礎理論", "過去問2023秋"
type Textbook struct {
	ID              string       `gorm:"primaryKey;type:varchar(36)" json:"id"`
	FolderID        string       `json:"folderId" gorm:"type:varchar(36)"`
	UserID          string       `json:"userId" gorm:"type:varchar(36)"` // 上野追加部分．この問題集の所有者ID
	Name            string       `json:"name"`
	Type            TextbookType `json:"type"`                                    // "4択問題形式" など
	StudyMaterialID *string      `json:"studyMaterialId" gorm:"type:varchar(36)"` // 元資料ID(任意)

	// ↓【追加】仕様変更に対応
	// 過去のスコア履歴 (例: [80.5, 90.0, ...])
	ScoreHistory []float64 `gorm:"serializer:json" json:"score"`

	// 解いた回数
	PlayTimes int `json:"times"`

	// Textbookを削除したら中身のQuestionも消える設定
	Questions []Question `gorm:"foreignKey:TextbookID;constraint:OnDelete:CASCADE" json:"questions"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
