SET
    NAMES utf8mb4;

-- テストユーザー
-- email: test@example.com / password: password123
INSERT INTO
    users (email, email_verified, password_hash)
VALUES
    (
        'test@example.com',
        1,
        '$2a$10$vGSoAdPqLwV5hlvxF4dulOWMzwxclLYO1.mSZHI/W7zRQRAhUjRP2'
    );

-- カテゴリ（user_id=1）
INSERT INTO
    categories (user_id, name, position, is_system)
VALUES
    (1, 'すべて', 0, TRUE),
    (1, '仕事', 1, FALSE),
    (1, '個人', 2, FALSE);

-- サンプルタスク（user_id=1）
-- category_id: NULL=未分類, 2=仕事, 3=個人
INSERT INTO
    todos (user_id, text, done, category_id, position)
VALUES
    (1, 'Claude Codeを試してみる', TRUE,  NULL, 0),
    (1, 'タスク管理アプリを作る',  TRUE,  NULL, 1),
    (1, 'メール認証を実装する',    TRUE,  NULL, 2),
    (1, '機能追加を検討する',      FALSE, NULL, 3),
    (1, '企画書を作成する',        FALSE, 2,    0),
    (1, 'ミーティングの準備',      FALSE, 2,    1),
    (1, '請求書を送付する',        TRUE,  2,    2),
    (1, '本を読む',                FALSE, 3,    0),
    (1, '運動する',                FALSE, 3,    1),
    (1, '部屋を片付ける',          TRUE,  3,    2);