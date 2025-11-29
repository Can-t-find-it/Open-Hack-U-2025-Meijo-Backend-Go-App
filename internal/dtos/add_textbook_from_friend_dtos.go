package dtos

type InputAddTextbook struct {
	TextbookID     uint `json:"textbook_id"`
	TargetFolderID uint `json:"target_folder_id"`
}
