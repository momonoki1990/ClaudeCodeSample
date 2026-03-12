CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         INT         AUTO_INCREMENT PRIMARY KEY,
    user_id    INT         NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at DATETIME    NOT NULL,
    created_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_refresh_tokens_expires_at (expires_at)
);
