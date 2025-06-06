package handler

import (
    "context"
    "net/http"
    "mood-bridge-v2/server/internal/service"

    "github.com/gin-gonic/gin"
)

type AIChatHandler struct {
    AIService *service.DialoGPTService
}

func NewAIChatHandler(aiService *service.DialoGPTService) *AIChatHandler {
    return &AIChatHandler{AIService: aiService}
}

type ChatRequest struct {
    UserID  int    `json:"user_id"`
    Message string `json:"message"`
}

type ChatResponse struct {
    Response string `json:"response"`
}

func (h *AIChatHandler) HandleChat(c *gin.Context) {
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    reply, err := h.AIService.Chat(context.Background(), req.UserID, req.Message)
    if err != nil {
        c.Error(err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get response from AI"})
        return
    }

    c.JSON(http.StatusOK, ChatResponse{
        Response: reply,
    })
}

