package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// プロジェクトのルートが 'hacku_2025_meijo/backend' であることを前提
	"hacku_2025_meijo/internal/handlers"
)

// SetupRouter はアプリケーションのすべてのルーティングを設定します。
// Ginのルーターインスタンスを返します。
func SetupRouter() *gin.Engine {
	// Ginのデフォルト設定で初期化
	r := gin.Default()

	// --- 共通ミドルウェア (CORS設定) ---
	// フロントエンドからのアクセスを許可するために必須
	r.Use(func(c *gin.Context) {
		// すべてのオリジンからのアクセスを許可
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	// --- ルートパス/ヘルスチェック ---
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "問題生成APIが正常に動作しています"})
	})
	// ご要望のヘルスチェックも追加
	r.GET("/health", handlers.HealthCheck)

	// --- APIルーティングの定義 ---
	// ここに問題生成に関するすべてのエンドポイントを追加します
	api := r.Group("/")
	{
		// 1. 単一問題生成 (GET)
		api.GET("/generate_4choice/", handlers.GenerateQuestion4ChoiceHandler)
		
		// 2. 単一問題生成 (POST) - 四択
		api.POST("/generate_question_4choice_api/", handlers.GenerateQuestion4ChoiceAPIHandler)

		// 3. 複数問題生成 (POST) - 一問一答/穴埋め
		api.POST("/generate_workbook_for_q_and_a/", handlers.GenerateWorkbookForQAndAHandler)
		
		// 4. 複数問題生成 (POST) - 四択
		api.POST("/generate_4_choice_workbook_for_q_and_a/", handlers.Generate4ChoiceWorkbookForQAndAHandler)

		// 5. 統合問題生成 (POST)
		api.POST("/generate_problem/", handlers.GenerateProblemHandler)
	}

	return r
}