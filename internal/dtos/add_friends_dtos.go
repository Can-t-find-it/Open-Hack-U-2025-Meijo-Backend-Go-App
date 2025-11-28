package dtos

type AddFriends struct {
	Friends []int `json:"add_friends"`
}

type SoloAddFriend struct {
	Friend int    `json:"friend"`
	Name   string `json:"name"`
}

type ResponseAddFriends struct {
	Friends []SoloAddFriend `json:"added_friends"`
}
