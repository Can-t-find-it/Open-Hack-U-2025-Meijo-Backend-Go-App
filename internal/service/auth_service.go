package service

import (
	"errors"
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/dtos"
	"hacku_2025_meijo/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct{}

// Login機能
func (s *AuthService) Login(input dtos.LoginInput) (*dtos.LoginResponse, error) {
	var user models.User

	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("ユーザーが見つかりません")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		return nil, errors.New("パスワードが間違っています")
	}

	token, err := GenerateToken(user.ID, user.Name)

	if err != nil {
		return nil, errors.New("エラーによりトークンが発行されませんでした")
	}

	return &dtos.LoginResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			IconURL: user.IconURL,
		},
	}, nil

}

// Signup機能
func (s *AuthService) Signup(input dtos.UserSignup) (*dtos.LoginResponse, error) {
	hashpassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: string(hashpassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	token, err := GenerateToken(user.ID, user.Name)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponse{
		Token: token,
		User: dtos.UserResponse{
			ID:      user.ID,
			Name:    user.Name,
			Email:   user.Email,
			IconURL: user.IconURL,
		},
	}, nil
}

func GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), //トークンの有効期限，発行後２４時間
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		secret = "develop_secret_key"
	}

	return token.SignedString([]byte(secret))
}
