package main

import (
	"fmt"
	"os"

	// ユーザーのプロジェクト構造に基づくインポート
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/router"

	// Ginフレームワークのインポート (routerパッケージで使用されていると仮定)
	// "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func init() {
    godotenv.Load()
}
func main() {
	// --- 1. 環境変数の設定と確認 ---
	
	// OpenAI APIキーの存在チェック (サービス層のコア機能に必須)
	if os.Getenv("OPENAI_API_KEY") == "" {
		// 致命的なエラーとしてプログラムを終了
		fmt.Println("致命的なエラー: OPENAI_API_KEY 環境変数が設定されていません。API機能は動作しません。")
		os.Exit(1)
	}

	// サーバーのポートを設定 (環境変数 "PORT" があればそれを使用し、なければ "8080" をデフォルトとする)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// --- 2. データベース接続 ---
	
	// internal/database パッケージの Connect 関数を呼び出し、データベースとの接続を確立する
	// (この関数内で接続処理やマイグレーションが行われていると仮定)
	database.Connect()

	// --- 3. ルーターのセットアップとサーバー起動 ---
	
	// Ginのデフォルト設定でルーターエンジンを初期化
	// (gin.Default() はLoggerとRecoveryミドルウェアを含む)
	// router.SetupRouter() が *gin.Engine を返すことを前提とする
	r := router.SetupRouter()

	fmt.Printf("=================================================================\n")
	fmt.Printf(" 🚀 サーバーをポート %s で実行中\n", port)
	fmt.Printf(" 🔗 アクセス先: http://localhost:%s\n", port)
	fmt.Printf("=================================================================\n")

	// サーバーを起動
	// r.Run() は ListenAndServe をラップしたものであり、ブロックされます。
	if err := r.Run(":" + port); err != nil {
		// サーバー起動失敗時（例: ポートが既に使用されている場合）のエラー処理
		fmt.Printf("サーバー実行中にエラーが発生しました: %v\n", err)
		os.Exit(1)
	}
}