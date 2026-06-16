# adopt adapter pattern for AI tools

- Date: 2026-06-16
- Type: ADR
- Status: Active
- Author: masahiro.kasatani

## Context

sumikaのコア機能はAIツールに依存すべきではない。
Claude Code以外のユーザー(cursor, aider等を使うユーザー)にも価値を提供するために、AIツール依存のコードを分離する必要がある。
将来的に複数のAIツールをサポートする拡張性も必要。

## Decision

AIツールとのインタラクションをアダプターパターンで実装する。
`internal/adapter/adapter.go`にインターフェースを定義し、各ツールの実装を`internal/adapter/<toolname>/`に配置する。
Claude Codeを最初の公式アダプターとして実装する。

```go
type AIAdapter interface {
    Name() string
    IsAvailable() bool
    Launch(projectPath string) error
    GetSessionSummary(projectPath string) (*SessionSummary, error)
    GetContextFile(projectPath string) (string, error)
}
```

## Consequences

- **メリット**: コア機能(Layer 1)とAIツール連携(Layer 2)が分離され、Claude Code以外のユーザーでも価値が成立する。新しいAIツールの追加が容易。
- **トレードオフ**: 初期実装の抽象化コストがある。インターフェースの設計が各ツールの特性に合わない場合に変更が必要になることがある。

## Alternatives Considered

- **Claude Code直接呼び出し**: 実装が単純だが、他ツールへの対応が困難になる。
- **プラグインシステム**: 実行時ロードは複雑すぎる。コンパイル時のアダプターで十分。

## Related Files

internal/adapter/adapter.go
internal/adapter/claudecode/claudecode.go
