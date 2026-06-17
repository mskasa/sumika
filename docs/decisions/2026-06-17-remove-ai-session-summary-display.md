# Remove AI Session Summary Display from Dashboard

- Date: 2026-06-17
- Type: ADR
- Status: Accepted
- Author: mskasa

## Context

Phase 3 でダッシュボードの各プロジェクトカードに「前回AIセッション」として、
Claude Code の JSONL セッションファイルから取得した最後のアシスタントメッセージを表示する機能を実装した。

実際に利用してみたところ、アシスタントの回答の断片（200文字）だけでは
「前回どんな作業をしていたか」を思い出すには不十分であることが分かった。

## Decision

「前回AIセッション」の表示機能をダッシュボードから削除する。

セッション開始時に Claude Code 自身に「前回の作業内容を教えて」と問えば、
Claude Code はセッション履歴を参照して作業内容を教えてくれる。
sumika 側で要約を表示する必要はない。

## Consequences

- ダッシュボードのカードがシンプルになる
- `SessionSummary` 構造体・`GetSessionSummary` メソッド・JSONL パース処理を削除できる
- `adapter.AIAdapter` インターフェースから `GetSessionSummary` を除去
- `server.New` がアダプターを引数に取らなくなる（現時点でサーバーはアダプターを使わない）

## Alternatives Considered

- **表示文字数を増やす**: 根本的な解決にならない
- **ユーザーメッセージも一緒に表示**: 情報量は増えるが作業想起には依然不十分
- **Claude API でサマリー生成**: コストと遅延が発生する。ユーザー自身が claude に聞けば済む

## Related Files

