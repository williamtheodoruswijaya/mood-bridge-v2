package repository

import (
	"context"
	"database/sql"
	"mood-bridge-v2/server/internal/entity"
)

type PostRepository interface {
	Create(ctx context.Context, tx *sql.Tx, post *entity.Post) (*entity.Post, error)
}

type PostRepositoryImpl struct {
}

func NewPostRepository() PostRepository {
	return &PostRepositoryImpl{}
}

func (r *PostRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, post *entity.Post) (*entity.Post, error) {
	query := `INSERT INTO posts (userid, content, mood) VALUES ($1, $2, $3) RETURNING postid, userid, content, mood, createdat`

	row := tx.QueryRowContext(ctx, query, post.UserID, post.Content, post.Mood)

	var createdPost entity.Post
	err := row.Scan(&createdPost.PostID, &createdPost.UserID, &createdPost.Content, &createdPost.Mood, &createdPost.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &createdPost, nil
}