package entity

import "time"

type MessageStatus string

const (
	StatusSent      MessageStatus = "sent"
	StatusDelivered MessageStatus = "delivered"
	StatusRead      MessageStatus = "read"
	StatusFailed    MessageStatus = "failed"
)

type Message struct {
	ID          int 		`gorm:"primaryKey;autoIncrement" json:"id"`
	SenderID    int 		`json:"senderid"`
	RecipientID int 		`json:"recipientid"`
	Content     string 		`json:"content"`
	Timestamp   time.Time 	`json:"timestamp"`
	Status MessageStatus 	`json:"status,omitempty"`
}

// function ini ibaratnya constructor untuk membuat message baru
func NewMessage(senderID, recipientID int, content string) *Message {
	return &Message{
		SenderID:   senderID,
		RecipientID: recipientID,
		Content:    content,
		Timestamp:  time.Now(),
		Status:     StatusSent,
	}
}