# Todo App

JWT認証・複数ユーザー対応のTodoアプリです。

## 構成

| サービス       | URL                    | 説明                      |
| -------------- | ---------------------- | ------------------------- |
| フロントエンド | https://localhost:5173 | React + TypeScript (Vite) |
| バックエンド   | http://localhost:8080  | Go (Echo)                 |
| DB             | localhost:3306         | MySQL 8.0                 |

## セットアップ

### 1. 環境変数の設定

```bash
cp .env.sample .env
```

必要に応じて `.env` の値を編集してください（特に `JWT_SECRET`）。

### 2. mkcert のインストールと初期化

ローカルHTTPS用の証明書を発行するために必要です。

```bash
# mkcert のインストール（Mac）
brew install mkcert

# ローカルCAをOSとブラウザの信頼リストに登録（初回のみ）
mkcert -install
```

### 3. 証明書の発行

```bash
mkdir -p front/certs
mkcert \
  -cert-file front/certs/localhost.pem \
  -key-file  front/certs/localhost-key.pem \
  localhost
```

### 4. 起動

```bash
mise run up
```

ブラウザで https://localhost:5173 にアクセスしてください。

## 開発

### DB の再作成

スキーマを変更した場合はボリュームごと削除して再作成してください。

```bash
mise run restart
```
