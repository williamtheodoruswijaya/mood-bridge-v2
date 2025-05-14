package response

import "time"

type CreateCommentResponse struct {
	CommentID int         `json:"commentid"`
	PostID    int         `json:"postid"`
	UserID    int         `json:"userid"`
	User      UserSummary `json:"user"`
	Content   string      `json:"content"`
	CreatedAt time.Time   `json:"created_at"`
}