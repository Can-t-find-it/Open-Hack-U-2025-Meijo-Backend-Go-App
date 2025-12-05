package dtos

type UserSearchResponse struct {
    ID       string `json:"id"`
    UserName string `json:"userName"`
    IsFriend bool   `json:"isFriend"` 
}