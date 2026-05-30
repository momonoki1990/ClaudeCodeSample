# Claude Code Workflow Guide

## プロジェクト概要

- **プロジェクト名**: TODO App API
- **言語**: Rust
- **目的**: TODO リストを作成・管理できる API バックエンド
- **開発ワークフロー**: Git Flow + Claude Code

---

## Claude の役割と自律実行範囲

権限は [.claude/settings.json](.claude/settings.json) の `permissions` で制御する。
CLAUDE.md は方針、settings.json はハード制約という役割分担。

### 自律実行（承認不要 / settings.json の allow）

- `feature/*` ブランチの作成・commit・push
- `cargo build / test / fmt / clippy` の実行
- ファイル編集（Edit / Write / Read）
- PR の作成（`gh pr create`）

### 必ず人間と対話してから実施

- **実装方針の決定**: プランモード（`/plan`）で方針を提示し、人間が承認してから実装に着手
- **設計上の判断**: 仕様が不明瞭な場合は実装を止めて確認

### 禁止（settings.json の deny で強制）

- `main` / `develop` ブランチへの直接 push・切り替え
- `--force` push（`git push --force` / `-f`）
- `.env` / `*.pem` / `*.key` など秘密情報の読み取り
- 本番環境への deploy

---

## 開発フロー（全体像）

```
1. Planning（対話）   ← 人間が承認するまで実装しない
      ↓
2. 実装（自律）       ← feature/* で実装
      ↓
3. /code-review       ← Claude がローカルで自己レビュー → 修正
      ↓
4. commit             ← pre-commit が cargo fmt + clippy を自動実行
      ↓
5. push（承認不要）   ← feature/* へ push
      ↓
6. PR 作成            ← gh pr create で develop 向け PR
      ↓
7. CI（GitHub）       ← CodeRabbit(AIレビュー) + GitHub Actions(clippy/test/build)
      ↓
8. 人間レビュー       ← 最終確認 → develop へマージ
```

### 1. Planning（対話）

```
人間: 「〇〇機能を実装したい」
Claude: /plan で実装方針を提示（フレームワーク・データ構造・エラー形式など）
人間: 承認 or フィードバック
```

Claude が止まるのはこの工程のみ。承認後の 2〜6 は自律実行する。

### 2. 実装（自律）

```bash
git checkout -b feature/xxx develop
# コード実装
```

### 3. ローカル AI レビュー（push 前）

```
/code-review   # Claude が差分を自己レビューし、明らかな問題を修正
```

### 4〜6. commit / push / PR

```bash
git commit -m "feat: ..."        # pre-commit が fmt + clippy を実行
git push origin feature/xxx      # 承認不要
gh pr create --base develop --title "..." --body "..."
```

### 7. CI（GitHub 上の自動チェック）

- **CodeRabbit**: PR に AI レビューコメントを自動投稿（第三者視点）
- **GitHub Actions**: `cargo clippy` / `cargo test` / `cargo build` の品質ゲート

ローカル(Claude) と CI(CodeRabbit) でレビュー主体を変え、視点の偏りを防ぐ。

### 8. 人間レビュー → マージ

CI 全パス + AI レビュー確認後、人間が `develop` へマージする。

---

## ブランチ戦略

```
main (本番リリース)        ← 人間がマージ
  ↑
develop (開発統合)         ← PR 経由でマージ（CI 必須）
  ↑
feature/* (機能開発)       ← Claude が作業・push（承認不要）
  例: feature/create-todo, feature/list-todos
```

---

## コード品質保証

### Pre-commit Hook（フォーマット + Lint）

commit 時に自動実行。fmt エラーまたは clippy 警告があれば commit をブロック:

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: cargo-fmt
        name: cargo fmt
        entry: cargo fmt --
        language: system
        types: [rust]
        pass_filenames: false
      - id: cargo-clippy
        name: cargo clippy
        entry: cargo clippy -- -D warnings
        language: system
        types: [rust]
        pass_filenames: false
