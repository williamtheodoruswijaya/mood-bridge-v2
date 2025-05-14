package request

type MoodPredictionRequest struct {
	Input string `json:"input" binding:"required"`
}