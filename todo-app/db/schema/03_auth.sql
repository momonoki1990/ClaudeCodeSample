SET
  NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS users (
    id            INT          AUTO_INCREMENT PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at    DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE categories
    DROP INDEX name,
    ADD COLUMN user_id INT NULL AFTER id,
    ADD CONSTRAINT fk_categories_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    ADD UNIQUE KEY unique_user_category (user_id, name);

ALTER TABLE todos
    ADD COLUMN user_id INT NULL AFTER id,
    ADD CONSTRAINT fk_todos_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
