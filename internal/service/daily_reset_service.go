package service

import (
	"fmt"
	"hacku_2025_meijo/internal/database"
	"hacku_2025_meijo/internal/models"
)

func (s *GetTextbooksService) ResetTodayProgress() error {
	result := database.DB.Model(&models.Textbook{}).Where("today_progress > 0").Update("today_progress",0)

	if result.Error != nil {
		return result.Error
	}
	fmt.Printf("TodayProgressリセット完了")
	return nil
}