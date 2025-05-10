package entity

type Friend struct {
	FriendID     int `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int `gorm:"not null" json:"user_id"`        // orang yang menambahkan teman
	FriendUserID int `gorm:"not null" json:"friend_user_id"` // orang yang jadi temannya

	User       User `gorm:"foreignKey:UserID;references:ID"`
	FriendUser User `gorm:"foreignKey:FriendUserID;references:ID"`
}

// tar kita bisa nyari temen dengan cara

/*
	db.Preload("User").Preload("FriendUser").Find(&friends)
*/
