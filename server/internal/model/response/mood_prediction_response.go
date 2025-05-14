package response

type MoodPredictionResponse struct {
	Prediction string `json:"prediction"`
}

type MoodPredictionResponseList struct {
	Anxiety             float64 `json:"anxiety"`
	Bipolar             float64 `json:"bipolar"`
	Depression          float64 `json:"depression"`
	Normal              float64 `json:"normal"`
	PersonalityDisorder float64 `json:"personality_disorder"`
	Stress              float64 `json:"stress"`
	Suicidal            float64 `json:"suicidal"`
}
