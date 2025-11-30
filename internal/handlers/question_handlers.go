package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"
	"hacku_2025_meijo/internal/models"
	"hacku_2025_meijo/internal/database"
)

// ==========================================
// 共通: CORS ミドルウェア
// ==========================================
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}

// ==========================================
// 1. AI生成機能
// ==========================================

// ---- GET 四択 ----
func GenerateQuestion4ChoiceHandler(c *gin.Context) {
	word := c.Query("word")

	if word == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "word is required"})
		return
	}

	result, err := service.GenerateSingle4ChoiceQuestion(word, "1問1答", []string{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"word":        word,
		"question":    result.Question,
		"options":     result.Options,
		"answer":      word,
		"explanation": result.Explanation,
	})
}

// ---- POST 問題生成・保存 (教科書の設定に従う版) ----
func GenerateProblemHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// 1. 教科書IDのチェック
	if body.TextbookID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "textbookId is required"})
		return
	}

	// 2. DBから教科書の情報を取得して、Type（形式）を確認する
	var textbook models.Textbook
	// modelsパッケージのインポートが必要です
	if err := database.DB.First(&textbook, body.TextbookID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Textbook not found"})
		return
	}

	// 教科書の設定を「パターン」として取得
	pattern := string(textbook.Type) 

	var resultItems []dtos.ResultItem
	var targetAnswers []string
	var err error

	// 3. 単語リストの整理
	if len(body.Answers) > 0 {
		targetAnswers = body.Answers
	} else if body.Answer != "" {
		targetAnswers = []string{body.Answer}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "answer or answers is required"})
		return
	}

	// 4. 生成ロジック (教科書のタイプで分岐)
	// 定数 (models.Type...) を使って判定します
	switch models.TextbookType(pattern) {
	
	case models.Type4Choice, models.TypeFillIn4Choice:
		// --- 4択系 (普通の4択 or 穴埋め4択) ---
		// Generate4ChoiceWorkbookForQAndA は pattern 文字列をそのままAIへの指示に使います
		resultItems, err = service.Generate4ChoiceWorkbookForQAndA(targetAnswers, pattern)

	case models.TypeFillIn, models.TypeInput:
		// --- 入力系 (穴埋め入力 or 完全入力) ---
		// GenerateWorkbookForQAndA は pattern 文字列をAIへの指示に使います
		resultItems, err = service.GenerateWorkbookForQAndA(targetAnswers, pattern)

	default:
		// 未知のタイプの場合
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported textbook type: " + pattern})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI Generation Error: " + err.Error()})
		return
	}

	if len(resultItems) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成された問題が0件でした"})
		return
	}

	// 5. DB保存
	var finalResponse []dtos.ResultItem

	for i, item := range resultItems {
		
		currentAnswer := ""
		if i < len(targetAnswers) {
			currentAnswer = targetAnswers[i]
		}


		id, err := service.SaveQuestionToDB(body.TextbookID, item, currentAnswer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB保存エラー: " + err.Error()})
			return
		}
		item.ID = id // 生成されたIDをセット
		finalResponse = append(finalResponse, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"questions": finalResponse,
	})
}

// ---- POST 四択 複数 (個別API) ----
func Generate4ChoiceWorkbookForQAndAHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if len(body.Answers) == 0 || body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "`解答` list and `pattern` are required"})
		return
	}

	resultItems, err := service.Generate4ChoiceWorkbookForQAndA(body.Answers, body.Pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultItems)
}

// ---- POST 四択 単体 (個別API) ----
func GenerateQuestion4ChoiceAPIHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if body.Answer == "" || body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "解答 and pattern are required"})
		return
	}

	result, err := service.GenerateSingle4ChoiceQuestion(body.Answer, body.Pattern, body.ExistingQuestions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, []dtos.ResultItem{result})
}

// ---- POST 一問一答・穴埋め 複数 (個別API) ----
func GenerateWorkbookForQAndAHandler(c *gin.Context) {
	var body dtos.RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	if len(body.Answers) == 0 || body.Pattern == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "`解答` list and `pattern` are required"})
		return
	}

	resultItems, err := service.GenerateWorkbookForQAndA(body.Answers, body.Pattern)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultItems)
}

// ==========================================
// 2. 教科書・DB管理機能
// ==========================================

// GET /api/textbooks - 自分の問題集一覧取得
func GetTextbooksHandler(c *gin.Context) {
	// 認証機能を入れたので、本来はここ（Context）からユーザーIDを取ります
	// 今はテスト用に 1 固定、もしくは middleware から取得する形
	// userID := c.GetUint("userId") // ミドルウェア実装次第
	userID := uint(1) // 仮置き

	result, err := service.GetUserTextbooks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ★ここで JSON の形を決めています
	// { "folder": [ ... ] } という形になります
	c.JSON(http.StatusOK, gin.H{
		"folder": result,
	})
}

