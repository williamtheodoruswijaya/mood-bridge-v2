package repository

import (
	"context"
	"database/sql"
	"mood-bridge-v2/server/internal/entity"
	"strconv"
)

type CommentRepository interface {
	Create(ctx context.Context, tx *sql.Tx, comment *entity.Comment) (*entity.Comment, error)
	GetAllByPostID(ctx context.Context, db *sql.DB, postID int) ([]*entity.Comment, error)
	Delete(ctx context.Context, tx *sql.Tx, commentID int) (string, error)
	GetByID(ctx context.Context, db *sql.DB, commentID int) (*entity.Comment, error)
}

type CommentRepositoryImpl struct {
}

func NewCommentRepository() CommentRepository {
	return &CommentRepositoryImpl{}
}

func (r *CommentRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, comment *entity.Comment) (*entity.Comment, error) {
	query := `INSERT INTO comments (postid, userid, content, createdat) VALUES ($1, $2, $3, $4) RETURNING commentid, postid, userid, content, createdat`

	row := tx.QueryRowContext(ctx, query, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)

	var createdComment entity.Comment
	err := row.Scan(&createdComment.CommentID, &createdComment.PostID, &createdComment.UserID, &createdComment.Content, &createdComment.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &createdComment, nil
}

func (r *CommentRepositoryImpl) GetAllByPostID(ctx context.Context, db *sql.DB, postID int) ([]*entity.Comment, error) {
	query := `SELECT commentid, postid, userid, content, createdat FROM comments WHERE postid = $1;`
	rows, err := db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*entity.Comment
	for rows.Next() {
		var comment entity.Comment
		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *CommentRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, commentID int) (string, error) {
	query := `DELETE FROM comments WHERE commentid = $1 RETURNING commentid`
	row := tx.QueryRowContext(ctx, query, commentID)

	var deletedCommentID int
	err := row.Scan(&deletedCommentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // Comment not found
		}
		return "", err // Other error
	}

	return "Comment with ID " + strconv.Itoa(deletedCommentID) + " deleted successfully", nil
}

func (r *CommentRepositoryImpl) GetByID(ctx context.Context, db *sql.DB, commentID int) (*entity.Comment, error) {
	query := `SELECT commentid, postid, userid, content, createdat FROM comments WHERE commentid = $1;`
	row := db.QueryRowContext(ctx, query, commentID)

	var selectedComment entity.Comment
	err := row.Scan(&selectedComment.CommentID, &selectedComment.PostID, &selectedComment.UserID, &selectedComment.Content, &selectedComment.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Comment not found
		}
		return nil, err // Other error
	}

	return &selectedComment, nil
}