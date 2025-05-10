package entity

import (
	"time"
)

type User struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	Username   string    `gorm:"type:varchar(100);unique;not null" json:"username"`
	Fullname   string    `gorm:"type:varchar(100);not null" json:"fullname"`
	Email      string    `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string    `gorm:"type:varchar(255);not null" json:"-"`
	ProfileUrl string    `gorm:"type:varchar(255)" json:"profile_url"`
	Posts      []Post    `gorm:"foreignKey:UserID"`
	CreatedAt  time.Time `json:"created_at"`
}
