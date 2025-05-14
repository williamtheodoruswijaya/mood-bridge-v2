package utils

import "fmt"

func ValidatePostInput(content string, mood string) error {
	if content == "" || mood == "" {
		return fmt.Errorf("userID, content, and mood are required")
	}
	if len(content) < 1 {
		return fmt.Errorf("content must be at least 1 character long")
	}

	// Kalau semua aman, kita return nil
	return nil
}