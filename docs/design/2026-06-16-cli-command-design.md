# CLI command design

- Date: 2026-06-16
- Type: Design
- Status: Draft
- Author: masahiro.kasatani

## Overview

sumikaのCLIコマンド体系を定義する。
個人開発者が日常的に使う操作をシンプルなサブコマンドで提供し、cobraフレームワークで実装する。

## Background

個人開発者が複数プロジェクトを横断管理するために、プロジェクトの登録・一覧・起動・状態確認をCLIで行えることが必要。

## Goals / Non-Goals

Goals:
- プロジェクトのCRUD操作(add/list/remove)
- エディタ+AIツールの起動(open)
- 全プロジェクトのgit状態一覧(status)
- Webダッシュボードの起動(serve)
- 設定初期化(init)

Non-Goals:
- インタラクティブTUI(初期は不要)
- プロジェクト間のタスク管理

## Design

### コマンド一覧

```
sumika init                  # 設定ファイルを初期化
sumika add <path>            # プロジェクトを登録
sumika list                  # プロジェクト一覧を表示
sumika open <name>           # エディタ + AIツールを起動
sumika serve                 # Webダッシュボードを起動
sumika status                # 全プロジェクトのgit status一覧
sumika remove <name>         # プロジェクトを登録解除
```

### フラグ定義

**`sumika add <path>`**
- `--name string`: プロジェクト名(デフォルト: ディレクトリ名)
- `--description string`: 説明文

**`sumika serve`**
- `--port int`: ポート番号(デフォルト: config.yamlの設定値、未設定なら8964)

### 出力形式

`sumika list` の表形式出力例:
```
NAME         PATH                        DESCRIPTION          TAGS
my-api       ~/projects/my-api           REST API サーバー    backend, go
my-frontend  ~/projects/my-frontend      Next.js フロント     frontend, ts
```

`sumika status` の出力例:
```
NAME         LAST COMMIT              CHANGES
my-api       2026-06-16 10:30:00 +09  3 uncommitted
my-frontend  2026-06-15 22:00:00 +09  clean
```

## Implementation Plan

1. `cmd/sumika/main.go`: cobraルートコマンド定義
2. 各サブコマンドをcmd配下またはroot.goに追加
3. 各コマンドから`internal/`パッケージのロジックを呼び出す

## Open Questions

- `sumika open`の起動プロセスはデタッチ(バックグラウンド)か同期か
- `sumika serve`はフォアグラウンド動作のみか、デーモン化オプションを持つか

## Related Files

cmd/sumika/main.go
internal/config/config.go
internal/project/project.go
internal/launcher/launcher.go
