package repository

import (
	"context"
	"database/sql"
	"mood-bridge-v2/server/internal/entity"
)

type PostRepository interface {
	Create(ctx context.Context, tx *sql.Tx, post *entity.Post) (*entity.Post, error)
	Find(ctx context.Context, db *sql.DB, postID int) (*entity.Post, error)
	FindAll(ctx context.Context, db *sql.DB) ([]*entity.Post, error)
	FindByUserID(ctx context.Context, db *sql.DB, postID int) ([]*entity.Post, error)
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

func (r *PostRepositoryImpl) Find(ctx context.Context, db *sql.DB, postID int) (*entity.Post, error) {
	query := `SELECT postid, userid, content, mood, createdat FROM posts WHERE postid = $1;`

	row := db.QueryRowContext(ctx, query, postID)

	var selectedPost entity.Post
	err := row.Scan(&selectedPost.PostID, &selectedPost.UserID, &selectedPost.Content, &selectedPost.Mood, &selectedPost.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Post not found
		}
		return nil, err // Other error
	}
	
	return &selectedPost, nil
}

func (r *PostRepositoryImpl) FindAll(ctx context.Context, db *sql.DB) ([]*entity.Post, error) {
	query := `SELECT postid, userid, content, mood, createdat FROM posts;`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.Mood, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostRepositoryImpl) FindByUserID(ctx context.Context, db *sql.DB, userID int) ([]*entity.Post, error) {
	query := `SELECT postid, userid, content, mood, createdat FROM posts WHERE userid = $1;`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var posts []*entity.Post
	for rows.Next() {
		var post entity.Post
		err := rows.Scan(&post.PostID, &post.UserID, &post.Content, &post.Mood, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}