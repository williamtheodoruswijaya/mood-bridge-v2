package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mood-bridge-v2/server/internal/entity"
	"mood-bridge-v2/server/internal/model/request"
	"mood-bridge-v2/server/internal/model/response"
	"mood-bridge-v2/server/internal/repository"
	"mood-bridge-v2/server/internal/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type PostService interface {
	Create(ctx context.Context, req request.CreatePostRequest) (*response.CreatePostResponse, error)
	Find(ctx context.Context, postID int) (*response.CreatePostResponse, error)
	FindAll(ctx context.Context, limit, offset int) ([]*response.CreatePostResponse, error)
	FindByUserID(ctx context.Context, userID int) ([]*response.CreatePostResponse, error)
	Update(ctx context.Context, postID int, req request.CreatePostRequest) (*response.CreatePostResponse, error)
	Delete(ctx context.Context, postID int) (string, error)
	// GetPostBySearch(ctx context.Context, query string) ([]*response.CreatePostResponse, error)
	GetFriendPosts(ctx context.Context, userID int) ([]*response.CreatePostResponse, error)
}

type PostServiceImpl struct {
	DB *sql.DB
	PostRepository repository.PostRepository
	UserRepository repository.UserRepository
	MoodService MoodPredictionService
	RedisClient *redis.Client
}

func NewPostService(db *sql.DB, postRepository repository.PostRepository, userRepository repository.UserRepository, moodService MoodPredictionService, redisClient *redis.Client) PostService {
	return &PostServiceImpl {
		DB: db,
		PostRepository: postRepository,
		UserRepository: userRepository,
		MoodService: moodService,
		RedisClient: redisClient,
	}
}

const cacheVersion = 1

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

	// Invalidate the "post:all" cache after post creation
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("post:all:v%d", cacheVersion)).Err()

	return result, nil
}

func (s *PostServiceImpl) Find(ctx context.Context, postID int) (*response.CreatePostResponse, error) {
	// Step 0: Check if post yang mau kita cari ada di cache
	cacheKey := fmt.Sprintf("post:%d:v%d", postID, cacheVersion)
	cached, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var postResp response.CreatePostResponse
		if err := json.Unmarshal([]byte(cached), &postResp); err == nil {
			return &postResp, nil
		}
	}

	// Start step 1 kalau step 0 gagal: cari post-nya di database
	post, err := s.PostRepository.Find(ctx, s.DB, postID)
	if err != nil {
		if err == sql.ErrNoRows || post == nil {
			return nil, fmt.Errorf("post with ID %d not found", postID)
		}
		return nil, err
	}

	// Step 2: Load user details
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

	// Cache the response
	jsonVal, err := json.Marshal(postResponse)
	if err == nil {
		_ = s.RedisClient.Set(ctx, cacheKey, jsonVal, 10*time.Minute).Err()
	}

	// Return response
	return &postResponse, nil
}

