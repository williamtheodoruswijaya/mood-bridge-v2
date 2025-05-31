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

type FriendService interface {
	AddFriend(ctx context.Context, req request.FriendRequest) (*response.FriendResponse, error)
	AcceptRequest(ctx context.Context, req request.FriendRequest) (*response.FriendResponse, error)
	GetFriends(ctx context.Context, userID int) ([]*response.FriendResponse, error)
	Delete(ctx context.Context, friendID int) (string, error)
	GetFriendRequests(ctx context.Context, userID int) ([]*response.FriendResponse, error)
}

type FriendServiceImpl struct {
	friendRepository repository.FriendRepository
	userRepository repository.UserRepository
	DB *sql.DB
	RedisClient *redis.Client
}

func NewFriendService(friendRepository repository.FriendRepository, userRepository repository.UserRepository, db *sql.DB, redisClient *redis.Client) FriendService {
	return &FriendServiceImpl{
		friendRepository: friendRepository,
		userRepository: userRepository,
		DB: db,
		RedisClient: redisClient,
	}
}

func (s *FriendServiceImpl) AddFriend(ctx context.Context, req request.FriendRequest) (*response.FriendResponse, error) {
	// step 1: start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	// step 2: rollback transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// step 3: validate request
	if err := utils.ValidateFriendInput(req.UserID, req.FriendUserID, 0); err != nil {
		return nil, err
	}

	// step 4: find user
	user, err := s.userRepository.FindByID(ctx, s.DB, req.UserID)
	if err != nil {
		if err == sql.ErrNoRows || user == nil {
			return nil, fmt.Errorf("user with id %d not found", req.UserID)
		}
		return nil, err
	}

	// step 5: check apakah user dan friend sudah berteman
	exists, err := s.friendRepository.IsFriendExist(ctx, s.DB, req.UserID, req.FriendUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("friend request already exists")
	}

	// step 6: ubah ke entity friend
	friend := &entity.Friend{
		UserID:       req.UserID,
		FriendUserID: req.FriendUserID,
		FriendStatus: false, // false = pending
		CreatedAt:    time.Now(),
	}

	// step 6: jalankan repository
	newFriend, err := s.friendRepository.AddFriend(ctx, tx, friend)
	if err != nil {
		return nil, err
	}
	if newFriend == nil {
		return nil, fmt.Errorf("failed to add friend")
	}

	// step 7: commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// step 8: return response
	resp := &response.FriendResponse{
		FriendID:    newFriend.FriendID,
		UserID: newFriend.UserID,
		FriendUserID: newFriend.FriendUserID,
		FriendStatus: newFriend.FriendStatus,
		CreatedAt: newFriend.CreatedAt,
		User: response.UserSummary{
			UserID: newFriend.UserID,
			Username: newFriend.User.Username,
			FullName: newFriend.User.Fullname,
		},
	}

	// step 9: cache logic
	// 9.1: hapus cache lama dari daftar friend request punya si target karena ada request add friend yang baru
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("friendrequest:%d:v%d", req.FriendUserID, cacheVersion)).Err()
	// opsional: hapus cache daftar teman dari user dan target
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("friend:%d:v%d", req.UserID, cacheVersion)).Err()
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("friend:%d:v%d", req.FriendUserID, cacheVersion)).Err()
	// NOTES: kita gaush simpan request ini ke cache karena akan menimpa semua friend request yang ada di cache

	// step 10: return response
	return resp, nil
}

func (s *FriendServiceImpl) AcceptRequest(ctx context.Context, req request.FriendRequest) (*response.FriendResponse, error) {
	// step 1: start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	// step 2: rollback transaction
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// step 3: validate request
	if err := utils.ValidateFriendInput(req.UserID, req.FriendUserID, 0); err != nil {
		return nil, err
	}

	// step 4: find user (siapa tau usernya gaada tapi malah friend request kan serem)
	user, err := s.userRepository.FindByID(ctx, s.DB, req.FriendUserID)
	if err != nil {
		if err == sql.ErrNoRows || user == nil {
			return nil, fmt.Errorf("user with id %d not found", req.UserID)
		}
		return nil, err
	}

	// step 5: check apakah user sudah berteman
	exists, err := s.friendRepository.IsFriendAlreadyAccepted(ctx, s.DB, req.UserID, req.FriendUserID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("friend request already accepted")
	}
	
	// step 6: jalankan repository-nya
	acceptFriend, err := s.friendRepository.AcceptRequest(ctx, tx, &entity.Friend{
		UserID:       req.UserID,
		FriendUserID: req.FriendUserID,
		FriendStatus: true,
		CreatedAt:    time.Now(), // update created at buat nunjukin kapan dia di-accept
	})

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("friend request not found")
		}
		return nil, err
	}

	// step 7: commit transaction
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// step 8: return response
	resp := &response.FriendResponse{
		FriendID:    acceptFriend.FriendID,
		UserID: acceptFriend.UserID,
		FriendUserID: acceptFriend.FriendUserID,
		FriendStatus: acceptFriend.FriendStatus,
		CreatedAt: acceptFriend.CreatedAt,
		User: response.UserSummary{
			UserID: user.ID,
			Username: user.Username,
			FullName: user.Fullname,
		},
	}

	// step 9: cache invalidation (ketika kita accept sebuah request, maka perubahan dalam database ada pada friend list dan friend request)
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("friend:%d:v%d", req.UserID, cacheVersion)).Err()
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("friend:%d:v%d", req.FriendUserID, cacheVersion)).Err()

	// Hapus juga cache friend request punya si pengirim (karena kalau dia accept, berarti friend request-nya udah gaada)
	_ = s.RedisClient.Del(ctx, fmt.Sprintf("friendrequest:%d:v%d", req.UserID, cacheVersion)).Err()

	// step 10: return response
	return resp, nil
}

