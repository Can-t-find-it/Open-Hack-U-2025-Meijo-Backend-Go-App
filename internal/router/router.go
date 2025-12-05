package router

import (
	"net/http"

	"hacku_2025_meijo/internal/handlers"
	"hacku_2025_meijo/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter はアプリケーションのすべてのルーティングを設定します。
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

	// --- ハンドラーの初期化 (ここでまとめて行う) ---
	authHandler := handlers.NewAuthHandler()
	friendChangeHandler := handlers.NewAddFriendsHandler()
	friendGetAllHandler := handlers.NewFriendGetAllHandler()
	friendStudyLogHandler := handlers.NewFriendStudyLogHandler()
	// ★追加: 教科書取得ハンドラを初期化
	friendTextbookHandler := handlers.NewGetTextBooksHandler()

	deviceHandler := handlers.NewDeviceHandler()

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
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{

			// 問題文生成

			protected.POST("/generate_statement/:question_id", handlers.GenerateAndAddStatementHandler)
			// 総合問題生成
			protected.POST("/generate_problem/:textbook_id", handlers.GenerateProblemHandler)

			// === 教科書・フォルダ管理機能 ===

			protected.POST("/folders", handlers.CreateFolderHandler)
			protected.DELETE("/folders", handlers.DeleteFoldersBatchHandler)
			protected.GET("/textbooks", handlers.GetTextbooksHandler)//?????
			protected.POST("/textbooks", handlers.CreateTextbookHandler)
			protected.GET("/textbook/:id", handlers.GetTextbookDetailHandler)
			protected.DELETE("/textbook/:id", handlers.DeleteTextbookHandler)
			protected.POST("/textbook/:id/result", handlers.UpdateTextbookResultHandler)

			protected.POST("/textbook/:id/studylog", handlers.NewStudyLogHandler().CreateStudyLogHandler)

			// === 問題・問題文の操作 ===
			protected.POST("/question", handlers.AddQuestionHandler)//?????
			protected.DELETE("/question/:id", handlers.DeleteQuestionHandler)
			protected.POST("/questionstatement", handlers.AddQuestionStatementHandler)//?????
			protected.DELETE("/questionstatement/:id", handlers.DeleteQuestionStatementHandler)

			// === ユーザー情報取得 ===
			protected.GET("/user/status", handlers.GetUserStatus)
			protected.GET("/user/studylog", handlers.GetStudyLogsHandler)

			// === その他 ===
			protected.GET("/textbook/:id/word", handlers.SuggestWordHandler)

			// === フレンド機能 ===
			// フレンド一覧取得
			protected.POST("/getfriends", friendGetAllHandler.GetAllFriends)

			// フレンド追加・削除・取得
			protected.POST("/friend/change/:id", friendChangeHandler.AddFriends)
			protected.DELETE("/friend/change/:id", friendChangeHandler.DeleteFriendsBatch)
			protected.GET("/friend/change", friendChangeHandler.GetFriends)

			protected.GET("/friend/studylog", friendStudyLogHandler.GetFriendLog)

			// フレンド問題集取得
			protected.GET("/friend/textbooks", friendTextbookHandler.GetTextbooks)
			protected.POST("/friend/textbooks", handlers.NewAddTextbookHandler().ImportTextbook)
			protected.GET("/friend/textbook/:textbook_id", friendTextbookHandler.GetFriendTextbookDetail)
			// === ペナルティ通知機能 ===
			protected.POST("/penalty/check", handlers.NewPenaltyHandler().CheckPenalty)
			protected.POST("/device_token", deviceHandler.SaveDeviceToken)

			//=== PDFアップロード機能 ===
			protected.POST("/upload_pdf", handlers.UploadPDFHandler)
			protected.POST("/extract_keywords/:textbook_id", handlers.ExtractKeywordsFromTextHandler)

		}
	}

	return r
}
