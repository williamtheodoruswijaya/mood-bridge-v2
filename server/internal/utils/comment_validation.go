package utils

import "fmt"

func ValidateCommentInput(Content string) error {
	if len(Content) < 1 || len(Content) > 500 {
		return fmt.Errorf("comment content must be between 1 and 500 characters")
	}
	return nil
}