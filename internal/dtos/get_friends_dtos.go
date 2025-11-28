package dtos

type SoloGetFriend struct {
	FriendsID int    `json:"id"`
	Name      string `json:"name"`
}
type GetFriends struct {
	Friends []SoloGetFriend `json:"friends"`
}
