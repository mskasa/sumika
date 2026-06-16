# use Go over other languages

- Date: 2026-06-16
- Type: ADR
- Status: Active
- Author: masahiro.kasatani

## Context

sumikaは個人開発者向けの軽量プロジェクトハブCLIツール。
配布の容易さ・クロスプラットフォーム対応・依存関係の少なさが重要な要件。
ユーザーは `brew install` や単一バイナリのダウンロードで即座に使えることを期待する。

## Decision

実装言語としてGoを採用する。
単一バイナリにコンパイルできるため配布が簡単で、クロスプラットフォームのクロスコンパイルが標準サポートされている。
標準ライブラリが充実しており、CLIツール・Webサーバー・ファイル操作・プロセス管理に必要な機能が揃っている。

## Consequences

- **メリット**: `go build` 一発で単一バイナリを生成できる。GoReleaser で macOS/Linux/Windows 向けのクロスコンパイルとリリースを自動化できる。`embed` パッケージでフロントエンドアセットをバイナリに同梱できる。
- **トレードオフ**: 動的型付け言語と比べると記述量が増えることがある。

## Alternatives Considered

- **Node.js (TypeScript)**: npmの依存関係管理が複雑になりやすく、配布に`pkg`等の追加ツールが必要。
- **Python**: 実行環境のバージョン依存があり、配布に手間がかかる。
- **Rust**: 単一バイナリ配布は可能だが、CLIツール開発のエコシステムの成熟度とビルド時間を考慮してGoを選択。

## Related Files

<!-- List files related to this decision (e.g. internal/search/search.go). -->
go.mod
cmd/sumika/main.go
