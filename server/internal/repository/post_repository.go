package repository

import (
	"context"
	"database/sql"
	"mood-bridge-v2/server/internal/entity"
	"strconv"
)

type PostRepository interface {
	Create(ctx context.Context, tx *sql.Tx, post *entity.Post) (*entity.Post, error)
	Find(ctx context.Context, db *sql.DB, postID int) (*entity.Post, error)
	FindAll(ctx context.Context, db *sql.DB) ([]*entity.Post, error)
	FindByUserID(ctx context.Context, db *sql.DB, postID int) ([]*entity.Post, error)
	Update(ctx context.Context, tx *sql.Tx, postID int, post *entity.Post) (*entity.Post, error)
	Delete(ctx context.Context, tx *sql.Tx, postID int) (string, error)
	GetFriendPosts(ctx context.Context, db *sql.DB, userID int) ([]*entity.Post, error)
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
	query := `SELECT postid, userid, content, mood, createdat FROM posts ORDER BY createdat DESC;`
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

func (r *PostRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, postID int, post *entity.Post) (*entity.Post, error) {
	// set query-nya
	query := `UPDATE posts SET content = $1, mood = $2 WHERE postid = $3 RETURNING postid, userid, content, mood, createdat`

	// jalankan query-nya
	row := tx.QueryRowContext(ctx, query, post.Content, post.Mood, postID)

	// buat variable untuk menampung hasil query
	var updatedPost entity.Post

	// scan hasil query ke variable
	err := row.Scan(&updatedPost.PostID, &updatedPost.UserID, &updatedPost.Content, &updatedPost.Mood, &updatedPost.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Post not found
		}
		return nil, err // Other error
	}

	return &updatedPost, nil
}

func (r *PostRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, postID int) (string, error) {
	// step 1: set query-nya
	query := `DELETE FROM posts WHERE postid = $1 RETURNING postid`

	// step 2: jalankan query-nya
	row := tx.QueryRowContext(ctx, query, postID)

	// step 3: buat variable untuk menampung hasil query
	var deletedPostID int

	// step 4: scan hasil query ke variable
	err := row.Scan(&deletedPostID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Post not found
		}
		return "", err // Other error
	}

	// step 5: return hasilnya
	return "Post with ID " + strconv.Itoa(deletedPostID) + " deleted successfully", nil
}

func (r *PostRepositoryImpl) GetFriendPosts(ctx context.Context, db *sql.DB, userID int) ([]*entity.Post, error) {
	query := `
		SELECT p.postid, p.userid, p.content, p.mood, p.createdat
		FROM posts p
		WHERE 
    		p.userid = $1
    		OR p.userid IN (
        	SELECT
            	CASE 
                	WHEN userid = $1 THEN frienduserid
                	ELSE userid
            	END AS friend_userid
        	FROM friends
        	WHERE (userid = $1 OR frienduserid = $1)
        	AND friendstatus = TRUE
    	)
	ORDER BY p.createdat DESC;
	`

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