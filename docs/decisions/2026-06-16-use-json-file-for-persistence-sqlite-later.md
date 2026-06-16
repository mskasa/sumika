# use JSON file for persistence, SQLite later

- Date: 2026-06-16
- Type: ADR
- Status: Active
- Author: masahiro.kasatani

## Context

プロジェクト情報の永続化方式を決定する必要がある。
初期のMVP段階では、プロジェクト数は数十件程度を想定。
設定ファイル(`~/.config/sumika/config.yaml`)とは別に、動的なデータ(タグ・メモ等)を保存する場所が必要になる可能性がある。

## Decision

初期実装ではYAMLファイル(`~/.config/sumika/config.yaml`)に設定とプロジェクト情報を一元管理する。
プロジェクト数が増えてパフォーマンスが問題になった場合、SQLite(CGO不要な`modernc.org/sqlite`)に移行する。

## Consequences

- **メリット**: 外部DBプロセス不要でシンプル。設定ファイルがGitで管理可能(dotfiles)。人間が直接編集できる。
- **トレードオフ**: 同時書き込み時の競合リスク。プロジェクト数が多い場合の全件読み込みコスト。

## Alternatives Considered

- **PostgreSQL/MySQL**: Docker前提になりエンタープライズ寄りすぎる。個人開発者向けツールとして不適切。
- **SQLite (初期から)**: CGO依存でクロスコンパイルが複雑になる可能性がある。初期はYAMLで十分。
- **BoltDB**: 組み込みKVSだが、人間が読めないバイナリ形式のためdotfiles管理に不向き。

## Related Files

internal/config/config.go
