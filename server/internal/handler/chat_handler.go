package handler

import (
	"log"
	"mood-bridge-v2/server/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler interface {
	HandleWebSocketConnection(c *gin.Context)
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