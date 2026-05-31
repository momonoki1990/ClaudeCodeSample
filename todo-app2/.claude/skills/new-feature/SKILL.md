---
name: new-feature
description: develop ブランチを最新化し、引数で指定された名前で feature/<name> ブランチを作成する。引数が未指定なら名前を尋ねる。
---

# new-feature

新しい feature ブランチを作るときの定型作業（fetch → develop 更新 → feature 切り出し）を一括で実行する。

## 引数

`$ARGUMENTS` に feature 名（kebab-case 推奨、`feature/` プレフィックスは付けない）。例:

- `/new-feature create-todo` → `feature/create-todo` を作成
- `/new-feature list-todos` → `feature/list-todos` を作成

引数が空のときは、ユーザーに名前を尋ねてから進める。勝手に決めない。

## 手順

以下を順に実行する。途中で失敗したら即座に停止し、原因をユーザーに報告すること。

### 1. 引数のバリデーション

- 引数が空 → ユーザーに「feature 名を教えてください（kebab-case 推奨）」と尋ねる。
- 引数に `feature/` が含まれている → ユーザーに警告し、プレフィックスを外した名前で進めてよいか確認。
- 引数に空白や `/` が含まれる、または英数字とハイフン以外を含む → ユーザーに確認。

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
