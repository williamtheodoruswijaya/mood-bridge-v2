package utils

import (
	"fmt"
	"mood-bridge-v2/server/internal/entity"
	"regexp"
)

func ValidateUserInput(user *entity.User) error {
	if user.Username == "" || user.Fullname == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("username, fullname, email, and password are required")
	}
	if len(user.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(user.Username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}
	if len(user.Fullname) < 3 {
		return fmt.Errorf("fullname must be at least 3 characters long")
	}

	// check email format using regex
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, user.Email)
	if err != nil || !matched {
		return fmt.Errorf("invalid email format")
	}

	// Kalau semua aman, kita return nil
	return nil
}

func ValidateUserLoginInput(username, password string) error {
	if username == "" || password == "" {
		return fmt.Errorf("username and password are required")
	}
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	if len(username) < 3 {
		return fmt.Errorf("username must be at least 3 characters long")
	}

	// Kalau semua aman, kita return nil
	return nil
}
