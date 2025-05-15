package response

import "time"

type FriendResponse struct {
	FriendID     int         `json:"id"`           // ID teman
	UserID       int         `json:"userid"`       // ID orang yang menambahkan teman
	FriendUserID int         `json:"frienduserid"` // ID orang yang jadi temannya
	FriendStatus bool        `json:"friendstatus"` // Status teman
	CreatedAt    time.Time   `json:"createdat"`    // Waktu saat teman ditambahkan
	User         UserSummary `json:"user"`         // Data pengguna yang menambahkan teman
}