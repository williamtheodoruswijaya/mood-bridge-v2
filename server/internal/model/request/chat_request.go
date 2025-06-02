package request

type PrivateMessagePayload struct {
	RecipientID int    `json:"recipientid" binding:"required"`
	Content     string `json:"content" binding:"required,max=1024"`
}

type MarkAsReadPayload struct {
	MessageID int `json:"messageid" binding:"required"`
}