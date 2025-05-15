package entity

import "time"

type Friend struct {
	FriendID     	int       	`gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       	int       	`gorm:"not null" json:"userid"`            // orang yang menambahkan teman
	FriendUserID 	int       	`gorm:"not null" json:"frienduserid"`      // orang yang jadi temannya
	FriendStatus 	bool      	`gorm:"default:false" json:"friendstatus"` // status teman
	CreatedAt    	time.Time 	`gorm:"autoCreateTime" json:"createdat"`
	User       		*User 		`gorm:"foreignKey:FriendUserID;references:ID"`
}

// tar kita bisa nyari temen dengan cara

/*
	db.Preload("User").Preload("FriendUser").Find(&friends)
*/
