CREATE TABLE IF NOT EXISTS messages (
    MessageID BIGSERIAL PRIMARY KEY,
    SenderID BIGINT REFERENCES users(UserID) ON DELETE CASCADE,
    RecipientID BIGINT REFERENCES users(UserID) ON DELETE CASCADE,
    Content TEXT NOT NULL,
    Timestamp TIMESTAMP DEFAULT NOW(),
    Status VARCHAR(50) DEFAULT 'sent' -- Possible values: 'sent', 'delivered', 'read', 'failed' [cite: 50]
);