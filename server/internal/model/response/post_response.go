package response

import (
	"time"
)

type CreatePostResponse struct {
	PostID    	int       		`json:"postid"`
	UserID    	int       		`json:"userid"`
	User 		UserSummary 	`json:"user"`
	Content   	string    		`json:"content"`
	Mood      	string    		`json:"mood"`
	CreatedAt 	time.Time 		`json:"createdat"`
}

type UserSummary struct {
	UserID    	int       		`json:"userid"`
	Username  	string    		`json:"username"`
	FullName  	string    		`json:"fullname"`
}