package handler

import (
	"mood-bridge-v2/server/internal/service"
	"net/http"

	"github.com/gorilla/websocket"
)

type ChatHandler interface {
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
				return origin == "http://localhost:3000" // TODO: ganti dengan origin yang sesuai (link vercel)
			},
		},
	}
}