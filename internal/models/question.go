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
