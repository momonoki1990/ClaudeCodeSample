CREATE TABLE IF NOT EXISTS categories (
    id        INT          AUTO_INCREMENT PRIMARY KEY,
    name      VARCHAR(255) NOT NULL UNIQUE,
    position  INT          NOT NULL DEFAULT 0,
    is_system BOOLEAN      NOT NULL DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS todos (
    id          INT     AUTO_INCREMENT PRIMARY KEY,
    text        TEXT    NOT NULL,
    done        BOOLEAN NOT NULL DEFAULT FALSE,
    category_id INT     NULL,
    position    INT     NOT NULL DEFAULT 0,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
);
