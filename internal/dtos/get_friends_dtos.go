package dtos

type SoloGetFriend struct {
	FriendsID string    `json:"friend_id"`
	Name      string `json:"name"`
}
type GetFriends struct {
	Friends []SoloGetFriend `json:"friends"`
}
