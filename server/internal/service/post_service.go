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
	Create(ctx context.Context, req request.CreatePostRequest) (*response.CreatePostResponse, error)
	Find(ctx context.Context, postID int) (*response.CreatePostResponse, error)
	FindAll(ctx context.Context) ([]*response.CreatePostResponse, error)
	FindByUserID(ctx context.Context, userID int) ([]*response.CreatePostResponse, error)
	Update(ctx context.Context, postID int, req request.CreatePostRequest) (*response.CreatePostResponse, error)
	Delete(ctx context.Context, postID int) (string, error)
}

type PostServiceImpl struct {
	DB *sql.DB
	PostRepository repository.PostRepository
	UserRepository repository.UserRepository
	MoodService MoodPredictionService
}

func NewPostService(db *sql.DB, postRepository repository.PostRepository, userRepository repository.UserRepository, moodService MoodPredictionService) PostService {
	return &PostServiceImpl {
		DB: db,
		PostRepository: postRepository,
		UserRepository: userRepository,
		MoodService: moodService,
	}
}

func (s *PostServiceImpl) Create(ctx context.Context, req request.CreatePostRequest) (*response.CreatePostResponse, error) {
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Validate if user exists
	user, err := s.UserRepository.FindByID(ctx, s.DB, req.UserID)
	if err != nil || user == nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", req.UserID)
		}
		return nil, err
	}

	// Predict mood from content using MoodPredictionService
	moodInput := request.MoodPredictionRequest{
		Input: req.Content,
	}
	moodResp, err := s.MoodService.PredictMood(ctx, moodInput)
	if err != nil {
		return nil, fmt.Errorf("failed to predict mood: %v", err)
	}

	post := entity.Post{
		UserID: req.UserID,
		Content: req.Content,
		Mood: moodResp.Prediction,
		CreatedAt: time.Now(),
	}

	if err := utils.ValidatePostInput(post.Content, post.Mood); err != nil {
		return nil, err
	}

	// Create Post (store ke dalam database)
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

func (s *PostServiceImpl) FindAll(ctx context.Context) ([]*response.CreatePostResponse, error) {
	posts, err := s.PostRepository.FindAll(ctx, s.DB)
	if err != nil {
		if err == sql.ErrNoRows || posts == nil {
			return nil, fmt.Errorf("no posts found")
		}
	}

	var postResponses []*response.CreatePostResponse
	for _, post := range posts {
		user, err := s.UserRepository.FindByID(ctx, s.DB, post.UserID)
		if err != nil {
			if err == sql.ErrNoRows || user == nil {
				return nil, fmt.Errorf("user with ID %d not found", post.UserID)
			}
			return nil, err
		}

		postResponse := &response.CreatePostResponse{
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
		postResponses = append(postResponses, postResponse)
	}

	// Check if any posts were found
	if len(postResponses) == 0 {
		return nil, fmt.Errorf("no posts found")
	}
	return postResponses, nil
}

func (s *PostServiceImpl) FindByUserID(ctx context.Context, userID int) ([]*response.CreatePostResponse, error) {
	// step 1: validate if user exists
	user, err := s.UserRepository.FindByID(ctx, s.DB, userID)
	if err != nil || user == nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", userID)
		}
		return nil, err
	}

	// step 2: find posts by userID
	posts, err := s.PostRepository.FindByUserID(ctx, s.DB, userID)
	if err != nil {
		if err == sql.ErrNoRows || posts == nil {
			return nil, fmt.Errorf("no posts found for user with ID %d", userID)
		}
		return nil, err
	}

	var postResponses []*response.CreatePostResponse
	for _, post := range posts {
		user, err := s.UserRepository.FindByID(ctx, s.DB, post.UserID)
		if err != nil {
			if err == sql.ErrNoRows || user == nil {
				return nil, fmt.Errorf("user with ID %d not found", post.UserID)
			}
			return nil, err
		}
		postResponse := &response.CreatePostResponse{
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
		postResponses = append(postResponses, postResponse)
	}

	// Check if any posts were found
	if len(postResponses) == 0 {
		return nil, fmt.Errorf("no posts found for user with ID %d", userID)
	}

	return postResponses, nil
}

func (s *PostServiceImpl) Update(ctx context.Context, postID int, req request.CreatePostRequest) (*response.CreatePostResponse, error) {
	// Validate if post exists
	post, err := s.PostRepository.Find(ctx, s.DB, postID)
	if err != nil {
		if err == sql.ErrNoRows || post == nil {
			return nil, fmt.Errorf("post with ID %d not found", postID)
		}
		return nil, err
	}

	// Validate if user exists
	user, err := s.UserRepository.FindByID(ctx, s.DB, post.UserID)
	if err != nil {
		if err == sql.ErrNoRows || user == nil {
			return nil, fmt.Errorf("user with ID %d not found", post.UserID)
		}
		return nil, err
	}

	// Ambil data mood dari MoodPredictionService
	moodPrediction := request.MoodPredictionRequest{
		Input: req.Content,
	}
	moodResp, err := s.MoodService.PredictMood(ctx, moodPrediction)
	if err != nil {
		return nil, fmt.Errorf("failed to predict mood: %v", err)
	}

	// Validate post content and mood
	if err := utils.ValidatePostInput(req.Content, moodResp.Prediction); err != nil {
		return nil, err
	}

	// Update post
	post.Content = req.Content
	post.Mood = moodResp.Prediction
	post.CreatedAt = time.Now()

	updatedPost, err := s.PostRepository.Update(ctx, s.DB, postID, post)
	if err != nil {
		return nil, err
	}

	// Convert to response
	postResponse := &response.CreatePostResponse{
		PostID:    updatedPost.PostID,
		UserID:    updatedPost.UserID,
		User: response.UserSummary{
			UserID: user.ID,
			Username: user.Username,
			FullName: user.Fullname,
		},
		Content:   updatedPost.Content,
		Mood:      updatedPost.Mood,
		CreatedAt:  updatedPost.CreatedAt,
	}

	return postResponse, nil
}

func (s *PostServiceImpl) Delete(ctx context.Context, postID int) (string, error) {
	// Validate if post exists
	post, err := s.PostRepository.Find(ctx, s.DB, postID)
	if err != nil {
		if err == sql.ErrNoRows || post == nil {
			return "", fmt.Errorf("post with ID %d not found", postID)
		}
		return "", err
	}

	// Delete post
	message, err := s.PostRepository.Delete(ctx, s.DB, postID)
	if err != nil {
		return "", err
	}

	return message, nil
}