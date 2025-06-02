package response

import (
	"mood-bridge-v2/server/internal/entity"
	"time"
)

type ChatMessage struct {
	ID          int    `json:"id"`
	SenderID    int    `json:"senderid"`
	RecipientID int    `json:"recipientid"`
	Content     string    `json:"content"`
	Timestamp   time.Time `json:"timestamp"`
	Status entity.MessageStatus `json:"status,omitempty"`
}

type WebSocketMessage struct {
	Type string 	 `json:"type"`
	Payload interface{} `json:"payload"` // ini payload isinya bisa berupa ChatMessage, ErrorMessage, dll
}

type ErrorMessage struct {
	Code string `json:"code"`
	Message string `json:"message"`
}