func (s *PostServiceImpl) FindAll(ctx context.Context, limit, offset int) ([]*response.CreatePostResponse, error) {
	// step 0: Check cache-nya dulu (apakah data yang diretrieve ada perubahan atau engga)
	cacheKey := fmt.Sprintf("post:all:v%d:limit:%d:offset:%d", cacheVersion, limit, offset)
	cached, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var postResp []*response.CreatePostResponse
		if err := json.Unmarshal([]byte(cached), &postResp); err == nil {
			return postResp, nil
		}
	}

	// step 1: find all posts (fetch dari database jika tidak ada di cache)
	posts, err := s.PostRepository.FindAll(ctx, s.DB, limit, offset)
	if err != nil {
		if err == sql.ErrNoRows || posts == nil {
			return nil, fmt.Errorf("no posts found")
		}
		return nil, err
	}

	// step 2: kumpulkan semua user id buat di proses sebagai batch queries (menghindari N+1 query)
	userIDMap := make(map[int]bool)
	for _, post := range posts {
		userIDMap[post.UserID] = true
	}

	// step 3: ambil semua id user yang sudah kita kumpulkan (ambil yang True aja)
	var userIDs []int
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	// step 4: query ke database untuk ambil semua user yang ada di userIDMap secara sekaligus
	users, err := s.UserRepository.FindByIDs(ctx, s.DB, userIDs)
	if err != nil {
		if err == sql.ErrNoRows || users == nil {
			return nil, fmt.Errorf("no users found")
		}
		return nil, err
	}

	// step 5: buat map untuk memudahkan pengambilan user berdasarkan ID
	userMap := make(map[int]*entity.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var postResponses []*response.CreatePostResponse
	for _, post := range posts {
		user, ok := userMap[post.UserID]
		if !ok {
			fmt.Printf("user with ID %d not found for post ID %d\n", post.UserID, post.PostID)
			continue // Skip this post if user not found
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

	// Simpan response ke dalam cache
	jsonVal, err := json.Marshal(postResponses) // Convert data ke JSON
	if err == nil {
		_ = s.RedisClient.Set(ctx, cacheKey, jsonVal, 10*time.Minute).Err() // Simpan ke Redis
	}

	// Return response
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

	// Kalau udah aman, baru kita update post-nya
	post.Content = req.Content
	post.Mood = moodResp.Prediction

	// Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	// Rollback kalau terjadi error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Update post di DB
	updatedPost, err := s.PostRepository.Update(ctx, tx, postID, post)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Invalidate the cache for the updated post
	s.RedisClient.Del(ctx, fmt.Sprintf("post:%d:v%d", postID, cacheVersion))
	s.RedisClient.Del(ctx, fmt.Sprintf("post:all:v%d", cacheVersion))

	// Cari post yang udah diupdate untuk direturn sebagai response
	postResponse, err := s.Find(ctx, updatedPost.PostID)
	if err != nil {
		return nil, err
	}

	return postResponse, nil
}

func (s *PostServiceImpl) Delete(ctx context.Context, postID int) (string, error) {
	// Start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return "", err
	}

	// Rollback kalau terjadi error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check if post exists
	_, err = s.PostRepository.Find(ctx, s.DB, postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("post with ID %d not found", postID)
		}
		return "", err
	}

	// Delete from DB
	message, err := s.PostRepository.Delete(ctx, tx, postID)
	if err != nil {
		return "", err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	// Invalidate the cache for the deleted post
	s.RedisClient.Del(ctx, fmt.Sprintf("post:%d:v%d", postID, cacheVersion))
	s.RedisClient.Del(ctx, fmt.Sprintf("post:all:v%d", cacheVersion))

	return message, nil
}

// func (s *PostServiceImpl) GetPostBySearch(ctx context.Context, query string) ([]*response.CreatePostResponse, error) {
// 	// nanti diisinya besok.
// }

func (s *PostServiceImpl) GetFriendPosts(ctx context.Context, userID int) ([]*response.CreatePostResponse, error) {
	cacheKey := fmt.Sprintf("friend_posts:%d:v%d", userID, cacheVersion)
	cached, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var postResp []*response.CreatePostResponse
		if err := json.Unmarshal([]byte(cached), &postResp); err == nil {
			return postResp, nil
		}
	}

	// Step 1: Get friend posts from repository
	friendPosts, err := s.PostRepository.GetFriendPosts(ctx, s.DB, userID)
	if err != nil {
		if err == sql.ErrNoRows || friendPosts == nil {
			return nil, fmt.Errorf("no friend posts found for user with ID %d", userID)
		}
		return nil, err
	}

	// step 2: Collect all user IDs from friend posts
	userIDMap := make(map[int]bool)
	for _, post := range friendPosts {
		userIDMap[post.UserID] = true
	}

	// step 3: Get all user IDs as a slice
	var userIDs []int
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	// step 4: Query the database to get all users in one go
	users, err := s.UserRepository.FindByIDs(ctx, s.DB, userIDs)
	if err != nil {
		if err == sql.ErrNoRows || users == nil {
			return nil, fmt.Errorf("no users found for friend posts")
		}
		return nil, err
	}

	// step 5: Create a map for easy user lookup by ID
	userMap := make(map[int]*entity.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	var postResponses []*response.CreatePostResponse
	for _, post := range friendPosts {
		user, ok := userMap[post.UserID]
		if !ok {
			fmt.Printf("user with ID %d not found for post ID %d\n", post.UserID, post.PostID)
			continue // Skip this post if user not found
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

	if len(postResponses) == 0 {
		return nil, fmt.Errorf("no friend posts found for user with ID %d", userID)
	}

	// Step 6: Cache the response
	jsonVal, err := json.Marshal(postResponses)
	if err == nil {
		_ = s.RedisClient.Set(ctx, cacheKey, jsonVal, 10*time.Minute).Err()
	}

	return postResponses, nil
}