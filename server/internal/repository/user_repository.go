package repository

import (
	"context"
	"database/sql"
	"mood-bridge-v2/server/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
	Find(ctx context.Context, db *sql.DB, username string) (*entity.User, error)
}

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error) {
	// step 1: define query-nya
	query := `INSERT INTO users (username, fullname, email, password) VALUES ($1, $2, $3, $4) RETURNING userid, username, fullname, email, password, profileurl, createdat`

	// step 2: execute query-nya
	row := tx.QueryRowContext(ctx, query, user.Username, user.Fullname, user.Email, user.Password)

	// step 3: scan hasilnya ke dalam struct user untuk di return.
	var createdUser entity.User
	err := row.Scan(&createdUser.ID, &createdUser.Username, &createdUser.Fullname, &createdUser.Email, &createdUser.Password, &createdUser.ProfileUrl, &createdUser.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &createdUser, err
}

func (r *UserRepositoryImpl) Find(ctx context.Context, db *sql.DB, username string) (*entity.User, error) {
	// step 1: define query
	query := `SELECT userid, username, fullname, profileurl, email, password, createdat FROM users WHERE username = $1;`

	// step 2: execute query
	row := db.QueryRowContext(ctx, query, username)

	// step 3: scan row-nya ke struct user
	var selectedUser entity.User
	err := row.Scan(&selectedUser.ID, &selectedUser.Username, &selectedUser.Fullname, &selectedUser.ProfileUrl, &selectedUser.Email, &selectedUser.Password, &selectedUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	if !selectedUser.ProfileUrl.Valid {
		selectedUser.ProfileUrl.String = "https://upload.wikimedia.org/wikipedia/commons/a/ac/Default_pfp.jpg" // Default value if NULL
	}

	return &selectedUser, err
}
