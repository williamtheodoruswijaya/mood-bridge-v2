package entity

import "time"

type Comment struct {
	CommentID int       `gorm:"primaryKey;autoIncrement" json:"comment_id"`
	UserID    int       `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"` // Relasi ke User
	PostID    int       `gorm:"not null" json:"post_id"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post"` // Relasi ke Post
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
