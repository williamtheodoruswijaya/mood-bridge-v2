package response

type CreatePostResponse struct {
	PostID    int    `json:"postid"`
	UserID    int    `json:"userid"`
	Content   string `json:"content"`
	Mood      string `json:"mood"`
	CreatedAt string `json:"createdat"`
}