# CLAUDE.md - sumika

## プロジェクト概要

**sumika** は、AI駆動開発時代の個人開発者向け軽量プロジェクトハブ。
複数のプロジェクトを横断するダッシュボードと、AIコンテキストの可視化を核とした OSS。

> "Backstage, but for solo developers."

### 解決する課題
- 個人開発者が複数プロジェクトを並行して扱う際の「どこまでやったか思い出す時間」の削減
- プロジェクトごとに散在するリポジトリ・ドキュメント・ローカル環境・AIコンテキストの一元管理
- 1コマンドでエディタ/Claude Codeを該当プロジェクトで起動する体験の実現

### 既存OSSとの差別化
| ツール | 領域 | sumikaとの違い |
|---|---|---|
| Backstage | 企業向けポータル | Kubernetes前提で個人には重すぎる |
| ctx-cli | 作業環境の切り替え | ダッシュボード・AI可視化なし |
| cli-continues | AIセッション引き継ぎ | プロジェクト管理・環境起動なし |
| ProjectHub-Mcp | タスク管理 | Docker/Postgres前提でエンタープライズ寄り |
| shiori | ブックマーク管理 | 開発プロジェクト管理ではない |

---

## 技術スタック

| 領域 | 選定 | 備考 |
|---|---|---|
| 言語 | Go | 単一バイナリ配布、クロスプラットフォーム |
| CLIフレームワーク | `cobra` | サブコマンド・フラグ管理 |
| Webサーバー | `net/http` + `chi` | 軽量・net/http互換 |
| フロントエンド | HTMX + `html/template` | npmレス・`embed`で単一バイナリに同梱 |
| 設定管理 | YAML(`gopkg.in/yaml.v3`) | 人間が読み書きしやすい |
| プロセス管理 | `os/exec` | Claude Code・各プロジェクト起動コマンド実行 |
| Git情報取得 | `os/exec`でgit CLI呼び出し | 依存最小・確実 |
| データ永続化 | YAMLファイル(将来SQLite移行可) | `~/.config/sumika/config.yaml`に一元管理 |
| 配布 | GitHub Releases + GoReleaser | クロスコンパイル・リリース自動化 |

---

## アーキテクチャ方針

### レイヤー構成
```
┌─────────────────────────────────┐
│  Layer 3: 推奨統合               │
│  Claude Code特化機能(起動・CLAUDE.md管理・セッション表示) │
├─────────────────────────────────┤
│  Layer 2: アダプター層           │
│  AIツールアダプター(claude-code / cursor / aider ...) │
├─────────────────────────────────┤
│  Layer 1: コア機能(ツール非依存) │
│  プロジェクト管理・起動・git状態・ダッシュボード │
└─────────────────────────────────┘
```

### 設計原則
- **コア機能はツール非依存**: Layer 1はClaude Code以外のユーザーでも価値が成立する
- **Claude Codeは最初の公式アダプター**: Layer 2に`claude-code`アダプターを実装し、将来的に他ツールも追加可能な設計
- **シングルバイナリ**: フロントエンドは`embed`パッケージで同梱し、外部依存なしで動作
- **設定ファイルはGit管理可能**: `~/.config/sumika/config.yaml`をdotfilesで管理できる設計

### ディレクトリ構成
```
sumika/
├── cmd/
│   └── sumika/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── config/              # 設定ファイルの読み書き
│   ├── project/             # プロジェクト管理ロジック
│   ├── git/                 # git情報取得
│   ├── launcher/            # エディタ・Claude Code起動
│   ├── adapter/             # AIツールアダプター
│   │   ├── adapter.go       # インターフェース定義
│   │   └── claudecode/      # Claude Codeアダプター実装
│   └── server/              # Webサーバー・ハンドラー
├── web/
│   ├── templates/           # html/templateテンプレート
│   └── static/              # CSS・JS(HTMX)
├── go.mod
├── go.sum
├── .goreleaser.yaml
└── CLAUDE.md                # このファイル
```

---

## CLIコマンド体系

```
sumika init                  # 設定ファイルを初期化(~/.config/sumika/config.yaml作成)
sumika add <path>            # プロジェクトを登録
sumika list                  # プロジェクト一覧を表示
sumika open <name>           # エディタ + Claude Codeを該当ディレクトリで起動
sumika serve                 # Webダッシュボードを起動(デフォルト: localhost:8964)
sumika status                # 全プロジェクトのgit status・起動状態を一覧表示
sumika remove <name>         # プロジェクトを登録解除
```

---

## 設定ファイル仕様

**配置場所**: `~/.config/sumika/config.yaml`

