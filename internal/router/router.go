package router

import (
    "net/http"

<<<<<<< HEAD
    "github.com/gin-gonic/gin"
    // プロジェクトのルートが 'hacku_2025_meijo' であることを前提
    "hacku_2025_meijo/internal/handlers"
=======
	"github.com/gin-gonic/gin"
	// プロジェクトのルートが 'hacku_2025_meijo/backend' であることを前提
	"hacku_2025_meijo/internal/handlers"
	"hacku_2025_meijo/internal/middleware"
>>>>>>> 531e805c3e2e6b4a117b402828578c53add3195d
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
        c.JSON(http.StatusOK, gin.H{"message": "APIサーバーが正常に動作しています"})
    })
    r.GET("/health", handlers.HealthCheck)

<<<<<<< HEAD
    // --- APIルーティングの定義 ---
    // ここにすべての機能のエンドポイントをまとめます
    // グループ化しておくと、後で認証ミドルウェアなどを一括適用しやすくなります
    api := r.Group("/api") // フロントエンドに合わせて /api プレフィックスをつけるのが一般的ですが、変更も可能です
    {
        // ==========================================
        // 1. AI問題生成機能 (既存のルート)
        // ==========================================
        
        // 単一問題生成 (GET)
        api.GET("/generate_4choice/", handlers.GenerateQuestion4ChoiceHandler)
        
        // 単一問題生成 (POST) - 四択
        api.POST("/generate_question_4choice_api/", handlers.GenerateQuestion4ChoiceAPIHandler)

        // 複数問題生成 (POST) - 一問一答/穴埋め
        api.POST("/generate_workbook_for_q_and_a/", handlers.GenerateWorkbookForQAndAHandler)
        
        // 複数問題生成 (POST) - 四択
        api.POST("/generate_4_choice_workbook_for_q_and_a/", handlers.Generate4ChoiceWorkbookForQAndAHandler)

        // 統合問題生成 (POST)
        api.POST("/generate_problem/", handlers.GenerateProblemHandler)

        // ==========================================
        // 2. 学習機能・CRUD管理 (今回追加分)
        // ==========================================

        // 教科書・フォルダ一覧取得
        api.GET("/textbooks", handlers.GetTextbooksHandler)

        // 教科書（問題集）の作成
        api.POST("/textbooks", handlers.CreateTextbookHandler)

        // 教科書（問題集）の詳細取得（中の問題リスト含む）
        api.GET("/textbook/:id", handlers.GetTextbookDetailHandler)

        // 教科書（問題集）の削除
        api.DELETE("/textbooks/:id", handlers.DeleteTextbookHandler)


        // 生成した問題を教科書に保存する
        api.POST("/question", handlers.AddQuestionHandler)

		// 問題削除
        api.DELETE("/question/:id", handlers.DeleteQuestionHandler)

        // 問題文（Statement）単体の追加・削除
        api.POST("/questionstatements", handlers.AddQuestionStatementHandler)
        api.DELETE("/questionstatements/:id", handlers.DeleteQuestionStatementHandler)

        // 覚えたい単語提案機能
        api := r.Group("/api")
        {
            // ... 既存のルート ...
            // 6. 単語提案 (GET)
            api.GET("/word", handlers.SuggestWordHandler)
        }
    }

    // 注意: もしフロントエンドが "/api" なしでアクセスしている場合は、
    // api := r.Group("/") に戻すか、クライアント側のURLを修正してください。

    return r
}
=======
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
			// 1. 単一問題生成 (GET)
			protected.GET("/generate_4choice", handlers.GenerateQuestion4ChoiceHandler)

			// 2. 単一問題生成 (POST) - 四択
			protected.POST("/generate_question_4choice_api", handlers.GenerateQuestion4ChoiceAPIHandler)

			// 3. 複数問題生成 (POST) - 一問一答/穴埋め
			protected.POST("/generate_workbook_for_q_and_a", handlers.GenerateWorkbookForQAndAHandler)

			// 4. 複数問題生成 (POST) - 四択
			protected.POST("/generate_4_choice_workbook_for_q_and_a", handlers.Generate4ChoiceWorkbookForQAndAHandler)

			// 5. 統合問題生成 (POST)
			protected.POST("/generate_problem", handlers.GenerateProblemHandler)
      protected.POST("/question/", handlers.GenerateProblemHandler)
      protected.DELETE("/question/:id", handlers.DeleteQuestionHandler)
		}
	}

	return r
}
>>>>>>> 531e805c3e2e6b4a117b402828578c53add3195d
