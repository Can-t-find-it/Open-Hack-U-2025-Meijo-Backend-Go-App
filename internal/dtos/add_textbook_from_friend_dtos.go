package dtos

type InputAddTextbook struct {
	TextbookID     string `json:"friendTextbookId"`
	TargetFolderID string `json:"folderId"`
}
