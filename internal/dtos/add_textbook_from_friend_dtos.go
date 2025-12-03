package dtos

type InputAddTextbook struct {
	TextbookID     uint `json:"friendTextbookId"`
	TargetFolderID uint `json:"folderId"`
}
