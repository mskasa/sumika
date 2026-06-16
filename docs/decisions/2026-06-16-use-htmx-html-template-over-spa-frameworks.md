# use HTMX + html/template over SPA frameworks

- Date: 2026-06-16
- Type: ADR
- Status: Active
- Author: masahiro.kasatani

## Context

sumikaのWebダッシュボードを実装するにあたり、フロントエンドの技術選定が必要。
「シングルバイナリ」という制約のもと、外部の静的ファイルサーバーや`node_modules`なしで動作させたい。
React/Vue等のSPAフレームワークはビルドステップとnpmエコシステムへの依存が発生する。

## Decision

フロントエンドにHTMXと`html/template`を採用する。
Goの`embed`パッケージでテンプレート・CSS・JSをバイナリに同梱し、npmレスのシングルバイナリを実現する。
HTMXのポーリング機能(hx-trigger="every 30s")でgit statusの自動更新を実装する。

## Consequences

- **メリット**: `npm install`不要でビルドが簡単。GoのHTMLテンプレートと自然に統合できる。バイナリに全リソースを`embed`できる。
- **トレードオフ**: 複雑なUI状態管理はSPAより難しい。コンポーネント単位の再利用性はReact等に劣る。

## Alternatives Considered

- **React/Next.js**: 開発体験は良いが、npmビルドステップが必須でシングルバイナリ配布と相性が悪い。
- **Vue.js + Vite**: 同様にビルドステップが必要。
- **Templ**: Goのコンポーネントテンプレートエンジン。将来の移行候補だが初期は標準の`html/template`で十分。

## Related Files

web/templates/
web/static/
internal/server/server.go
