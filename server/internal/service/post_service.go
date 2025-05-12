package service

import (
	"context"
	"database/sql"
	"fmt"
	"mood-bridge-v2/server/internal/entity"
	"mood-bridge-v2/server/internal/model/request"
	"mood-bridge-v2/server/internal/model/response"
	"mood-bridge-v2/server/internal/repository"
	"mood-bridge-v2/server/internal/utils"
	"time"
)

type PostService interface {
	Create(ctx context.Context, request request.CreatePostRequest) (*response.CreatePostResponse, error)
	Find(ctx context.Context, postID int) (*response.CreatePostResponse, error)
}

type PostServiceImpl struct {
	DB *sql.DB
	PostRepository repository.PostRepository
	UserRepository repository.UserRepository
}

func NewPostService(db *sql.DB, postRepository repository.PostRepository, userRepository repository.UserRepository) PostService {
	return &PostServiceImpl {
		DB: db,
		PostRepository: postRepository,
		UserRepository: userRepository,
	}
}

func (s *PostServiceImpl) Create(ctx context.Context, request request.CreatePostRequest) (*response.CreatePostResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	post := entity.Post{
		UserID: request.UserID,
		Content: request.Content,
		Mood: request.Mood,
		CreatedAt: time.Now(),
	}

	if err := utils.ValidatePostInput(post.Content, post.Mood); err != nil {
		return nil, err
	}

	// Validate if user exists\
	user, err := s.UserRepository.FindByID(ctx, s.DB, post.UserID)
	if err != nil || user == nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", post.UserID)
		}
		return nil, err
	}

	// Pada bagian ini kita akan hit api streamlit untuk mendapatkan mood dari content

	// Create Post
	createdPost, err := s.PostRepository.Create(ctx, tx, &post)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	result, err := s.Find(ctx, createdPost.PostID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *PostServiceImpl) Find(ctx context.Context, postID int) (*response.CreatePostResponse, error) {
	post, err := s.PostRepository.Find(ctx, s.DB, postID)
	if err != nil {
		if err == sql.ErrNoRows || post == nil {
			return nil, fmt.Errorf("post with ID %d not found", postID)
		}
		return nil, err
	}

	// Load user details
	user, err := s.UserRepository.FindByID(ctx, s.DB, post.UserID)
	if err != nil {
		if err == sql.ErrNoRows || user == nil {
			return nil, fmt.Errorf("user with ID %d not found", post.UserID)
		} 
		return nil, err
	}

	// Convert to response
	postResponse := response.CreatePostResponse{
		PostID:    post.PostID,
		UserID:    post.UserID,
		User: response.UserSummary{
			UserID: user.ID,
			Username: user.Username,
			FullName: user.Fullname,
		},
		Content:   post.Content,
		Mood:      post.Mood,
		CreatedAt:  post.CreatedAt,
	}
	return &postResponse, nil
}





