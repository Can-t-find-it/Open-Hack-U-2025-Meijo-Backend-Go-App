package dtos

type InputFriendGetAll struct {
	UserID string `json:"user_id"`
}

type ResponseFriendGetAll struct {
	User string `json:"user"`
	FriendsUserID []string `json:"friends_id"`
}