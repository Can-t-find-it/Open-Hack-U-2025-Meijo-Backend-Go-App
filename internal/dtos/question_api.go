package dtos

// --- 1. APIエンドポイントの入出力構造体 (DTOs) ---

type RequestBody struct {
    Answers           []string `json:"words"`             // 複数解答用
    Answer            string   `json:"answer"`              // 単体問題用
    Pattern           string   `json:"pattern"`             // 問題形式
    ExistingQuestions []string `json:"existing_questions"`
    
    TextbookID        string   `json:"textbookId"`         // 保存先の教科書ID// 重複回避用
}

type ResultItem struct {
    ID          string   `json:"id"`
    Question    string   `json:"question"`
    Options     []string `json:"options,omitempty"`       // 四択のみ
    Explanation string   `json:"explanation"`
}

// --- 2. OpenAI API通信用構造体 ---

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type OpenAIChatCompletionRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    Temperature float64   `json:"temperature"`
    MaxTokens   int       `json:"max_tokens"`
}

type Choice struct {
    Message Message `json:"message"`
}

type OpenAIChatCompletionResponse struct {
    Choices []Choice `json:"choices"`
}