```yaml
version: 1

settings:
  port: 8964
  editor: code          # 起動するエディタコマンド(code / cursor / nvim 等)
  ai_tool: claude       # 起動するAIツールコマンド(claude / aider 等)

projects:
  - name: my-api
    path: ~/projects/my-api
    description: "REST API サーバー"
    tags:
      - backend
      - go
    launch:
      editor: true      # sumika open時にエディタを起動するか
      ai: true          # sumika open時にAIツールを起動するか
      commands:         # 追加で実行するコマンド(開発サーバー起動等)
        - "docker compose up -d"
    links:
      - label: "仕様書"
        url: "https://notion.so/..."
      - label: "staging"
        url: "https://staging.example.com"

  - name: my-frontend
    path: ~/projects/my-frontend
    description: "Next.js フロントエンド"
    tags:
      - frontend
      - typescript
    launch:
      editor: true
      ai: true
      commands:
        - "npm run dev"
```

---

## AIアダプターインターフェース

```go
// internal/adapter/adapter.go
package adapter

type SessionSummary struct {
    ProjectName string
    LastActive  time.Time
    Summary     string    // セッションの最後の作業内容
    RawLog      string    // 生ログ
}

type AIAdapter interface {
    Name() string
    IsAvailable() bool                          // コマンドがPATHに存在するか
    Launch(projectPath string) error            // AIツールを起動
    GetSessionSummary(projectPath string) (*SessionSummary, error)  // セッション履歴取得
    GetContextFile(projectPath string) (string, error)              // CLAUDE.md等の取得
}
```

---

## Claude Codeアダプターの実装方針

### 起動
```go
// os/execでサブプロセスとして起動
cmd := exec.Command("claude")
cmd.Dir = projectPath
cmd.Stdin = os.Stdin
cmd.Stdout = os.Stdout
cmd.Stderr = os.Stderr
```

### セッション履歴の取得
- `~/.claude/projects/` 配下のJSONLファイルを読み込む
- ファイルパスはプロジェクトのパスをハッシュ化したディレクトリ名で対応
- 最終セッションの末尾N件のメッセージからサマリーを生成
- ⚠️ **注意**: `~/.claude/`配下のファイル形式はClaude Codeの内部仕様であり、バージョンアップで変更される可能性がある。変更に備えてアダプター層に封じ込め、バージョン検出ロジックを入れる

### CLAUDE.mdの管理
- 各プロジェクトルートの`CLAUDE.md`を読み込み、ダッシュボードに表示
- 複数プロジェクト間で共通ルールを管理する「グローバルCLAUDE.md」の同期機能(将来)

---

## ダッシュボードUI仕様

### プロジェクトカード(1プロジェクト = 1カード)
```
┌──────────────────────────────────────────┐
│ 🟢 my-api                    [Open] [↗]  │
│ REST API サーバー                         │
│                                           │
│ 📁 ~/projects/my-api                     │
│ 🕐 最終更新: 2時間前                      │
│ 📝 git: 3件の未コミット変更               │
│ 🤖 前回AI: 認証機能のリファクタリング中   │
│    テスト未完了                           │
│                                           │
│ 🏷️ backend  go                           │
└──────────────────────────────────────────┘
```

### ダッシュボード全体
- プロジェクト一覧をカードグリッドで表示
- タグでフィルタリング可能
- 「最終更新順」「プロジェクト名順」でソート
- 各カードの[Open]ボタンでエディタ+Claude Codeを起動
- HTMXによるポーリング(30秒ごと)でgit status・起動状態を自動更新

---

## 実装優先順位(MVP)

### Phase 1: コアCLI
- [x] `sumika init` — config.yaml生成
- [x] `sumika add <path>` — プロジェクト登録(`--name`, `--description`フラグ対応)
- [x] `sumika list` — プロジェクト一覧表示(ターミナル)
- [x] `sumika open <name>` — エディタ起動(Claude Code含む)
- [x] `sumika status` — git status一覧表示
- [x] `sumika remove <name>` — プロジェクト登録解除

### Phase 2: Webダッシュボード
- [x] `sumika serve` — Webサーバー起動(HTMX + html/template、embed同梱)
- [x] プロジェクトカード一覧の表示
- [x] git status・最終更新日時の表示(相対時刻)
- [x] [Open]ボタンでの起動連携

### Phase 3: AIコンテキスト可視化
- [ ] Claude Codeアダプターの実装
- [ ] `~/.claude/`配下のセッション履歴読み込み
- [ ] 各プロジェクトカードへの「前回AIセッション要約」表示
- [ ] CLAUDE.mdの内容表示・最終更新日時表示

### Phase 4: 拡張
- [ ] タグフィルタリング・ソート
- [ ] グローバルCLAUDE.mdの同期機能
- [ ] 他AIツールアダプター(cursor, aider等)
- [ ] SQLite移行(プロジェクト数が多い場合の性能対策)
- [ ] `sumika serve`の自動起動(launchd/systemd連携)

---

## 開発環境セットアップ

```bash
# リポジトリのクローン
git clone https://github.com/mskasa/sumika
cd sumika

# 依存関係のインストール
go mod download

# kizami(ドキュメント管理CLI)のインストール
go install github.com/mskasa/kizami@latest

# 開発用ビルド
go build -o sumika ./cmd/sumika

# テスト実行
go test ./...

# ローカルインストール
go install ./cmd/sumika
```

---

## ドキュメント管理

