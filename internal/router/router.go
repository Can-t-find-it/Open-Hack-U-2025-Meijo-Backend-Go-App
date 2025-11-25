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
	api := r.Group("/api")
	{
		api.POST("/question/", handlers.GenerateProblemHandler)
		api.DELETE("/question/:id", handlers.DeleteQuestionHandler)
		api.POST("/getfriends", handlers.NewFriendGetAllHandler().GetAllFriends)
	}

	return r
}