```

### GitHub Actions（Lint / Test / Build）

push・PR 時に自動実行:

```yaml
# .github/workflows/ci.yml
- run: cargo clippy -- -D warnings   # Lint（厳格モード）
- run: cargo test                     # テスト
- run: cargo build --release          # ビルド確認
```

CI 失敗時は PR がマージ不可。

### AI コードレビュー（2 段構え）

| 段階 | ツール | 場所 | タイミング |
|------|--------|------|-----------|
| 1次（自己） | Claude `/code-review` | ローカル | push 前 |
| 2次（第三者） | CodeRabbit | CI(GitHub App) | PR 作成時 |

人間レビューはこの 2 段の AI レビュー後に行う。

---

## 権限設定（.claude/settings.json）

詳細は [.claude/settings.json](.claude/settings.json) を参照。要点:

- **allow**: `cargo *`, `feature/*` への push, `gh pr create`, ファイル編集
- **deny**: `main`/`develop` への push・切替, `--force` push, 秘密情報の読み取り
- deny は allow より優先されるため、ガードレールが確実に効く

---

## 複数エージェントでの並行開発（将来オプション）

> 現状の TODO API 規模では**単一エージェントで直列開発**を推奨。
> 独立した機能が複数同時に走る段階になったら、以下で並行化する。

### git worktree による並行開発

各エージェントを別ディレクトリ・別ブランチに分離し、競合なく並行作業:

```bash
git worktree add ../todo-app2-create feature/create-todo
git worktree add ../todo-app2-list   feature/list-todos
```

各ディレクトリで独立した Claude Code セッションを起動する。

---

## TODO API 仕様

### エンドポイント

| メソッド | パス | 説明 |
|----------|------|------|
| `POST` | `/api/todos` | TODO 作成 |
| `GET` | `/api/todos` | 一覧取得 |
| `GET` | `/api/todos/{id}` | 詳細取得 |
| `PUT` | `/api/todos/{id}` | 更新 |
| `DELETE` | `/api/todos/{id}` | 削除 |

### フィールド

| フィールド | 型 | 必須 | デフォルト |
|-----------|-----|------|-----------|
| `id` | UUID | 自動採番 | - |
| `title` | String | ✅ | - |
| `description` | String | ❌ | null |
| `completed` | bool | ❌ | false |
| `created_at` | DateTime | 自動 | - |
| `updated_at` | DateTime | 自動 | - |

---

## セットアップで作成するファイル

- [x] `.claude/settings.json` - Claude 権限設定
- [ ] `Cargo.toml` - プロジェクト設定
- [ ] `src/main.rs` - エントリーポイント
- [ ] `.github/workflows/ci.yml` - GitHub Actions（Lint / Test / Build）
- [ ] `.pre-commit-config.yaml` - Pre-commit（`cargo fmt` + `cargo clippy`）
- [ ] `.gitignore` - git 除外設定
- [ ] CodeRabbit 連携（GitHub App インストール、任意で `.coderabbit.yaml`）

---

## Commit メッセージ規約

```
<type>: <subject>

<body>（任意）
```

| タイプ | 用途 |
|--------|------|
| `feat:` | 新機能 |
| `fix:` | バグ修正 |
| `refactor:` | リファクタリング |
| `test:` | テスト追加 |
| `docs:` | ドキュメント更新 |
| `chore:` | ビルド・依存関係の更新 |

---

## 実装前チェックリスト（Claude が確認）

- [ ] `feature/*` ブランチで作業しているか
- [ ] 実装方針を人間と合意済みか
- [ ] `/code-review` で自己レビュー済みか
- [ ] `cargo fmt` 済みか（pre-commit でも検証される）
- [ ] `cargo clippy -- -D warnings` が警告なしか（pre-commit でも検証される）
- [ ] `cargo test` が全パスか
- [ ] commit メッセージが規約に従っているか