func (s *FriendServiceImpl) GetFriends(ctx context.Context, userID int) ([]*response.FriendResponse, error) {
	// step 1: cache key
	friendCacheKey := fmt.Sprintf("friend:%d:v%d", userID, cacheVersion)
	
	// step 2: get cache based on cache key
	cachedFriends, err := s.RedisClient.Get(ctx, friendCacheKey).Result()
	if err == nil {
		var friends []*response.FriendResponse
		if err := json.Unmarshal([]byte(cachedFriends), &friends); err == nil {
			return friends, nil
		}
	}

	// step 3: get friends from repository
	friends, err := s.friendRepository.GetFriends(ctx, s.DB, userID)
	if err != nil {
		return nil, err
	}

	// step 4: convert to response
	var friendResponses []*response.FriendResponse
	for _, friend := range *friends {
		friendResponses = append(friendResponses, &response.FriendResponse{
			FriendID:    friend.FriendID,
			UserID:      friend.UserID,
			FriendUserID: friend.FriendUserID,
			FriendStatus: friend.FriendStatus,
			CreatedAt:    friend.CreatedAt,
			User: response.UserSummary{
				UserID:      friend.User.ID,
				Username:    friend.User.Username,
				FullName:    friend.User.Fullname,
			},
		})
	}

	// step 5: validasi kalau gaada yang ketemu (kosong misalnya)
	if len(friendResponses) == 0 {
		return []*response.FriendResponse{}, fmt.Errorf("no friends found")
	}

	// step 6: ubah ke dalam json untuk disimpan ke cache
	jsonData, err := json.Marshal(friendResponses)
	if err == nil {
		// step 7: kalau gaada error ketika perubahan ke json ini, kita simpan ke cache
		_ = s.RedisClient.Set(ctx, friendCacheKey, jsonData, 10*time.Minute).Err()
	}

	// step 8: return response
	return friendResponses, nil
}

func (s *FriendServiceImpl) Delete(ctx context.Context, friendID int) (string, error) {
	// step 1: start transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return "", err
	}

	// step 2: rollback transaction
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// step 3: jalankan repository-nya
	message, err := s.friendRepository.Delete(ctx, tx, friendID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("friend with id %d not found", friendID)
		}
		return "", err
	}

	// step 4: commit transaction
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	// step 5: delete cache (invalidate cache)
	friendCacheKey := fmt.Sprintf("friend:%d:v%d", friendID, cacheVersion)
	_ = s.RedisClient.Del(ctx, friendCacheKey).Err()

	// step 6: return response
	return message, nil
}

func (s *FriendServiceImpl) GetFriendRequests(ctx context.Context, userID int) ([]*response.FriendResponse, error) {
	// step 1: set cache key
	friendCacheKey := fmt.Sprintf("friendrequest:%d:v%d", userID, cacheVersion)

	// step 2: get cache based on cache key
	cachedFriendRequests, err := s.RedisClient.Get(ctx, friendCacheKey).Result()
	if err == nil {
		var friendRequests []*response.FriendResponse
		if err := json.Unmarshal([]byte(cachedFriendRequests), &friendRequests); err == nil {
			return friendRequests, nil
		}
	}

	// step 3: get friend requests from repository
	friendRequests, err := s.friendRepository.GetFriendRequests(ctx, s.DB, userID)
	if err != nil {
		return nil, err
	}

	// step 4: convert to response
	var friendRequestResponses []*response.FriendResponse
	for _, friendRequest := range *friendRequests {
		friendRequestResponses = append(friendRequestResponses, &response.FriendResponse{
			FriendID:    friendRequest.FriendID,
			UserID:      friendRequest.UserID,
			FriendUserID: friendRequest.FriendUserID,
			FriendStatus: friendRequest.FriendStatus,
			CreatedAt:    friendRequest.CreatedAt,
			User: response.UserSummary{
				UserID:      friendRequest.User.ID,
				Username:    friendRequest.User.Username,
				FullName:    friendRequest.User.Fullname,
			},
		})
	}

	// step 5: validasi kalau gaada yang ketemu (kosong misalnya)
	if len(friendRequestResponses) == 0 {
		return []*response.FriendResponse{}, fmt.Errorf("no friend requests found")
	}

	// step 6: ubah ke dalam json untuk disimpan ke cache
	jsonData, err := json.Marshal(friendRequestResponses)
	if err == nil {
		// step 7: kalau gaada error ketika perubahan ke json ini, kita simpan ke cache
		_ = s.RedisClient.Set(ctx, friendCacheKey, jsonData, 10*time.Minute).Err()
	}

	// step 8: return response
	return friendRequestResponses, nil
}