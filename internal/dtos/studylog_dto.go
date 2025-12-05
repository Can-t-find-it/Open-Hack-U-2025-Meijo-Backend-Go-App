package dtos

// import (
//     "time"
// )

// StudyLogResponse は、勉強ログデータをクライアントに返すための構造体です。
type StudyLogResponse struct {
    ID             string  `json:"id"`
    UserName       string  `json:"userName"`
    DateTime       string  `json:"dateTime"`
    TextbookName   string  `json:"textbookName"`
    Accuracy       float64 `json:"accuracy"`
}

// StudyLogsWrapper は、トップレベルのレスポンス構造
type StudyLogsWrapper struct {
    Logs []StudyLogResponse `json:"logs"`
}
