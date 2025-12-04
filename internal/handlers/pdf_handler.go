package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/models"
	"hacku_2025_meijo/internal/service"
)

// POST /api/upload_pdf
// PDFをアップロードして、新しい教科書を作成し、そこに問題を生成して保存する
func UploadPDFHandler(c *gin.Context) {
	// ---------------------------------------------------
	// 1. リクエストパラメータの取得 (教科書作成用)
	// ---------------------------------------------------
	
	// ファイル
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// フォルダID
	folderIDStr := c.PostForm("folder_id")
	if folderIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "folder_id is required"})
		return
	}
	folderID, _ := strconv.Atoi(folderIDStr)

	// 教科書名
	name := c.PostForm("name")
	if name == "" {
		name = "新規PDF問題集" // デフォルト名
	}

	// 問題タイプ (例: "4択問題形式", "穴埋め入力")
	typeStr := c.PostForm("type")
	if typeStr == "" {
		typeStr = string(models.Type4Choice) // デフォルトは4択
	}

	// ---------------------------------------------------
	// 2. 教科書を新規作成する
	// ---------------------------------------------------
	
	// ServiceのCreateTextbookを使ってDBに保存
	err = service.CreateTextbook(name, typeStr, folderIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Textbook creation failed: " + err.Error()})
		return
	}

	// 作成されたばかりの教科書のIDを取得したいので、検索して特定する
	// (CreateTextbookがIDを返さない仕様の場合は、最新の1件を取得するなどの工夫が必要ですが、
	//  今回は簡易的に「今作ったやつ」を取得します)
	var newTextbook models.Textbook
	// フォルダIDと名前で検索して、一番新しいものを取得
	if err := database.DB.Where("folder_id = ? AND name = ?", folderID, name).Order("created_at desc").First(&newTextbook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve created textbook"})
		return
	}
	textbookID := newTextbook.ID

	// ---------------------------------------------------
	// 3. PDF解析 & 単語抽出
	// ---------------------------------------------------
	text, err := service.ExtractTextFromPDF(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF: " + err.Error()})
		return
	}

	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF content is empty"})
		return
	}

	keywords, err := service.ExtractKeywordsFromText(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI extraction failed: " + err.Error()})
		return
	}

	if len(keywords) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No keywords found", "count": 0})
		return
	}

	// ---------------------------------------------------
	// 4. 問題生成 & 保存
	// ---------------------------------------------------
	_, err = service.GenerateAndSaveBatch(textbookID, typeStr, keywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Generation failed: " + err.Error()})
		return
	}

	// ---------------------------------------------------
	// 5. レスポンス 
	// ---------------------------------------------------
	textbookIDStr := textbookID
	resultDTO, err := service.GetTextbookDetail(textbookIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get textbook detail"})
		return
	}

	// 指定されたJSON形式 { "textbook": { ... } } で返す
	c.JSON(http.StatusOK, gin.H{
		"textbook": gin.H{
			"id":         resultDTO.ID,
			"name":       resultDTO.Name,
			"type":       resultDTO.Type,
			"questions":  resultDTO.Questions,
		},
	})
}

//pdfを受けって、重要単語を抽出して返す
func ExtractKeywordsFromTextHandler(c *gin.Context){
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	text, err := service.ExtractTextFromPDF(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF: " + err.Error()})
		return
	}

	keywords, err := service.ExtractKeywordsFromText(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI extraction failed: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"extractWords": keywords,
	})
}