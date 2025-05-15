package utils

import "fmt"

func ValidateFriendInput(userID, friendUserID, friendStatus int) error {
	if userID <= 0 {
		return fmt.Errorf("user_id must be greater than 0")
	}
	if friendUserID <= 0 {
		return fmt.Errorf("friend_user_id must be greater than 0")
	}
	if friendStatus < 0 || friendStatus > 2 {
		return fmt.Errorf("friend_status must be between 0 and 2")
	}
	if userID == friendUserID {
		return fmt.Errorf("user_id and friend_user_id cannot be the same")
	}
	return nil
}