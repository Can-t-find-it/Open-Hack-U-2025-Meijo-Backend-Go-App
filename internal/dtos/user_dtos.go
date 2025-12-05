package dtos

type UserStatusResponse struct {
	ID     string     `json:"id"`
	Name   string     `json:"name"`
	Status UserStatus `json:"status"`
}

type UserStatus struct {
	TextbookCount int64 `json:"textbookCount"`
	StreakDays    int   `json:"streakDays"`
	FriendCount   int64 `json:"friendCount"`
}