package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"hacku_2025_meijo/backend/internal/models" // User などのモデルをここでインポート
)

var DB *gorm.DB

func Connect() {
	// .env を読み込む
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables only")
	}

	// DSN（Data Source Name）を組み立てる
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	// DB 接続
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to initialize database, got error:", err)
	}

	DB = db
	fmt.Println("Connected to database!")

	// 必要なモデルを自動でテーブル作成
	AutoMigrate()
}

// AutoMigrate 関数
func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Question{},
		&models.Penalty{},
		&models.StudyLog{},
		&models.Team{},
	)
	if err != nil {
		log.Fatal("failed to migrate database, got error:", err)
	}
	fmt.Println("Database migration completed!")
}
