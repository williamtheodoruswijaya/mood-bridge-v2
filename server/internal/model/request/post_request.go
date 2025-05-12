package request

type CreatePostRequest struct {
	UserID  int    `json:"userid" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=500"`
	Mood    string `json:"mood" validate:"required,min=1,max=50"`
}