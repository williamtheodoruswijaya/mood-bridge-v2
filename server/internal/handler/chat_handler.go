package handler

import (
	"log"
	"mood-bridge-v2/server/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler interface {
	HandleWebSocketConnection(c *gin.Context)
	HandleFetchChatHistory(c *gin.Context)
	HandleMarkMessageAsRead(c *gin.Context)
}

type ChatHandlerImpl struct {
	chatService service.ChatService
	upgrader websocket.Upgrader
}

func NewChatHandler(chatService service.ChatService) ChatHandler {
	return &ChatHandlerImpl{
		chatService: chatService,
		upgrader: websocket.Upgrader{
			ReadBufferSize: 1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				log.Printf("Handler: WebSocket connection request from origin: %s", origin)
				return true // Allow all origins for simplicity, adjust as needed
			},
		},
	}
}

func (h *ChatHandlerImpl) HandleWebSocketConnection(c *gin.Context) {
	// step 1: ambil userID dari auth middleware
	userIDInterface := c.Request.Context().Value("userID")
	if userIDInterface == nil {
		log.Println("Handler: User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}

	// step 2: ubah ke int
	userID, ok := userIDInterface.(int)
	if !ok {
		log.Println("Handler: User ID is not an integer")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}
	
	log.Printf("Handler: Attempting WebSocket upgrade for UserID: %d", userID)// buat test doang

	// step 2: upgrade koneksi ke WebSocket
	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Handler: Failed to upgrade connection: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to upgrade connection",
		})
		return
	}

	log.Printf("Handler: WebSocket connection upgraded for UserID: %d", userID)

	// step 3: panggil service untuk menangani koneksi WebSocket
	err = h.chatService.HandleNewConnection(c.Request.Context(), userID, conn)
	if err != nil {
		log.Printf("Handler: Error handling new connection: %v", err)
		_ = conn.Close() // close the connection if there's an error
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to handle new connection",
		})
		return
	}
}

func (h *ChatHandlerImpl) HandleFetchChatHistory(c *gin.Context) {
	// step 1: ambil userID dari auth middleware
	userIDInterface := c.Request.Context().Value("userID")
	if userIDInterface == nil {
		log.Println("Handler: User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}

	// step 2: ubah ke int
	userID, ok := userIDInterface.(int)
	if !ok {
		log.Println("Handler: User ID is not an integer")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}

	// step 3: ambil recipientID dari query parameter
	recipientIDstr := c.Query("with_user_id")
	if recipientIDstr == "" {
		log.Println("Handler: Recipient ID not provided")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Recipient ID is required",
		})
		return
	}

	recipientID, err := strconv.Atoi(recipientIDstr)
	if err != nil || recipientID <= 0 {
		log.Printf("Handler: Invalid recipient ID: %s", recipientIDstr)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid recipient ID",
		})
		return
	}

	// step 4: ambil limit dan offset dari query parameter
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		log.Printf("Handler: Invalid limit: %s", limitStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid limit",
		})
		return
	}

	offsetStr := c.DefaultQuery("offset", "0")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		log.Printf("Handler: Invalid offset: %s", offsetStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid offset",
		})
		return
	}

	// step 5: panggil service untuk mengambil chat history
	messages, err := h.chatService.FetchConversationHistory(c.Request.Context(), userID, recipientID, limit, offset)
	if err != nil {
		log.Printf("Handler: Error fetching chat history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to fetch chat history",
		})
		return
	}

	log.Printf("Handler: Successfully fetched chat history for UserID: %d with RecipientID: %d", userID, recipientID)
	// step 6: kirim response dengan chat history
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Chat history fetched successfully",
		"data":    messages,
	})
	log.Printf("Handler: Chat history response sent for UserID: %d with RecipientID: %d", userID, recipientID)
}

func (h *ChatHandlerImpl) HandleMarkMessageAsRead(c *gin.Context) {
	// step 1: ambil userID dari auth middleware
	userIDInterface := c.Request.Context().Value("userID")
	if userIDInterface == nil {
		log.Println("Handler: User ID not found in context")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}

	// step 2: ubah ke int
	userID, ok := userIDInterface.(int)
	if !ok {
		log.Println("Handler: User ID is not an integer")
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Internal Server Error",
		})
		return
	}

	// step 3: ambil messageID dari query parameter
	messageIDStr :=  c.Param("message_id")
	if messageIDStr == "" {
		log.Println("Handler: Message ID not provided")
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Message ID is required",
		})
		return
	}

	// step 4: ubah messageID ke int
	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil || messageID <= 0 {
		log.Printf("Handler: Invalid message ID: %s", messageIDStr)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid message ID",
		})
		return
	}

	// step 5: panggil service untuk menandai pesan sebagai dibaca
	err = h.chatService.MarkMessageAsRead(c.Request.Context(), userID, messageID)
	if err != nil {
		log.Printf("Handler: Error marking message as read: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to mark message as read",
		})
		return
	}

	log.Printf("Handler: Successfully marked message as read for UserID: %d, MessageID: %d", userID, messageID)
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Message marked as read successfully",
		"data":    nil,
	})
}