package entity

type FriendRecommendation struct {
	User        User   `json:"user"`
	OverallMood string `json:"overall_mood"`
}