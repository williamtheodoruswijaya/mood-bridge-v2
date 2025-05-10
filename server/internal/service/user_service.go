package service

import (
	"context"
	"database/sql"
	"mood-bridge-v2/server/internal/entity"
	"mood-bridge-v2/server/internal/model/request"
	"mood-bridge-v2/server/internal/model/response"
	"mood-bridge-v2/server/internal/repository"
	"time"
)

type UserService interface {
	Create(ctx context.Context, request request.CreateUserRequest) (*response.CreateUserResponse, error)
	Find(ctx context.Context, username string) (*response.CreateUserResponse, error)
}

type UserServiceImpl struct {
	DB             *sql.DB
	UserRepository repository.UserRepository
}

func NewUserService(db *sql.DB, userRepository repository.UserRepository) UserService {
	return &UserServiceImpl{
		DB:             db,
		UserRepository: userRepository,
	}
}

func (s *UserServiceImpl) Create(ctx context.Context, request request.CreateUserRequest) (*response.CreateUserResponse, error) {
	// step 1: begin transaction
	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}

	// step 2: rollback
	defer func() error {
		if err != nil {
			tx.Rollback()
			return err
		}
		return nil
	}()

	// step 3: convert request ke model User
	user := entity.User{
		Username:   request.Username,
		Fullname:   request.Fullname,
		Email:      request.Email,
		Password:   request.Password,
		ProfileUrl: sql.NullString{String: "https://upload.wikimedia.org/wikipedia/commons/a/ac/Default_pfp.jpg", Valid: true}, // Default URL for new users
		CreatedAt:  time.Now(),
	}

	// step 4: call repository to create user
	createdUser, err := s.UserRepository.Create(ctx, tx, &user)
	if err != nil {
		return nil, err
	}

	// step 5: commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	result, err := s.Find(ctx, createdUser.Username)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *UserServiceImpl) Find(ctx context.Context, username string) (*response.CreateUserResponse, error) {
	// step 1: call repository to find user
	user, err := s.UserRepository.Find(ctx, s.DB, username)
	if err != nil {
		return nil, err
	}

	// step 2: convert result ke response
	searchedUser := response.CreateUserResponse{
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	// step 3: return response
	return &searchedUser, nil
}