ドキュメントは **[kizami](https://github.com/mskasa/kizami)** を使ってsumikaリポジトリ内の `docs/decisions/` で管理する。

kizamiはADR・設計ドキュメントをコードと並走して管理するための自作CLIツール。
ドキュメントとコードの乖離(drift)を自動検出する機能を持つ。

```bash
# インストール
go install github.com/mskasa/kizami@latest
```

### ドキュメントの配置

```
sumika/
└── docs/
    ├── decisions/              # ADR置き場(kizami adr で作成)
    │   ├── YYYY-MM-DD-use-go-over-other-languages.md
    │   ├── YYYY-MM-DD-use-htmx-html-template-over-spa-frameworks.md
    │   ├── YYYY-MM-DD-use-chi-as-http-router.md
    │   ├── YYYY-MM-DD-adopt-adapter-pattern-for-ai-tools.md
    │   └── YYYY-MM-DD-use-json-file-for-persistence-sqlite-later.md
    └── design/                 # 設計ドキュメント置き場(kizami design で作成)
        ├── YYYY-MM-DD-dashboard-ui-specification.md
        ├── YYYY-MM-DD-cli-command-design.md
        ├── YYYY-MM-DD-config-file-schema.md
        └── YYYY-MM-DD-claude-code-adapter-implementation.md
```

### 基本的な使い方

```bash
# インストール(v0.9.5以降)
go install github.com/mskasa/kizami@latest

# 初期化(docs/decisions/, docs/design/ディレクトリを生成)
kizami init

# ADRを作成($EDITORで開く)
kizami adr "use Go over other languages"

# 設計ドキュメントを作成
kizami design "dashboard UI specification"

# 一覧表示
kizami list

# キーワード検索
kizami search "adapter"

# ステータス更新(slug指定)
kizami status use-go-over-other-languages accepted
kizami status use-go-over-other-languages superseded --by use-rust

# コードとの乖離チェック(Related Filesに記載したファイルが存在するか検証)
kizami audit
```

### ドキュメントのステータス

| ステータス | 意味 |
|---|---|
| `Draft` | 作成直後の初期状態 |
| `Accepted` | 現在有効な決定 |
| `Deprecated` | 廃止、後継なし |
| `Superseded` | 別の決定に置き換え済み(`--by <slug>`で後継を指定) |

### ドキュメントフォーマット

kizamiが生成するMarkdownテンプレート:

**ADR** (`kizami adr`):
```markdown
# [タイトル]

- Date: YYYY-MM-DD
- Type: ADR
- Status: Draft
- Author: mskasa

## Context
## Decision
## Consequences
## Alternatives Considered
## Related Files
```

**設計ドキュメント** (`kizami design`):
```markdown
# [タイトル]

- Date: YYYY-MM-DD
- Type: Design
- Status: Draft
- Author: mskasa

## Overview
## Background
## Goals / Non-Goals
## Design
## Implementation Plan
## Open Questions
## Related Files
```

### 初期ADR一覧(作成済み)

| タイトル | 種別 | ステータス |
|---|---|---|
| use Go over other languages | ADR | Active |
| use HTMX + html/template over SPA frameworks | ADR | Active |
| use chi as HTTP router | ADR | Active |
| adopt adapter pattern for AI tools | ADR | Active |
| use JSON file for persistence, SQLite later | ADR | Active |
| dashboard UI specification | design | Draft |
| CLI command design | design | Draft |
| config file schema | design | Draft |
| Claude Code adapter implementation | design | Draft |

### ドキュメント更新フロー

```
1. 新機能・仕様変更の検討
   ↓
2. ADRが必要な意思決定 → kizami adr "<title>"
   設計仕様の記述     → kizami design "<title>"
   ↓
3. sumikaリポジトリで実装
   (Related Filesに実装ファイルのパスを記載)
   ↓
4. kizami audit でdriftがないか確認
   ↓
5. 大きな変更の場合はCLAUDE.mdも更新
```

### drift検出について

`## Related Files` セクションに実装ファイルのパスを記載しておくと、
`kizami audit` がファイルの削除・移動を検出してドキュメントの陳腐化を防ぐ。

```bash
# 例: リファクタリングでファイルを移動した後に実行
kizami audit
# → Related Filesに記載されたファイルが存在しないADRを一覧表示
```

---

## コーディング規約

- **エラーハンドリング**: `fmt.Errorf("context: %w", err)` でラップし、呼び出し元でハンドリング
- **ログ出力**: `log/slog`を使用(Go 1.21以降)
- **テスト**: テーブル駆動テストを基本とする
- **コミットメッセージ**: Conventional Commits準拠(`feat:`, `fix:`, `docs:` 等)
- **ブランチ戦略**: `main`ブランチへのPRベース開発

---

## ライセンス

MIT License

---

## TODO(未決定事項)

- [ ] ロゴ・タグラインの作成
- [ ] READMEの初稿(コンセプト・既存OSSとの差別化・インストール手順)
- [ ] CONTRIBUTING.md・Issue/PRテンプレートの作成
