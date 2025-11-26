package dtos

type InputFriendGetAll struct {
	UserID int `json:"user_id"`
}

type ResponseFriendGetAll struct {
	User int `json:"user"`
	FriendsUserID []int `json:"friends_id"`
}