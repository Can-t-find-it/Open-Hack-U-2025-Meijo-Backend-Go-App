/*
package models

import "gorm.io/gorm"

// 問題情報を管理するモデル
type Question struct {
    gorm.Model
    Pattern     string   `json:"pattern"`
    Question    string   `json:"question"`
    OptionsJSON string   `json:"options_json"` // 四択用 JSON文字列
    Explanation string   `json:"explanation"`
}
*/

package models

import (
	"time"
)

// Question: 1つの問題単位（親）
// 以前の Content string は削除し、子テーブル(QuestionStatement)で管理します
type Question struct {
	ID                 uint                `gorm:"primaryKey" json:"id"`
	TextbookID         uint                `json:"textbookId"`
	Textbook           Textbook            `gorm:"foreignKey:TextbookID" json:"textbook,omitempty"`
	
	Answer             string              `json:"answer"` // 正解 (例: "2", "ア")
	
	// 1つの問題に対して複数の「聞き方・選択肢」を持つ
	QuestionStatements []QuestionStatement `gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE" json:"questionStatements"`

	CreatedAt          time.Time           `json:"createdAt"`
	UpdatedAt          time.Time           `json:"updatedAt"`
}

// QuestionStatement: 実際の問題文・選択肢・解説（子）
type QuestionStatement struct {
	ID         uint     `gorm:"primaryKey" json:"id"`
	QuestionID uint     `json:"questionId"`
	
	Statement  string   `json:"questionStatement"` // 問題文: "1+1は？"
	Explain    string   `json:"explain"`           // 解説
	
	// PostgreSQLの場合、GORMのserializer機能でJSONとして保存されます
	Choices    []string `gorm:"serializer:json" json:"choices"` 
	
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}