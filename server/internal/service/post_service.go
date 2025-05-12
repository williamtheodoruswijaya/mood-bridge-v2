package service

import (
	"database/sql"
	"mood-bridge-v2/server/internal/repository"
)

type PostService interface {
	// Create(ctx context.Context, request request.CreatePostRequest) (*response.CreatePostResponse, error)
}

type PostServiceImpl struct {
	DB *sql.DB
	PostRepository repository.PostRepository
}

func NewPostService(db *sql.DB, postRepository repository.PostRepository) PostService {
	return &PostServiceImpl {
		DB: db,
		PostRepository: postRepository,
	}
}

// func (s *PostServiceImpl) Create(ctx context.Context, request request.CreatePostRequest) (*response.CreatePostResponse, error) {
// 	tx, err := s.DB.Begin()
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	post := entity.Post{
// 		UserID: request.UserID,
// 		Content: request.Content,
// 		Mood: request.Mood,
// 		CreatedAt: time.Now(),
// 	}

// 	if err := utils.ValidatePostInput(post.Content, post.Mood); err != nil {
// 		return nil, err
// 	}

// 	// Validate if user exists
// }





