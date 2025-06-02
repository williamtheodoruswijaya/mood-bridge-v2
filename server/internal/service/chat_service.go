package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mood-bridge-v2/server/internal/entity"
	"mood-bridge-v2/server/internal/model/response"
	"mood-bridge-v2/server/internal/repository"
	"time"

	"github.com/gorilla/websocket"
)

type ChatService interface {
	HandleNewConnection(ctx context.Context, userID int, conn *websocket.Conn) error
	HandleIncomingMessage(ctx context.Context, senderID int, recipientID int, content string) error
	FetchConversationHistory(ctx context.Context, senderID, recipientID, limit, offset int) ([]*response.ChatMessage, error)
	MarkMessageAsRead(ctx context.Context, messageID, userID int) error
}

type ChatServiceImpl struct {
	messageRepo repository.ChatRepository
	hub Hub
}

func NewChatService(msgRepo repository.ChatRepository, hub Hub) ChatService {
	return &ChatServiceImpl{
		messageRepo: msgRepo,
		hub: hub,
	}
}

func (s *ChatServiceImpl) HandleNewConnection(ctx context.Context, userID int, conn *websocket.Conn) error {
	client := NewClient(userID, s.hub, conn, s)
	s.hub.RegisterClient(client)
	log.Printf("ChatService: User %d connected. Client registered with Hub.", userID)
	unreadMessages, err := s.messageRepo.GetUnreadMessagesForUser(ctx, userID, time.Now().Add(-7*24*time.Hour))
	if err != nil {
		log.Printf("ChatService: Error fetching unread messages for user %d: %v", userID, err)
	}
	if len(unreadMessages) > 0 {
		log.Printf("ChatService: Sending %d unread messages to user %d", len(unreadMessages), userID)
		for _, msg := range unreadMessages {
			wsMsg := response.WebSocketMessage{
				Type: "offline_message",
				Payload: response.ChatMessage{
					ID: msg.ID,
					SenderID: msg.SenderID,
					RecipientID: msg.RecipientID,
					Content: msg.Content,
					Timestamp: msg.Timestamp,
					Status: msg.Status,
				},
			}
			payloadBytes, err := json.Marshal(wsMsg)
			if err != nil {
				log.Printf("ChatService: Error marshalling offline message %d for user %d: %v", msg.ID, userID, err)
				continue
			}
			select {
			case client.Send <- payloadBytes:
				log.Printf("ChatService: Offline message %d sent to user %d", msg.ID, userID)
			default:
				log.Printf("ChatService: Client %d's send channel is full, skipping message %d", userID, msg.ID)
			}
		}
	}
	return nil
}

func (s *ChatServiceImpl) HandleIncomingMessage(ctx context.Context, senderID int, recipientID int, content string) error {
	if content == "" {
		return errors.New("message content cannot be empty")
	}
	if recipientID <= 0 {
		return errors.New("invalid recipient ID")
	}
	if senderID == recipientID {
		return errors.New("sender and recipient cannot be the same")
	}

	msg := entity.NewMessage(senderID, recipientID, content)
	if err := s.messageRepo.SaveMessage(ctx, msg); err != nil {
		log.Printf("ChatService: Error saving message from %d to %d: %v", senderID, recipientID, err)
		return fmt.Errorf("failed to save message: %w", err)
	}
	log.Printf("ChatService: Message from %d to %d saved successfully", senderID, recipientID)
	s.hub.RoutePrivateMessage(msg)
	return nil
}

func (s *ChatServiceImpl) FetchConversationHistory(ctx context.Context, senderID, recipientID, limit, offset int) ([]*response.ChatMessage, error) {
	if limit <= 0 || offset < 0 {
		limit = 20 // default limit
		offset = 0 // default offset
	}
	messages, err := s.messageRepo.GetMessagesForConversation(ctx, senderID, recipientID, limit, offset)
	if err != nil {
		log.Printf("ChatService: Error fetching conversation history between %d and %d: %v", senderID, recipientID, err)
		return nil, fmt.Errorf("failed to fetch conversation history: %w", err)
	}

	var chatMessages []*response.ChatMessage
	for _, msg := range messages {
		chatMessages = append(chatMessages, &response.ChatMessage{
			ID: msg.ID,
			SenderID: msg.SenderID,
			RecipientID: msg.RecipientID,
			Content: msg.Content,
			Timestamp: msg.Timestamp,
			Status: msg.Status,
		})
	}

	log.Printf("ChatService: Fetched %d messages for conversation between %d and %d", len(chatMessages), senderID, recipientID)
	return chatMessages, nil
}

func (s *ChatServiceImpl) MarkMessageAsRead(ctx context.Context, messageID, userID int) error {
	if messageID <= 0 || userID <= 0 {
		return errors.New("invalid message ID or user ID")
	}

	err := s.messageRepo.UpdateMessageStatus(ctx, messageID, entity.StatusRead)
	if err != nil {
		log.Printf("ChatService: Error marking message %d as read for user %d: %v", messageID, userID, err)
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	log.Printf("ChatService: Message %d marked as read for user %d", messageID, userID)
	return nil
}