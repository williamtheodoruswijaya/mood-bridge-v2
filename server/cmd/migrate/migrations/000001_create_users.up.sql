CREATE TABLE IF NOT EXISTS users (
	UserID BIGSERIAL PRIMARY KEY,
	Username VARCHAR(50) NOT NULL UNIQUE,
	Fullname VARCHAR(50) NOT NULL,
	ProfileUrl VARCHAR(255),
	Email VARCHAR(50) NOT NULL UNIQUE,
    Password VARCHAR(255) NOT NULL,
	CreatedAt TIMESTAMP DEFAULT NOW()
);