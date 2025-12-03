package service

import (
    "hacku_2025_meijo/internal/database"
    "hacku_2025_meijo/internal/models"
	"strconv"
	"strings"
)

func DeleteFolders(folderIDs []string) error {
    uintIDs := make([]uint, 0, len(folderIDs))

    for _, id := range folderIDs {
        // 先頭の "id" を取り除く
        clean := strings.TrimPrefix(id, "id")

        n, err := strconv.ParseUint(clean, 10, 64)
        if err != nil {
            return err
        }

        uintIDs = append(uintIDs, uint(n))
    }

    return database.DB.Where("id IN ?", uintIDs).Delete(&models.Folder{}).Error
}
