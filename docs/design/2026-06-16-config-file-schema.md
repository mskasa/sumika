# config file schema

- Date: 2026-06-16
- Type: Design
- Status: Draft
- Author: masahiro.kasatani

## Overview

sumikaの設定ファイル(`~/.config/sumika/config.yaml`)のスキーマを定義する。
人間が読み書きしやすいYAML形式で、dotfilesリポジトリでのGit管理を想定した設計。

## Background

プロジェクトの登録情報・グローバル設定を永続化する場所が必要。
外部DBなしで、`~/.config/sumika/config.yaml`に全情報を集約する。

## Goals / Non-Goals

Goals:
- プロジェクト情報のYAML永続化
- グローバル設定(エディタ・AIツール・ポート)の管理
- dotfilesでのGit管理を想定したパス設計

Non-Goals:
- 暗号化(APIキー等は別途環境変数で管理)
- 複数設定ファイルの分割(初期は1ファイル)

## Design

### スキーマ

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
      commands:         # 追加で実行するコマンド
        - "docker compose up -d"
    links:
      - label: "仕様書"
        url: "https://notion.so/..."
      - label: "staging"
        url: "https://staging.example.com"
```

### Goの構造体マッピング

```go
type Config struct {
    Version  int      `yaml:"version"`
    Settings Settings `yaml:"settings"`
    Projects []Project `yaml:"projects"`
}

type Settings struct {
    Port   int    `yaml:"port"`
    Editor string `yaml:"editor"`
    AITool string `yaml:"ai_tool"`
}

type Project struct {
    Name        string   `yaml:"name"`
    Path        string   `yaml:"path"`
    Description string   `yaml:"description"`
    Tags        []string `yaml:"tags"`
    Launch      Launch   `yaml:"launch"`
    Links       []Link   `yaml:"links"`
}

type Launch struct {
    Editor   bool     `yaml:"editor"`
    AI       bool     `yaml:"ai"`
    Commands []string `yaml:"commands"`
}

type Link struct {
    Label string `yaml:"label"`
    URL   string `yaml:"url"`
}
```

### デフォルト値

- `settings.port`: 8964
- `settings.editor`: 空文字(スキップ)
- `settings.ai_tool`: 空文字(スキップ)
- `project.launch.editor`: true
- `project.launch.ai`: true

## Open Questions

- `version`フィールドの移行ロジック(将来v2になった場合の対応)

## Related Files

internal/config/config.go
