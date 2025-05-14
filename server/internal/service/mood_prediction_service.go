package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mood-bridge-v2/server/internal/model/request"
	"mood-bridge-v2/server/internal/model/response"
	"net/http"
	"os"
)

type MoodPredictionService interface {
	PredictMood(ctx context.Context, request request.MoodPredictionRequest) (*response.MoodPredictionResponse, error)
	PredictMoodProba(ctx context.Context, request request.MoodPredictionRequest) (*response.MoodPredictionResponseList, error)
}

type MoodPredictionServiceImpl struct {
}

func NewMoodPredictionService() MoodPredictionService {
	return &MoodPredictionServiceImpl{}
}

func (s *MoodPredictionServiceImpl) PredictMood(ctx context.Context, request request.MoodPredictionRequest) (*response.MoodPredictionResponse, error) {
	apiURL := os.Getenv("MIC_PREDICT_URL") // ini route API sementara
	
	// step 1: set payload buat request ke API (anggep aja dia ni body-nya kalau di postman)
	payload := map[string]string{"input": request.Input}
	
	// step 2: lakuin json.marshal untuk convert payload ke json
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// step 3: hit api-nya
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // step 4: close body-nya

	// step 5: check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get mood prediction, status code: %d", resp.StatusCode)
	}

	// step 6: decode response body ke struct
	var result response.MoodPredictionResponse
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// step 7: return result
	return &result, nil
}

func (s *MoodPredictionServiceImpl) PredictMoodProba(ctx context.Context, request request.MoodPredictionRequest) (*response.MoodPredictionResponseList, error) {
	apiURL := os.Getenv("MIC_PREDICT_MANY_URL")

	// step 1: set payload buat request ke API (anggep aja dia ni body-nya kalau di postman)
	payload := map[string]string{"input": request.Input}

	// step 2: lakuin json.marshal untuk convert payload ke json
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// step 3: hit api-nya
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close() // step 4: close body-nya

	// step 5: check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get mood prediction, status code: %d", resp.StatusCode)
	}

	// step 6: decode response body ke struct
	var result response.MoodPredictionResponseList
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	// step 7: return result
	return &result, nil
}