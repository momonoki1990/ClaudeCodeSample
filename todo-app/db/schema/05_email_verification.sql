ALTER TABLE users
    ADD COLUMN email_verified TINYINT(1) NOT NULL DEFAULT 0 AFTER email;

CREATE TABLE IF NOT EXISTS email_verification_tokens (
    id         INT         AUTO_INCREMENT PRIMARY KEY,
    user_id    INT         NOT NULL,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    expires_at DATETIME    NOT NULL,
    created_at DATETIME    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_evt_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
