package service

import (
    "hacku_2025_meijo/internal/database"
    "hacku_2025_meijo/internal/models"
)

func DeleteFolders(folderIDs []string) error {

    result:= database.DB.Unscoped().Where("id IN ?", folderIDs).Delete(&models.Folder{})
    return result.Error

}
