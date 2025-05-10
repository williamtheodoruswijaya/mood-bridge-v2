package entity

import "time"

type Post struct {
	PostID    int       `gorm:"primaryKey;autoIncrement"`
	UserID    int       `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"` // Ini perlu biar bisa db.Preload("User").Find(&posts)
	Content   string    `json:"content"`
	Mood      string    `json:"mood"`
	CreatedAt time.Time `json:"created_at"`
}
