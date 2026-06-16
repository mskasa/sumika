# Claude Code adapter implementation

- Date: 2026-06-16
- Type: Design
- Status: Draft
- Author: masahiro.kasatani

## Overview

Claude Codeアダプターの実装方針を定義する。
`AIAdapter`インターフェースを実装し、Claude Codeの起動・セッション履歴取得・CLAUDE.md読み込みを提供する。

## Background

sumikaのPhase 3でAIコンテキスト可視化を実現するために、Claude Codeの内部データ(`~/.claude/`配下)を読み込む必要がある。
内部仕様はバージョンアップで変更される可能性があるため、アダプター層に封じ込めてリスクを最小化する。

## Goals / Non-Goals

Goals:
- Claude Codeの起動(`os/exec`でサブプロセス)
- セッション履歴の取得(`~/.claude/projects/`配下のJSONL読み込み)
- CLAUDE.mdの読み込み
- バージョン変更に備えたバージョン検出ロジック

Non-Goals:
- Claude CodeのAPIを直接呼び出すこと(対象外)
- セッション履歴の書き込み

## Design

### 起動

```go
func (a *ClaudeCodeAdapter) Launch(projectPath string) error {
    cmd := exec.Command("claude")
    cmd.Dir = projectPath
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    return cmd.Start()
}
```

### セッション履歴の取得

- `~/.claude/projects/` 配下にプロジェクトパスをエンコードしたディレクトリが存在する
- ディレクトリ名の対応ルールはClaude Codeの内部仕様に従う
- 各ディレクトリ配下のJSONLファイルから最新セッションを読み込む
- 末尾N件のメッセージからサマリーを生成する

### CLAUDE.mdの管理

```go
func (a *ClaudeCodeAdapter) GetContextFile(projectPath string) (string, error) {
    p := filepath.Join(projectPath, "CLAUDE.md")
    data, err := os.ReadFile(p)
    if errors.Is(err, os.ErrNotExist) {
        return "", nil
    }
    return string(data), err
}
```

### バージョン検出

```go
func (a *ClaudeCodeAdapter) IsAvailable() bool {
    _, err := exec.LookPath("claude")
    return err == nil
}
```

## Implementation Plan

1. `internal/adapter/adapter.go`: `AIAdapter`インターフェースと`SessionSummary`型の定義
2. `internal/adapter/claudecode/claudecode.go`: Claude Codeアダプターの実装
3. セッションJSONLのパースロジックの実装
4. `internal/server/server.go`: アダプターをダッシュボードハンドラーに注入

## Open Questions

- `~/.claude/projects/`のディレクトリ名のエンコーディングルール(要調査)
- セッションJSONLのフォーマット変更を検知する方法

## Related Files

internal/adapter/adapter.go
internal/adapter/claudecode/claudecode.go
