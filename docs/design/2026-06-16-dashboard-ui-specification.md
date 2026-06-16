# dashboard UI specification

- Date: 2026-06-16
- Type: Design
- Status: Draft
- Author: masahiro.kasatani

## Overview

sumikaのWebダッシュボードは、登録済みプロジェクトをカードグリッドで一覧表示し、git状態・最終更新・AIセッション要約を可視化する。
HTMXによる30秒ポーリングでリアルタイム更新を行い、[Open]ボタンからエディタ+AIツールを起動できる。

## Background

個人開発者が複数プロジェクトを並行して管理する際に、「どこまでやったか」を素早く把握できるUIが必要。
CLIの`sumika status`を補完するグラフィカルなビューとして機能する。

## Goals / Non-Goals

Goals:
- プロジェクトカードのグリッド表示
- git status・最終コミット日時の表示
- AIセッション要約の表示(Phase 3)
- [Open]ボタンでエディタ+AIツールを起動
- タグによるフィルタリング・ソート

Non-Goals:
- リアルタイムWebSocket通信(HTMXポーリングで代替)
- プロジェクトの作成・削除(CLIで行う)
- モバイル向け最適化(初期は不要)

## Design

### プロジェクトカードレイアウト

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

### ダッシュボード全体構成

- ヘッダー: sumikaロゴ・グローバルフィルター(タグ選択)・ソート選択
- メインコンテンツ: プロジェクトカードのグリッド(2〜3カラム)
- フッター: バージョン情報

### HTMXポーリング

```html
<div hx-get="/api/projects" hx-trigger="every 30s" hx-swap="innerHTML">
```

### APIエンドポイント

| メソッド | パス | 説明 |
|---|---|---|
| GET | `/` | ダッシュボードHTML |
| GET | `/api/projects` | プロジェクトカードHTML(HTMX用) |
| POST | `/api/projects/{name}/open` | エディタ+AIツール起動 |

## Implementation Plan

1. `internal/server/server.go`: chiルーター・ハンドラー実装
2. `web/templates/index.html`: ダッシュボードベーステンプレート
3. `web/templates/project_card.html`: カードコンポーネント
4. `web/static/style.css`: グリッドレイアウト・カードスタイル

## Open Questions

- カードの起動状態(エディタが開いているか)をどう検知するか
- AIセッション要約のポーリング間隔(30秒は重いか)

## Related Files

internal/server/server.go
web/templates/
web/static/
