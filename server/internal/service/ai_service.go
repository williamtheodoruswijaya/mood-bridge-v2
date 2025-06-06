package service

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "os"
    "strings"
    "sync"
    "time"
)

type HuggingFaceResponse struct {
    GeneratedText string `json:"generated_text"`
}

type DialoGPTService struct {
    client  *http.Client
    apiURL  string
    apiKey  string

    mu      sync.Mutex
    history map[int]string // chat history per user ID
}

func NewDialoGPTService() *DialoGPTService {
    return &DialoGPTService{
        client:  &http.Client{Timeout: 60 * time.Second},
        apiURL:  "https://api-inference.huggingface.co/models/HuggingFaceH4/zephyr-7b-beta",
        apiKey:  os.Getenv("HUGGINGFACE_API_TOKEN"),
        history: make(map[int]string),
    }
}

// Chat sends the user input with conversation history and returns the chatbot response
func (s *DialoGPTService) Chat(ctx context.Context, userID int, input string) (string, error) {
    if s.apiKey == "" {
        return "", errors.New("missing Hugging Face API token")
    }

    s.mu.Lock()
    // Retrieve existing chat history or initialize with system prompt if first time
    chatHistory, ok := s.history[userID]
    if !ok {
        chatHistory = "You are a compassionate mental health support assistant. " +
            "Listen carefully and respond empathetically to the user's messages. " +
            "Always be polite, supportive, and encouraging.\n\n"
    }

    // Append user input with prefix
    promptText := chatHistory + "User message: " + input + "\nAssistant:"
    s.mu.Unlock()

    // Construct payload with parameters including stop tokens to prevent runaway generation
    payload := map[string]interface{}{
        "inputs": promptText,
        "parameters": map[string]interface{}{
            "max_new_tokens": 500,
            "temperature":    0.7,
            "stop":           []string{"User message:", "Assistant:"},
        },
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return "", err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL, bytes.NewBuffer(body))
    if err != nil {
        return "", err
    }

    req.Header.Set("Authorization", "Bearer "+s.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := s.client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return "", errors.New("failed to get response from Hugging Face API: " + resp.Status + " - " + string(bodyBytes))
    }

    var hfResp []HuggingFaceResponse
    err = json.NewDecoder(resp.Body).Decode(&hfResp)
    if err != nil {
        return "", err
    }

    if len(hfResp) == 0 {
        return "", errors.New("empty response from model")
    }

    fullOutput := hfResp[0].GeneratedText

    // Use strings.HasPrefix and TrimPrefix to remove prompt from model output if present
    var assistantReply string
    if strings.HasPrefix(fullOutput, promptText) {
        assistantReply = strings.TrimPrefix(fullOutput, promptText)
    } else {
        assistantReply = fullOutput
    }
    assistantReply = strings.TrimSpace(assistantReply)

	assistantReply = strings.TrimSpace(assistantReply)

	// Remove trailing "User message:" if it exists
	if strings.HasSuffix(assistantReply, "User message:") {
		assistantReply = strings.TrimSuffix(assistantReply, "User message:")
		assistantReply = strings.TrimSpace(assistantReply)
	}


    // Update chat history by appending user input + assistant's reply
    s.mu.Lock()
    s.history[userID] = promptText + " " + assistantReply + "\n"
    s.mu.Unlock()

    return assistantReply, nil
}
