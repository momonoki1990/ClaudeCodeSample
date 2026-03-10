SET
  NAMES utf8mb4;

INSERT INTO
  categories (name, position, is_system)
VALUES
  ('仕事',         1, FALSE),
  ('プライベート',  2, FALSE),
  ('買い物',        3, FALSE);

INSERT INTO
  todos (text, done, category_id, position)
VALUES
  ('週次レポートを提出する', false, 2, 0),
  ('チームミーティングの準備をする', false, 2, 1),
  ('運動する', false, 3, 2),
  ('本を読む', false, 3, 3),
  ('牛乳を買う', false, 4, 4),
  ('野菜を買う', true, 4, 5),
  ('やることリストを整理する', false, NULL, 6);
