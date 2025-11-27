package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// プロジェクトのルートが 'hacku_2025_meijo/backend' であることを前提
	"hacku_2025_meijo/internal/handlers"
	"hacku_2025_meijo/internal/middleware"
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

	authHandler := handlers.NewAuthHandler()

	// --- APIルーティングの定義 ---
	// ここに問題生成に関するすべてのエンドポイントを追加します
	api := r.Group("/api")
	{
		// ログイン機能
		api.POST("/login", authHandler.Login)

		//サインアップ機能
		api.POST("/signup", authHandler.SignUp)

		// ログインしないと使えないAPIエリア
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// // 1. 単一問題生成 (GET)
			// protected.GET("/generate_4choice", handlers.GenerateQuestion4ChoiceHandler)

			// // 2. 単一問題生成 (POST) - 四択
			// protected.POST("/generate_question_4choice_api", handlers.GenerateQuestion4ChoiceAPIHandler)

			// // 3. 複数問題生成 (POST) - 一問一答/穴埋め
			// protected.POST("/generate_workbook_for_q_and_a", handlers.GenerateWorkbookForQAndAHandler)

			// // 4. 複数問題生成 (POST) - 四択
			// protected.POST("/generate_4_choice_workbook_for_q_and_a", handlers.Generate4ChoiceWorkbookForQAndAHandler)

			// 5. 統合問題生成 (POST)
			protected.POST("/generate_problem", handlers.GenerateProblemHandler)
			protected.POST("/question/", handlers.GenerateProblemHandler)
			protected.DELETE("/question/:id", handlers.DeleteQuestionHandler)
			protected.POST("/friend/studylog", handlers.NewFriendStudyLogHandler().GetFriendLog)
		}
	}

	return r
}
