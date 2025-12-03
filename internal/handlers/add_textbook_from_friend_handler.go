package handlers

import (
	"net/http"

	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/service"

	"github.com/gin-gonic/gin"
)

// AddTextbookHandler : 教科書コピー機能の担当者
type AddTextbookHandler struct {
	service service.TextbookCopyService
}

// NewAddTextbookHandler : コンストラクタ
func NewAddTextbookHandler() *AddTextbookHandler {
	return &AddTextbookHandler{
		service: service.TextbookCopyService{},
	}
}

// ImportTextbook : フレンドの教科書をコピーして取り込む
func (h *AddTextbookHandler) ImportTextbook(c *gin.Context) {
	// 1. 【本人確認】 自分のIDを取得
	userID := c.GetUint("userID")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ログインしてください"})
		return
	}

	// 2. 【注文確認】 JSONボディを受け取る
	// (コピー元の教科書IDと、保存先のフォルダIDが入っています)
	var input dtos.InputAddTextbook
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "入力データが不正です: " + err.Error()})
		return
	}

	// 3. 【厨房へ伝達】 Serviceを呼ぶ
	// DTOの中身を分解して、Serviceの引数に合わせて渡します
	err := h.service.AddTextbook(userID, input.TargetFolderID, input.TextbookID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 4. 【提供】 成功レスポンス (201 Created)
	// c.JSON(http.StatusCreated, gin.H{"message": "教科書をインポートしました"})
	c.Status(http.StatusNoContent)
}
