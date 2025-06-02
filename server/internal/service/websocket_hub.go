package service

/*
	WebScoket Hub Service:
	- berfungsi sebagai PENGHUBUNG UTAMA (pusat koordinasi = hub) antara client yang terhubung dengan server secara real-time
	- memiliki use-case sebagai berikut:
		1. Run(): menjalankan hub untuk menangani koneksi WebSocket, menerima pesan dari client, dan mengirim pesan ke client.
		2. RegisterClient(client *Client): mendaftarkan client baru yang terhubung ke hub saat koneksi WebSocket dibuka.
		3. UnregisterClient(client *Client): menghapus client yang terputus dari hub saat koneksi WebSocket ditutup.
		4. RoutePrivateMessage(message *entity.Message):
			- meneruskan pesan pribadi dari satu client ke client lain yang dituju.
			- misal: client A kirim pesan ke client B -> hub akan menerima pesan tersebut dan mengirimkannya ke client B.
*/

import (
	"context"
	"encoding/json"
	"log"
	"mood-bridge-v2/server/internal/entity"
	"mood-bridge-v2/server/internal/model/response"
	"mood-bridge-v2/server/internal/repository"
	"sync"
)

type Hub interface {
	Run()
	RegisterClient(client *Client)
	UnregisterClient(client *Client)
	RoutePrivateMessage(message *entity.Message)
}

type HubImpl struct { // berfungsi untuk menyimpan daftar client yang aktif (yang terhubung via WebSocket) dan menyediakan cara untuk mengatur koneksi tersebut.
	// menyimpan daftar clients yang terhubung dengan userID sebagai key
	clients map[int]*Client
	// mengamankan akses clients agar thread-safe saat ada banyak koneksi WebSocket (intinya saat ada banyak client yang terhubung ke server, kita perlu mengamankan akses ke map clients)
	clientsMutex sync.RWMutex
	// menyimpan pesan ke database, jadi tidak hanya dikirim tapi juga disimpan untuk riwayat chat
	messageRepo repository.ChatRepository
}

func NewConcreteHub(msgRepo repository.ChatRepository) Hub {
	return &HubImpl{
		clients: make(map[int]*Client),
		messageRepo: msgRepo,
	}
}

func (h *HubImpl) Run() {
	log.Println("WebSocket Hub is running...")
}

func (h *HubImpl) RegisterClient(client *Client) {
	// step 1: kunci thread (mutex) untuk mengamankan akses ke map clients
	h.clientsMutex.Lock()

	// step 2: unlock mutex jika sudah selesai
	defer h.clientsMutex.Unlock()

	// step 3: tambahkan client ke map clients
	log.Printf("Registering client: %d", client.UserID)
	h.clients[client.UserID] = client

	// step 4: jalankan writePump dan readPump untuk client yang baru terdaftar
	go client.WritePump()
	go client.ReadPump()
}


func (h *HubImpl) UnregisterClient(client *Client) {
	// step 1: kunci thread (mutex) untuk mengamankan akses ke map clients
	h.clientsMutex.Lock()

	// step 2: unlock mutex jika sudah selesai
	defer h.clientsMutex.Unlock()

	// step 3: hapus client dari map clients
	if _, ok := h.clients[client.UserID]; ok {
		log.Printf("Unregistering client: %d", client.UserID)
		delete(h.clients, client.UserID)
		close(client.Send)
	}
}

func (h *HubImpl) RoutePrivateMessage(message *entity.Message) {
	// step 1: kunci thread (mutex) untuk mengamankan akses ke map clients
	h.clientsMutex.RLock()

	// step 2: tentukan apakah penerima client ada di map clients
	recipientClient, ok := h.clients[message.RecipientID]

	// step 3: unlock mutex jika sudah selesai
	defer h.clientsMutex.RUnlock()

	// step 4: jika penerima client ada, kirim pesan ke penerima
	log.Printf("Routing private message from %d to %d: %s", message.SenderID, message.RecipientID, message.Content)
	if ok {
		// jika recipientClient online, coba kirim pesan ke recipientClient dalam bentuk json (step ini pertama ubah dulu ke bentuk json)
		wsMsg := response.WebSocketMessage{
			Type: "new_private_message",
			Payload: response.ChatMessage{
				ID: message.ID,
				SenderID: message.SenderID,
				RecipientID: message.RecipientID,
				Content: message.Content,
				Timestamp: message.Timestamp,
				Status: message.Status, // status pesan ubah jadi "sent"
			},
		}
		payloadBytes, err := json.Marshal(wsMsg)
		if err != nil {
			log.Printf("Error marshalling message for recipient %d: %v", message.RecipientID, err)
			return
		}

		// step 5: kirim pesan ke recipientClient
		select {
		case recipientClient.Send <- payloadBytes:
			log.Printf("Message sent to recipient %d: %s", message.RecipientID, message.Content)
			// jika pesan berhasil dikirim, update status pesan di database
			if h.messageRepo != nil {
				go func() {
					err := h.messageRepo.UpdateMessageStatus(context.Background(), message.ID, entity.StatusDelivered)
					if err != nil {
						log.Printf("Error updating message status for message ID %d: %v", message.ID, err)
					} else {
						log.Printf("Message status updated to 'delivered' for message ID %d", message.ID)
					}
				}()
			}
		default:
			log.Printf("Hub: Send channel full for recipient %d. Message ID: %d", message.RecipientID, message.ID) // artinya pesan sudah disimpan di DB tapi akan diambil saat client reconnext
		}
	} else {
		// step 6: jika client penerima tidak ada (offline), simpan pesan ke database dan akan diambil saat client reconnect
		log.Printf("Recipient client %d is offline. Message ID: %d stored in DB", message.RecipientID, message.ID)
	}
}