package dtos

type AddFriends struct {
	Friends []string `json:"friends_id"`
}

type SoloAddFriend struct {
	Friend string `json:"friend_id"`
	Name   string `json:"name"`
}

type ResponseAddFriends struct {
	Friends []SoloAddFriend `json:"added_friends"`
}
