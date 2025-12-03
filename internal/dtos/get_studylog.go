package dtos

type StudyLogResponse struct {
    ID           string  `json:"id"`
    UserName     string  `json:"userName"`
    DateTime     string  `json:"dateTime"`
    TextbookName string  `json:"textbookName"`
    Accuracy     float64 `json:"accuracy"`
}
