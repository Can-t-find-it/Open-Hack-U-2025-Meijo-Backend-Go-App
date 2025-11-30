package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"hacku_2025_meijo/internal/service"
)

// POST /api/upload_pdf
func UploadPDFHandler(c *gin.Context) {
	// 1. フォームからファイルを取得 ("file" というキーで送られてくる想定)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// 2. PDFからテキスト抽出
	text, err := service.ExtractTextFromPDF(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF: " + err.Error()})
		return
	}

	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "PDF content is empty or unreadable"})
		return
	}

	// 3. AIで重要単語を抽出
	keywords, err := service.ExtractKeywordsFromText(text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI extraction failed: " + err.Error()})
		return
	}

	// 4. 結果を返す
	// フロントエンドはこの "keywords" をリスト表示し、ユーザーに選ばせる
	c.JSON(http.StatusOK, gin.H{
		"keywords": keywords,
	})
}