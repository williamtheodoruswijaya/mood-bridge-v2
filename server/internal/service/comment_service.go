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

type CommentService interface {
	Create(ctx context.Context, req request.CreateCommentRequest) (*response.CreateCommentResponse, error)
	GetAllByPostID(ctx context.Context, postID int) ([]*response.CreateCommentResponse, error)
	Delete(ctx context.Context, commentID int) (string, error)
	GetByID(ctx context.Context, commentID int) (*response.CreateCommentResponse, error)
}

type CommentServiceImpl struct {
	commentRepository repository.CommentRepository
	userRepository repository.UserRepository
	postRepository repository.PostRepository
	DB                *sql.DB
	RedisClient *redis.Client
}

func NewCommentService(commentRepository repository.CommentRepository, userRepository repository.UserRepository, postRepository repository.PostRepository, db *sql.DB, redisClient *redis.Client) CommentService {
	return &CommentServiceImpl{
		commentRepository: commentRepository,
		userRepository:    userRepository,
		postRepository:    postRepository,
		DB:                db,
		RedisClient: redisClient,
	}
}

func (s *CommentServiceImpl) Create(ctx context.Context, req request.CreateCommentRequest) (*response.CreateCommentResponse, error) {
	// Step 1: Start a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	// Step 2: Siap-siap buat rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// step 3: Validate kalau user exist
	user, err := s.userRepository.FindByID(ctx, s.DB, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows || user == nil{
			return nil, fmt.Errorf("user with ID %d not found", req.UserID)
		}
		return nil, err
	}

	// step 4: Validate kalau post exist
	post, err := s.postRepository.Find(ctx, s.DB, req.PostID)
	if err != nil {
		if err == sql.ErrNoRows || post == nil {
			return nil, fmt.Errorf("post with ID %d not found", req.PostID)
		}
		return nil, err
	}

	// step  5: validate struktur input komentar-nya
	if err := utils.ValidateCommentInput(req.Content); err != nil {
		return nil, err
	}

	// step 6: Insert ke dalam database
	comment := entity.Comment{
		PostID:   req.PostID,
		UserID:   req.UserID,
		Content:  req.Content,
		CreatedAt: time.Now(),
	}
	createdComment, err := s.commentRepository.Create(ctx, tx, &comment)
	if err != nil {
		return nil, err
	}

	// step 7: Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// step 8: Cari comment yang baru saja di-insert
	commentResult, err := s.GetByID(ctx, createdComment.CommentID)
	if err != nil {
		return nil, err
	}

	// step 9: Convert ke dalam response
	commentResp := &response.CreateCommentResponse{
		CommentID: commentResult.CommentID,
		PostID:   commentResult.PostID,
		UserID:   commentResult.UserID,
		User: response.UserSummary{
			UserID:    user.ID,
			Username:  user.Username,
			FullName: user.Fullname,
		},
		Content:   commentResult.Content,
		CreatedAt: commentResult.CreatedAt,
	}

	// step 10: Invalidate cache pada comment dan post (post tetap harus di-invalidate karena ada perubahan pada komentar)
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("comment:%d:v%d", commentResult.CommentID, cacheVersion)).Err()
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("post:%d:v%d", commentResult.PostID, cacheVersion)).Err()

	// step 11: Kembalikan commentResp-nya
	return commentResp, nil
}

func (s *CommentServiceImpl) GetAllByPostID(ctx context.Context, postID int) ([]*response.CreateCommentResponse, error) {
	// step 1: Ambil cacheKey-nya, ini contoh cara penyimpanannya di redis (cacheKey -> comment:postID:v1)
	cacheKey := fmt.Sprintf("comment:post:%d:v%d", postID, cacheVersion)

	// step 2: Ambil data-nya berdasarkan cacheKey-nya
	cached, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil { // Kalau gaada error
		var commentsResp []*response.CreateCommentResponse
		// step 3: Decode JSON-nya dari redis ke dalam struct
		if err := json.Unmarshal([]byte(cached), &commentsResp); err == nil {
			return commentsResp, nil
		}
	}

	// step 4: cari postingan-nya dalam database (pastiin postingannya ada)
	post, err := s.postRepository.Find(ctx, s.DB, postID)
	if err != nil {
		if err == sql.ErrNoRows || post == nil {
			return nil, fmt.Errorf("post with ID %d not found", postID)
		}
		return nil, err
	}

	// step 5: cari semua komentar-nya
	comments, err := s.commentRepository.GetAllByPostID(ctx, s.DB, postID)
	if err != nil {
		if err == sql.ErrNoRows || comments == nil {
			return nil, fmt.Errorf("comments for post with ID %d not found", postID)
		}
		return nil, err
	}
	
	// step 6: Lakukan batch query untuk ambil user-nya
	userIDMap := map[int]bool{}
	for _, comment := range comments {
		userIDMap[comment.UserID] = true
	}

	var userIDs []int
	for id := range userIDMap {
		userIDs = append(userIDs, id)
	}

	users, err := s.userRepository.FindByIDs(ctx, s.DB, userIDs)
	if err != nil {
		if err == sql.ErrNoRows || users == nil {
			return nil, fmt.Errorf("users for post with ID %d not found", postID)
		}
		return nil, err
	}

	userMap := map[int]*entity.User{}
	for _, user := range users {
		userMap[user.ID] = user
	}

	// step 7: Convert ke dalam response
	var commentResps []*response.CreateCommentResponse
	for _, comment := range comments {
		user := userMap[comment.UserID]
		if user == nil {
			return nil, fmt.Errorf("user with ID %d not found", comment.UserID)
		}

		commentResp := &response.CreateCommentResponse{
			CommentID: comment.CommentID,
			PostID:   comment.PostID,
			UserID:   comment.UserID,
			User: response.UserSummary{
				UserID:    user.ID,
				Username:  user.Username,
				FullName: user.Fullname,
			},
			Content:   comment.Content,
			CreatedAt: comment.CreatedAt,
		}
		commentResps = append(commentResps, commentResp)
	}

	// step 6: validasi kalau gaada komentar
	if len(commentResps) == 0 {
		return nil, fmt.Errorf("comments for post with ID %d not found", postID)
	}

	// step 7: Ubah ke dalam JSON untuk disimpan ke dalam cache
	jsonData, err := json.Marshal(commentResps)
	if err == nil {
		// step 8: Kalau gaada error, simpan ke dalam cache
		err = s.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err()
	}

	// step 9: Kembalikan commentResps-nya
	return commentResps, nil
}

