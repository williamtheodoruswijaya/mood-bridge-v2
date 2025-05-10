package db

import "database/sql"

func Migrate(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			UserID SERIAL PRIMARY KEY,
			Username VARCHAR(50) NOT NULL UNIQUE,
			Fullname VARCHAR(50) NOT NULL,
			ProfileUrl VARCHAR(255),
			Email VARCHAR(50) NOT NULL UNIQUE,
			Password VARCHAR(255) NOT NULL,
			CreatedAt TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS friends (
			FriendID SERIAL PRIMARY KEY,
			UserID INTEGER REFERENCES users(UserID) ON DELETE CASCADE,
			FriendUserID INTEGER REFERENCES users(UserID) ON DELETE CASCADE,
			FriendStatus VARCHAR(50) DEFAULT 'PENDING',
			CreatedAt TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS posts (
			PostID SERIAL PRIMARY KEY,
			UserID INTEGER REFERENCES users(UserID) ON DELETE CASCADE,
			Content TEXT NOT NULL,
			Mood VARCHAR(50) DEFAULT 'NORMAL',
			CreatedAt TIMESTAMP DEFAULT NOW()
		);`,
		`CREATE TABLE IF NOT EXISTS comments (
			CommentID SERIAL PRIMARY KEY,
			PostID INTEGER REFERENCES posts(PostID) ON DELETE CASCADE,
			UserID INTEGER REFERENCES users(UserID) ON DELETE CASCADE,
			Content TEXT NOT NULL,
			CreatedAt TIMESTAMP DEFAULT NOW()
		);`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			panic(err)
		}
	}
}
