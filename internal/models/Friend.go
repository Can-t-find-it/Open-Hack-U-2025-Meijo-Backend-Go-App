package models


type Friend struct {
    ID            uint `gorm:"primaryKey"`
    UserID        uint
    FriendUserID  uint
    NotifyEnabled bool
}
