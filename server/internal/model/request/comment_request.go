package request

type CreateCommentRequest struct {
	PostID  int    `json:"postid" validate:"required"`
	UserID  int    `json:"userid" validate:"required"`
	Content string `json:"content" validate:"required,min=1,max=500"`
}