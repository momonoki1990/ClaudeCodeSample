---
name: new-feature
description: develop を最新化して feature/<name> ブランチを作成する定型作業（fetch → develop pull → feature 切り出し）を実行する。Use when ユーザーが新しい機能開発・バグ修正・リファクタなどの作業を開始しようとしているとき。例:「〇〇機能を実装したい」「△△のバグを直す」「新しい feature ブランチを作って」「ブランチ切って」「develop から派生させて」など。ブランチ作成までに責務を限定し、雛形ファイル生成や実装そのものには踏み込まない。
---

# new-feature

新しい feature ブランチを作るときの定型作業（fetch → develop 更新 → feature 切り出し）を一括で実行する。

## 起動と引数

このスキルは 2 つの経路で起動する:

1. **ユーザーが `/new-feature <name>` で明示起動**: `$ARGUMENTS` に feature 名が入る。
2. **Claude が文脈から自律起動**: ユーザーの自然文発話（例:「create-todo 機能を作りたい」「ログイン周りのバグ修正始める」）から feature 名を推測する。推測が曖昧／複数解釈がある場合は必ずユーザーに確認する。**勝手に決めない**。

feature 名は kebab-case 推奨、`feature/` プレフィックスは付けない。例:

- `create-todo` → `feature/create-todo` を作成
- `list-todos` → `feature/list-todos` を作成

### 自律起動の判断基準

- ✅ 起動する: ユーザーが新規作業の開始を示唆し、かつ現在 `feature/*` ブランチ上にいない、または別テーマの作業に切り替える文脈
- ❌ 起動しない: すでに `feature/*` ブランチ上で進行中の作業の続き／単発の質問・調査依頼／既存ブランチの修正依頼

迷ったら起動せず、「新しい feature ブランチを切りますか？」と一言確認する。

## 手順

以下を順に実行する。途中で失敗したら即座に停止し、原因をユーザーに報告すること。

### 1. 名前のバリデーション

- 引数が空、または自律起動で名前を確信できない → ユーザーに「feature 名を教えてください（kebab-case 推奨）」と尋ねる。
- 名前に `feature/` が含まれている → ユーザーに警告し、プレフィックスを外した名前で進めてよいか確認。
- 名前に空白や `/` が含まれる、または英数字とハイフン以外を含む → ユーザーに確認。
- Claude が文脈から名前を推測した場合は、ブランチ作成前に「`feature/<推測した名前>` で進めます。問題なければそのまま作成します」と一言告げ、明確な反対がなければ続行する。

### 2. 作業ツリーが clean か確認

```bash
git status --porcelain
```

未コミットの変更があれば実装を停止し、ユーザーに対処方法を尋ねる（stash する／commit する／中断する のいずれか）。**勝手に stash や reset をしてはいけない。**

### 3. 既存ブランチとの衝突チェック

```bash
git rev-parse --verify --quiet refs/heads/feature/<name>
```

既に存在する場合は停止し、ユーザーに「別名にする／既存ブランチに切り替える」を尋ねる。

### 4. リモート最新化と develop の更新

```bash
git fetch origin
git checkout develop
git pull origin develop
```

`pull` で conflict が出た場合は停止して報告（develop の状態が想定外なので人間判断が必要）。

### 5. feature ブランチを作成

```bash
git checkout -b feature/<name>
```

### 6. 最終確認の出力

以下を 1 ブロックで報告する:

- 作成したブランチ名
- ベースにした develop の最新 commit（`git log -1 --oneline` の結果）
- 次にやること（実装方針を `/plan` で詰めるか、すぐ実装に入るか）の簡単な提案

## やってはいけないこと

- `main` への切り替え（settings.json で deny されている）
- `--force` 系の操作
- ユーザー確認なしの `git stash` / `git reset` / `git clean`
- 「ついでに」スキャフォルディング（雛形ファイル生成）まで踏み込むこと。このスキルはブランチ作成までに責務を限定する。
