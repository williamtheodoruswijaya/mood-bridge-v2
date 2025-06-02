package repository

import (
	"context"
	"database/sql"
	"fmt"
	"mood-bridge-v2/server/internal/entity"
	"time"
)

type ChatRepository interface {
	SaveMessage(ctx context.Context, msg *entity.Message) error
	GetMessagesForConversation(ctx context.Context, senderID, recipientID, limit, offset int) ([]*entity.Message, error)
	UpdateMessageStatus(ctx context.Context, messageID int, newStatus entity.MessageStatus) error
	GetUnreadMessagesForUser(ctx context.Context, userID int, afterTimestamp time.Time) ([]*entity.Message, error)
}

type ChatRepositoryImpl struct {
	DB *sql.DB // Tujuan ini ditaro disini adalah agar DB tidak harus ada dalam hub, client, dsbnya tapi cukup ada pada saat inisialisasi di routes
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &ChatRepositoryImpl{
		DB: db,
	}
}

func (r *ChatRepositoryImpl) SaveMessage(ctx context.Context, msg *entity.Message) error {
	// step 1: define query buat masukin pesan yang ada ke dalam database
	query := `INSERT INTO messages (senderid, recipientid, content, timestamp, status) VALUES ($1, $2, $3, $4, $5)`

	// step 2: jalankan query-nya (ExecContext) digunakan untuk eksekusi query yang tidak mengembalikan baris (INSERT, UPDATE, DELETE) (Berbeda dengan QueryRowContext yang digunakan untuk query yang mengembalikan satu baris dan QueryContext yang digunakan untuk query yang mengembalikan banyak baris)
	_, err := r.DB.ExecContext(ctx, query, msg.SenderID, msg.RecipientID, msg.Content, msg.Timestamp, msg.Status)
	if err != nil {
		return fmt.Errorf("error saving message: %w", err)
	}

	// step 3: jika berhasil, return nil
	return nil
}

func (r *ChatRepositoryImpl) GetMessagesForConversation(ctx context.Context, senderID, recipientID, limit, offset int) ([]*entity.Message, error) {
	// step 1: define query untuk mengambil pesan dari database baik yang sudah di read maupun yang belum di read
	query := `
	SELECT id, senderid, recipientid, content, timestamp, status
	FROM messages
	WHERE (senderid = $1 AND recipientid = $2) OR (senderid = $2 AND recipientid = $1)
	ORDER BY timestamp ASC
	LIMIT $3 OFFSET $4
	`

	// step 2: jalankan query-nya (QueryContext) digunakan untuk eksekusi query yang mengembalikan banyak baris (SELECT)
	rows, err := r.DB.QueryContext(ctx, query, senderID, recipientID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error retrieving messages: %w", err)
	}

	// step 3: pastikan untuk menutup rows setelah selesai digunakan
	defer rows.Close()

	// step 4: buat slice untuk menyimpan pesan yang diambil dari database
	var messages []*entity.Message
	
	// step 5: iterasi melalui rows untuk mengambil setiap pesan dan scan ke dalam struct entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.Timestamp, &msg.Status); err != nil {
			fmt.Printf("error scanning message: %v\n", err)
			continue
		}
		messages = append(messages, &msg)
	}

	// step 6: periksa apakah ada error saat iterasi
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over messages: %w", err)
	}

	// step 7: jika berhasil, return slice of messages
	return messages, nil
}

func (r *ChatRepositoryImpl) UpdateMessageStatus(ctx context.Context, messageID int, newStatus entity.MessageStatus) error {
	// step 1: define query untuk update status pesan
	query := `UPDATE messages SET status = $1 WHERE id = $2`

	// step 2: jalankan query-nya (ExecContext) digunakan untuk eksekusi query yang tidak mengembalikan baris (UPDATE)
	result, err := r.DB.ExecContext(ctx, query, newStatus, messageID)
	if err != nil {
		return fmt.Errorf("error updating message status for %d: %w", messageID, err)
	}

	// step 3: cek apakah ada baris yang terpengaruh (RowsAffected) untuk memastikan pesan dengan ID tersebut ada
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected for message %d: %w", messageID, err)
	}

	// step 4: jika tidak ada baris yang terpengaruh, berarti pesan dengan ID tersebut tidak ditemukan
	if rowsAffected == 0 {
		return fmt.Errorf("no message found with ID %d to update status", messageID)
	}

	// step 5: jika berhasil, return nil
	return nil
}

func (r *ChatRepositoryImpl) GetUnreadMessagesForUser(ctx context.Context, userID int, afterTimestamp time.Time) ([]*entity.Message, error) {
	// step 1: define query untuk mengambil pesan yang statusnya bukan read dan setelah timestamp tertentu
	query := `
	SELECT id, senderid, recipientid, content, timestamp, status
	FROM messages
	WHERE (recipientid = $1 AND status != $2) AND timestamp > $3
	ORDER BY timestamp ASC
	`

	// step 2: jalankan query-nya (QueryContext) digunakan untuk eksekusi query yang mengembalikan banyak baris (SELECT)
	rows, err := r.DB.QueryContext(ctx, query, userID, entity.StatusRead, afterTimestamp)
	if err != nil {
		return nil, fmt.Errorf("error retrieving unread messages for user %d: %w", userID, err)
	}

	// step 3: pastikan untuk menutup rows setelah selesai digunakan
	defer rows.Close()

	// step 4: buat slice untuk menyimpan pesan yang diambil dari database
	var messages []*entity.Message
	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.RecipientID, &msg.Content, &msg.Timestamp, &msg.Status); err != nil {
			fmt.Printf("error scanning unread message: %v\n", err)
			continue
		}
		messages = append(messages, &msg)
	}

	// step 5: periksa apakah ada error saat iterasi
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over unread messages for user %d: %w", userID, err)
	}

	// step 6: jika berhasil, return slice of unread messages
	return messages, nil
}