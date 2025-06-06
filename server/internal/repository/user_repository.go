package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mood-bridge-v2/server/internal/entity"
	"strings"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user *entity.User) (*entity.User, error)
	Find(ctx context.Context, db *sql.DB, username string) (*entity.User, error)
	FindByID(ctx context.Context, db *sql.DB, id int) (*entity.User, error)
	FindByEmail(ctx context.Context, db *sql.DB, email string) (*entity.User, error) // for login and validation
	FindAll(ctx context.Context, db *sql.DB) ([]*entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, id int, user *entity.User) (*entity.User, error)
	FindByIDs(ctx context.Context, db *sql.DB, ids []int) ([]*entity.User, error)
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
		selectedUser.ProfileUrl.String = defaultProfileUrl // Default value if NULL
	}

	return &selectedUser, err
}

func (r *UserRepositoryImpl) FindByID(ctx context.Context, db *sql.DB, id int) (*entity.User, error) {
	query := `SELECT userid, username, fullname, profileurl, email, password, createdat FROM users WHERE userid = $1;`
	row := db.QueryRowContext(ctx, query, id)
	var selectedUser entity.User
	err := row.Scan(&selectedUser.ID, &selectedUser.Username, &selectedUser.Fullname, &selectedUser.ProfileUrl, &selectedUser.Email, &selectedUser.Password, &selectedUser.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	if !selectedUser.ProfileUrl.Valid {
		selectedUser.ProfileUrl.String = defaultProfileUrl // Default value if NULL
	}
	return &selectedUser, err
}

func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, db *sql.DB, email string) (*entity.User, error) {
	// step 1: define query
	query := `SELECT userid, username, fullname, profileurl, email, password, createdat FROM users WHERE email = $1;`

	// step 2: execute query
	row := db.QueryRowContext(ctx, query, email)

	// step 3: scan row-nya ke struct user
	var selectedUser entity.User
	err := row.Scan(&selectedUser.ID, &selectedUser.Username, &selectedUser.Fullname, &selectedUser.ProfileUrl, &selectedUser.Email, &selectedUser.Password, &selectedUser.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) { // more informative error handling (saat ini dibiarin dulu biar saya bisa ngeliat errornya)
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	if !selectedUser.ProfileUrl.Valid {
		selectedUser.ProfileUrl.String = defaultProfileUrl // Default value if NULL
	}

	return &selectedUser, err
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context, db *sql.DB) ([]*entity.User, error) {
	// step 1: define query
	query := `SELECT userid, username, fullname, profileurl, email, password, createdat FROM users;`

	// step 2: execute query
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // close rows after use buat mencegah memory leak

	// step 3: scan row-nya ke struct user
	var users []*entity.User
	for rows.Next() {
		// Create a new User instance for each row and scan the values into it
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Fullname, &user.ProfileUrl, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Check if ProfileUrl is NULL and set default value
		if !user.ProfileUrl.Valid {
			user.ProfileUrl.String = defaultProfileUrl // Default value if NULL
		}

		// Append the user to the slice
		users = append(users, &user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// step 4: return the slice of users
	return users, nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, id int, user *entity.User) (*entity.User, error) {
	// step 1: define query-nya
	query := `UPDATE users SET username = $1, fullname = $2, email = $3, password = $4, profileurl = $5 WHERE userid = $6 RETURNING userid, username, fullname, email, password, profileurl, createdat`

	// step 2: execute query-nya
	row := tx.QueryRowContext(ctx, query, user.Username, user.Fullname, user.Email, user.Password, user.ProfileUrl.String, id)

	// step 3: scan hasilnya ke dalam struct user untuk di return.
	var updatedUser entity.User
	err := row.Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.Fullname, &updatedUser.Email, &updatedUser.Password, &updatedUser.ProfileUrl, &updatedUser.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &updatedUser, err
}

func (r *UserRepositoryImpl) FindByIDs(ctx context.Context, db *sql.DB, ids []int) ([]*entity.User, error) {
	if len(ids) == 0 {
		return []*entity.User{}, nil
	}

	// step 1: buat placeholder untuk query (sesuai jumlah ids)
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	// step 2: buat query-nya
	query := fmt.Sprintf("SELECT userid, username, fullname, profileurl, email, password, createdat FROM users WHERE userid IN (%s)", strings.Join(placeholders, ","))

	// step 3: execute query-nya
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // close rows after use buat mencegah memory leak

	// step 4: scan row-nya ke struct user
	var users []*entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.Username, &user.Fullname, &user.ProfileUrl, &user.Email, &user.Password, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		if !user.ProfileUrl.Valid {
			user.ProfileUrl.String = defaultProfileUrl // Default value if NULL
		}
		users = append(users, &user)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// step 5: return the slice of users
	return users, nil
}

const defaultProfileUrl = "https://upload.wikimedia.org/wikipedia/commons/a/ac/Default_pfp.jpg"