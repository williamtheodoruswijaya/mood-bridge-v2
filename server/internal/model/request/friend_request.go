package request

type FriendRequest struct {
	UserID       int `json:"userid" validate:"required"`       // orang yang menambahkan teman
	FriendUserID int `json:"frienduserid" validate:"required"` // orang yang jadi temannya
}