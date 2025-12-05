package dtos

type InputFriendGetAll struct {
	UserID string `json:"user_id"`
}

type ResponseFriendGetAll struct {
	User string `json:"user"`
	FriendsUserID []string `json:"friends_id"`
}


// FriendSearchResponse は、フレンド検索の結果として返すデータ構造
type FriendSearchResponse struct {
    ID           string `json:"id"`
    FriendUserID string `json:"friendUserID"`
    Name         string `json:"name"`
}