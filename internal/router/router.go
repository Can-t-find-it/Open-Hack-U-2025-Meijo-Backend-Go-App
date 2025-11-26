package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// プロジェクトのルートが 'hacku_2025_meijo' であることを前提
	"hacku_2025_meijo/internal/handlers"
	"hacku_2025_meijo/internal/middleware" // ← ご友人が追加したミドルウェア
)

// SetupRouter はアプリケーションのすべてのルーティングを設定します。
// Ginのルーターインスタンスを返します。
func SetupRouter() *gin.Engine {
	// Ginのデフォルト設定で初期化
	r := gin.Default()

	// --- 共通ミドルウェア (CORS設定) ---
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
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
		c.JSON(http.StatusOK, gin.H{"message": "APIサーバーが正常に動作しています"})
	})
	r.GET("/health", handlers.HealthCheck)

	// 認証ハンドラーの初期化（ご友人のコード）
	authHandler := handlers.NewAuthHandler()

	// --- APIルーティングの定義 ---
	api := r.Group("/api")
	{
		// ------------------------------------
		// 1. 誰でもアクセスできる機能 (ログイン・登録)
		// ------------------------------------
		api.POST("/login", authHandler.Login)
		api.POST("/signup", authHandler.SignUp)

		// ------------------------------------
		// 2. ログイン必須の機能 (AuthMiddleware適用)
		// ------------------------------------
		// 以下の機能を使うには、ログイン後に貰えるトークンが必要です
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// === AI問題生成機能 ===
			protected.GET("/generate_4choice", handlers.GenerateQuestion4ChoiceHandler)
			protected.POST("/generate_question_4choice_api", handlers.GenerateQuestion4ChoiceAPIHandler)
			protected.POST("/generate_workbook_for_q_and_a", handlers.GenerateWorkbookForQAndAHandler)
			protected.POST("/generate_4_choice_workbook_for_q_and_a", handlers.Generate4ChoiceWorkbookForQAndAHandler)
			protected.POST("/generate_problem", handlers.GenerateProblemHandler)

			// === 教科書・フォルダ管理機能 (りょうさん作成) ===
			protected.GET("/textbooks", handlers.GetTextbooksHandler)
			protected.POST("/textbooks", handlers.CreateTextbookHandler)
			protected.GET("/textbook/:id", handlers.GetTextbookDetailHandler)
			protected.DELETE("/textbooks/:id", handlers.DeleteTextbookHandler)
			protected.POST("/textbook_result", handlers.UpdateTextbookResultHandler) // 学習結果の保存

			// === 問題・問題文の操作 ===
			protected.POST("/question", handlers.AddQuestionHandler) // 問題保存
			protected.DELETE("/question/:id", handlers.DeleteQuestionHandler)
			protected.POST("/questionstatements", handlers.AddQuestionStatementHandler)
			protected.DELETE("/questionstatements/:id", handlers.DeleteQuestionStatementHandler)

			// === その他 ===
			protected.GET("/word", handlers.SuggestWordHandler) // 単語提案

            
		}
	}

	return r
}