func (s *CommentServiceImpl) Delete(ctx context.Context, commentID int) (string, error) {
	// step 1: start a transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return "", err
	}

	// step 2: siapkan rollback
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// step 3: cari comment-nya
	comment, err := s.commentRepository.GetByID(ctx, s.DB, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("comment with ID %d not found", commentID)
		}
		return "", err
	}

	// step 4: delete comment-nya
	message, err := s.commentRepository.Delete(ctx, tx, commentID)
	if err != nil {
		return "", err
	}

	// step 5: commit transaction
	if err := tx.Commit(); err != nil {
		return "", err
	}

	// step 6: invalidate cache pada comment dan post (post tetap harus di-invalidate karena ada perubahan pada komentar)
	s.RedisClient.Del(ctx, fmt.Sprintf("comment:%d:v%d", comment.CommentID, cacheVersion)).Err()
	s.RedisClient.Del(ctx, fmt.Sprintf("post:%d:v%d", comment.PostID, cacheVersion)).Err()

	// step 7: kembalikan message-nya
	return message, nil
}

func (s *CommentServiceImpl) GetByID(ctx context.Context, commentID int) (*response.CreateCommentResponse, error) {
	// Step 1: Ambil cacheKey-nya, ini contoh cara penyimpanannya di redis (cacheKey -> comment:commentID:v1)
	cacheKey := fmt.Sprintf("comment:%d:v%d", commentID, cacheVersion)
	
	// Step 2: Ambil data-nya berdasarkan cacheKey-nya
	cached, err := s.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil { // Kalau gaada error
		var commentResp response.CreateCommentResponse
		// Step 3: Decode JSON-nya dari redis ke dalam struct
		if err := json.Unmarshal([]byte(cached), &commentResp); err != nil {
			return &commentResp, nil
		}
	}

	// Step 4: Cari komentar-nya dalam database
	comment, err := s.commentRepository.GetByID(ctx, s.DB, commentID)
	if err != nil {
		if err == sql.ErrNoRows || comment == nil {
			return nil, fmt.Errorf("comment with ID %d not found", commentID)
		}
	}

	// Step 5: Load User dan Post-nya
	user, err := s.userRepository.FindByID(ctx, s.DB, comment.UserID)
	if err != nil {
		if err == sql.ErrNoRows || user == nil {
			return nil, fmt.Errorf("user with ID %d not found", comment.UserID)
		}
	}
	post, err := s.postRepository.Find(ctx, s.DB, comment.PostID)
	if err != nil {
		if err == sql.ErrNoRows || post == nil {
			return nil, fmt.Errorf("post with ID %d not found", comment.PostID)
		}
	}

	// Step 6: Convert ke response
	commentResp := &response.CreateCommentResponse{
		CommentID: comment.CommentID,
		PostID:   comment.PostID,
		UserID:   comment.UserID,
		User: response.UserSummary{
			UserID:    user.ID,
			Username:  user.Username,
			FullName: user.Fullname,
		},
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
	}

	// Step 7: Ubah ke dalam JSON untuk disimpan ke dalam cache
	jsonData, err := json.Marshal(commentResp)
	if err == nil {
		// Step 8: Kalau gaada error, simpan ke dalam cache
		_ = s.RedisClient.Set(ctx, cacheKey, jsonData, 10*time.Minute).Err()
	}

	// step 9: Kembalikan commentResp-nya
	return commentResp, nil
}