// POST /api/textbooks - 教科書作成
func CreateTextbookHandler(c *gin.Context) {
	var req struct {
		Name     string `json:"name"`
		Type     string `json:"type"`
		FolderID uint   `json:"folderId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.CreateTextbook(req.Name, req.Type, req.FolderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Textbook created"})
}

// GET /api/textbook/:id - 教科書詳細取得
func GetTextbookDetailHandler(c *gin.Context) {
	textbookID := c.Param("id")
	result, err := service.GetTextbookDetail(textbookID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Textbook not found"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DELETE /api/textbooks/:id - 教科書削除
func DeleteTextbookHandler(c *gin.Context) {
	textbookID := c.Param("id")
	err := service.DeleteTextbook(textbookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Textbook deleted"})
}

// POST /api/textbook_result - 学習結果保存
func UpdateTextbookResultHandler(c *gin.Context) {
	var req struct {
		TextbookID uint    `json:"textbookId"`
		Score      float64 `json:"score"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	idStr := strconv.Itoa(int(req.TextbookID))
	err := service.UpdateTextbookStatus(idStr, req.Score)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Score updated successfully"})
}

// POST /api/question - 生成した問題をDBに保存 (単体)
func AddQuestionHandler(c *gin.Context) {
	var req struct {
		TextbookID uint            `json:"textbookId"`
		Answer     string          `json:"answer"`
		ResultItem dtos.ResultItem `json:"resultItem"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	err := service.AddQuestionToTextbook(req.TextbookID, req.ResultItem, req.Answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Question added successfully"})
}

// ==========================================
// 3. 詳細な編集機能 (削除・追加)
// ==========================================

// DELETE /api/question/:id - 問題（親）ごとの削除
func DeleteQuestionHandler(c *gin.Context) {
	questionID := c.Param("id")
	err := service.DeleteQuestion(questionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Question deleted"})
}

// POST /api/questionstatements - 既存の問題に新しい聞き方を追加
func AddQuestionStatementHandler(c *gin.Context) {
	var req struct {
		QuestionID uint     `json:"questionId"`
		Statement  string   `json:"questionStatement"`
		Explain    string   `json:"explain"`
		Choices    []string `json:"choices"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	err := service.AddQuestionStatement(req.QuestionID, req.Statement, req.Explain, req.Choices)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Statement added successfully"})
}

// DELETE /api/questionstatements/:id - 問題文（子）単体の削除
func DeleteQuestionStatementHandler(c *gin.Context) {
	statementID := c.Param("id")
	err := service.DeleteQuestionStatement(statementID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Statement deleted"})
}

// GET /api/word - 覚えたい単語提案 (AI版)
func SuggestWordHandler(c *gin.Context) {

	textbookID := c.Query("textbook_id")

	var suggestedWords []string // 配列に変更
	var err error

	if textbookID != "" {
		// ★AI提案ルート
		currentWords, err := service.GetWordsInTextbook(textbookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch textbook words"})
			return
		}

		if len(currentWords) == 0 {
			// 教科書が空ならランダムに3つ返す
			suggestedWords, err = service.GetRandomSuggestedWords(3)
		} else {
			// 教科書があるならAIに3つ考えさせる
			suggestedWords, err = service.SuggestNewWordsViaAI(currentWords)
		}

	} else {
		// ★完全ランダムルート（教科書指定なし）
		suggestedWords, err = service.GetRandomSuggestedWords(3)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Suggestion failed: " + err.Error()})
		return
	}

	// 配列で返す { "words": ["A", "B", "C"] }
	c.JSON(http.StatusOK, gin.H{
		"words": suggestedWords,
	})
}
// POST /api/generate_statement
func GenerateAndAddStatementHandler(c *gin.Context) {
	var req struct {
		QuestionID uint `json:"questionId"`
		// Pattern string `json:"pattern"`  <-- これはもう不要！削除！
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// 引数が QuestionID だけになりました
	newStmt, err := service.GenerateAndAddStatement(req.QuestionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Generation failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "New variation generated successfully",
		"statement": newStmt,
	})
}

// POST /api/folders - フォルダ作成
func CreateFolderHandler(c *gin.Context) {
	// リクエストの受け皿
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// 本来はトークンからユーザーIDを取得しますが、今はテスト用ID(1)を使います
	// (AuthMiddlewareで取得したIDを使う実装に変えるときはここを修正)
	userID := uint(1)

	// Serviceを呼ぶ
	folder, err := service.CreateFolder(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Folder created successfully",
		"folder":  folder,
	})
}

// DELETE /api/folders/:id - フォルダ削除
func DeleteFolderHandler(c *gin.Context) {
	folderID := c.Param("id")
	
	err := service.DeleteFolder(folderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}