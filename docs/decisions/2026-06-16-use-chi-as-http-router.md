# use chi as HTTP router

- Date: 2026-06-16
- Type: ADR
- Status: Active
- Author: masahiro.kasatani

## Context

GoのWebサーバーのルーターとして、標準の`net/http`のみではURLパラメータの扱いやミドルウェアの合成が煩雑になる。
軽量かつ`net/http`互換のルーターが必要。

## Decision

HTTPルーターとして`github.com/go-chi/chi`を採用する。
`net/http`の`http.Handler`インターフェースに完全準拠しており、将来的な置き換えが容易。

## Consequences

- **メリット**: `net/http`互換なので標準ライブラリとの混在が容易。URLパラメータ(`chi.URLParam`)やミドルウェアチェーンをシンプルに記述できる。依存が軽量。
- **トレードオフ**: `gin`や`echo`と比べてヘルパー関数は少ないが、sumikaの規模では問題ない。

## Alternatives Considered

- **net/http (標準のみ)**: 外部依存なしだが、URLパラメータやミドルウェア合成が煩雑。
- **gin**: 高機能だが`net/http`非互換な部分があり、ロックインリスクがある。
- **echo**: ginと同様に機能豊富だが、chiの軽量さで十分。

## Related Files

internal/server/server.go
go.mod
