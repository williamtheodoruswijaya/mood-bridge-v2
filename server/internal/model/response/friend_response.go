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

type FriendRecommendationResponse struct {
	UserID      int    `json:"userid"`
	Username    string `json:"username"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	OverallMood string `json:"overall_mood"`
}