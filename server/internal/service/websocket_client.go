package service

/*
	WebSocket Client Service:
	- WebSocket adalah protokol full-duplex (dua arah) yang memungkinkan komunikasi real-time antara client dan server.
	- Setiap client yang terhubung lewat WebSocket membutuhkan sebuah struct Client yang menyimpan informasi tentang koneksi, user ID, dan saluran untuk mengirim pesan.
	- Client memiliki dua goroutine utama: ReadPump untuk membaca pesan dari client dan WritePump untuk mengirim pesan ke client.
	- ReadPump menangani pesan masuk dari client, memprosesnya, dan meneruskan pesan tersebut ke ChatService untuk penanganan lebih lanjut.
	- WritePump menangani pengiriman pesan keluar ke client, termasuk pesan offline yang mungkin diterima saat client tidak terhubung.
	- Client juga menangani ping/pong untuk menjaga koneksi tetap hidup dan mendeteksi jika client terputus.
	- Client ini diibaratkan sebagai jembatan antara WebSocket dan ChatService
*/

import (
	"encoding/json"
	"log"
	"mood-bridge-v2/server/internal/model/request"
	"mood-bridge-v2/server/internal/model/response"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second // Timeout saat kirim pesan
	pongWait = 60 * time.Second // Timeout jika client tidak merespon ping
	pingPeriod = (pongWait * 9) / 10 // Waktu kirim ping ke client
	maxMessageSize = 1024 * 4 // Maksimum ukuran pesan yang diterima dari client
)

var newline = '\n'

type Client struct {
	UserID int // ID unik untuk mengidentifikasi client
	Hub Hub // Hub yang mengelola koneksi client
	Conn *websocket.Conn // Koneksi WebSocket untuk client yang aktif
	Send chan []byte // Channel untuk mengirim pesan ke client
	ChatService ChatService // Service yang menangani logika chat, seperti mengirim pesan, mengambil riwayat chat, dll.
}

func NewClient(userID int, hub Hub, conn *websocket.Conn, chatService ChatService) *Client {
	return &Client{
		UserID: userID,
		Hub: hub,
		Conn: conn,
		Send: make(chan []byte, 256),
		ChatService: chatService,
	}
}

func (c *Client) ReadPump() {
	// step 1: ketika koneksi terputus, unregister client dari Hub dan tutup koneksi
	defer func() {
		c.Hub.UnregisterClient(c)
		c.Conn.Close()
		log.Printf("Client %d disconnected, readPump closed", c.UserID)
	}()

	// step 2: atur batasan ukuran pesan yang diterima dari client dan atur timeout untuk membaca pesan
	c.Conn.SetReadLimit(maxMessageSize)
	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))

	// step 3: atur handler dengan memperbarui deadline setiap kali menerima pong dari client
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait));
		return nil
	})

	// step 4: baca pesan dari client dalam loop
	for {
		messageType, messageBytes, err := c.Conn.ReadMessage()
		if err != nil { // jika terjadi error, keluar dari loop
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message for client %d: %v", c.UserID, err)
			}
			break
		}

		// step 5: jika pesan yang diterima adalah teks, proses pesan tersebut
		if messageType == websocket.TextMessage {
			var msgPayload request.PrivateMessagePayload
			if err := json.Unmarshal(messageBytes, &msgPayload); err != nil {
				log.Printf("Error unmarshalling message from client %d: %v. Message: %s", c.UserID, err, string(messageBytes))
				errorMsg := response.WebSocketMessage{
					Type: "error",
					Payload: response.ErrorMessage{Code: "invalid_message", Message: "Invalid message format"},
				}
				errorBytes, _ := json.Marshal(errorMsg)
				select {
				case c.Send <- errorBytes:
				default:
					log.Printf("Error sending error message to client %d: send channel is full", c.UserID)
				}
				continue
			}

			// step 6: jika pesan valid, kirim pesan ke ChatService untuk diproses
			err := c.ChatService.HandleIncomingMessage(nil, c.UserID, msgPayload.RecipientID, msgPayload.Content)
			if err != nil {
				log.Printf("Error from ChatService.HandleIncomingMessage for client %d: %v", c.UserID, err)
				errMsgPayload := response.WebSocketMessage{
					Type: "message_send_failed",
					Payload: response.ErrorMessage{
						Code: "send_error",
						Message: err.Error(),
					},
				}
				errMsgBytes, _ := json.Marshal(errMsgPayload)
				select {
					case c.Send <- errMsgBytes:
					default:
					log.Printf("Error sending error message to client %d: send channel is full", c.UserID)
				}
			}
		}
	}
}

func (c *Client) WritePump() {
	// step 1: define timer untuk mengirim ping ke client secara berkala
	ticker := time.NewTicker(pingPeriod)
	defer func(){
		ticker.Stop()
		c.Conn.Close()
		log.Printf("Client %d disconnected, writePump closed", c.UserID)
	}()
	for {
		select {
		case message, ok := <-c.Send:
			// step 2: set timeout untuk menulis pesan ke client
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			// step 3: jika channel Send sudah ditutup, kirim pesan close ke client dan keluar dari loop
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// step 4: jika pesan yang diterima adalah teks, kirim pesan ke client
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting next writer for client %d: %v", c.UserID, err)
				return
			}
			_, _ = w.Write(message)

			// step 5: jika ada pesan lain yang ada di channel Send, kirim pesan tersebut ke client (biasanya pesan offline)
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, _ = w.Write([]byte{byte(newline)})
				_, _ = w.Write(<-c.Send)
			}
			if err := w.Close(); err != nil {
				return
			}

			// step 6: kirim ping tiap interval untuk menjaga koneksi tetap hidup
		case <- ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Error sending ping to client %d: %v", c.UserID, err)
				return
			}
		}
	}
}