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
	"strings"
	"time"
)

type UserService interface {
	Create(ctx context.Context, request request.CreateUserRequest) (*response.CreateUserResponse, error)
	Find(ctx context.Context, username string) (*response.CreateUserResponse, error)
	FindByEmail(ctx context.Context, email string) (*response.CreateUserResponse, error)
	FindAll(ctx context.Context) ([]*response.CreateUserResponse, error)
	Login(ctx context.Context, request request.ValidateUserRequest) (*string, error)
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
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// step 3: convert request ke model User
	user := entity.User{
		Username:   strings.TrimSpace(strings.ToLower(request.Username)),
		Fullname:   request.Fullname,
		Email:      strings.TrimSpace(strings.ToLower(request.Email)),
		Password:   request.Password,
		ProfileUrl: sql.NullString{String: "https://upload.wikimedia.org/wikipedia/commons/a/ac/Default_pfp.jpg", Valid: true}, // Default URL for new users
		CreatedAt:  time.Now(),
	}
	
	// Appendix: validate request
	if err := utils.ValidateUserInput(&user); err != nil {
		return nil, err
	}

	// Appendix: validate if username already exists
	existingUser, err := s.UserRepository.Find(ctx, s.DB, user.Username)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("username %s already exists", user.Username)
	}

	// Appendix: validate if email already exists
	existingEmail, err := s.UserRepository.FindByEmail(ctx, s.DB, user.Email)
	if err == nil && existingEmail != nil {
		return nil, fmt.Errorf("email %s already exists", user.Email)
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

	// step 6: Find the created user
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

func (s *UserServiceImpl) FindByEmail(ctx context.Context, email string) (*response.CreateUserResponse, error) {
	// step 1: call repository to find user
	user, err := s.UserRepository.FindByEmail(ctx, s.DB, email)
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

func (s *UserServiceImpl) FindAll(ctx context.Context) ([]*response.CreateUserResponse, error) {
	// step 1: call repository to find all users
	users, err := s.UserRepository.FindAll(ctx, s.DB)
	if err != nil {
		return nil, err
	}

	// step 2: convert result ke response
	var userResponses []*response.CreateUserResponse
	for _, user := range users {
		userResponse := &response.CreateUserResponse{
			Username:  user.Username,
			Fullname:  user.Fullname,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
		userResponses = append(userResponses, userResponse)
	}

	// step 3: return response
	return userResponses, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, request request.ValidateUserRequest) (*string, error) {
	// step 1: validate request
	err := utils.ValidateUserLoginInput(request.Username, request.Password)
	if err != nil {
		return nil, err
	}

	// step 2: call repository to find user by Username
	user, err := s.UserRepository.Find(ctx, s.DB, request.Username)
	if err != nil {
		return nil, err
	}

	// step 3: check if user is found
	if user == nil {
		return nil, fmt.Errorf("user %s not found", request.Username)
	}

	// step 4: validate password
	if user.Password != request.Password {
		return nil, fmt.Errorf("invalid password")
	}

	// step 5: get user response
	userResponse := &response.CreateUserResponse{
		Username:  user.Username,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	// step 6: generate token
	token, err := GenerateToken(userResponse)
	if err != nil {
		return nil, err
	}

	return token, nil
}
