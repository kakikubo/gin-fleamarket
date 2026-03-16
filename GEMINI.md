# Gin Fleamarket - Project Mandates (GEMINI.md)

このファイルは、Antigravity (Gemini CLI) がこのリポジトリで作業する際に遵守すべき基盤となる指針です。グローバルな設定よりも優先されます。

## 1. プロジェクト概要
Go (Gin framework) を使用したフリマアプリのバックエンドプロジェクトです。
クリーンアーキテクチャ（Controllers -> Services -> Repositories -> Models）を採用し、明確な関心の分離と依存性注入（DI）を行っています。

## 2. 技術スタック & ツール
- **Backend:** Go (Gin), GORM
- **Database:** PostgreSQL (Prod/Dev), SQLite (Test)
- **Auth:** JWT (Middleware 経由)
- **Validation:** `gotest.tools/v3/assert`
- **Infrastructure:** Docker Compose (PostgreSQL, pgAdmin)
- **Hot Reload:** Air (`.air.toml`)

## 3. 開発およびコーディング規約
### アーキテクチャ層
- **controllers/**: HTTP リクエスト/レスポンスのハンドリング。
- **services/**: ビジネスロジックの記述。
- **repositories/**: データベース操作の抽象化。
- **models/**: GORM モデル定義。
- **dto/**: リクエスト/レスポンス用の構造体。
- **infra/**: DB 接続、初期化などの基盤コード。

### 実装のルール
- **依存性注入 (DI)**: コンストラクタ関数（例: `NewItemService`）を使用して手動で DI を行います。
- **認証**: 保護されたルートには必ず `middlewares.AuthMiddleware` を適用します。
- **エラーハンドリング**: Gin のコンテキスト（`c.JSON`）を使用して一貫したエラーレスポンスを返します。

## 4. テスト戦略
- **環境**: テスト時は `.env.test` を読み込み、SQLite (In-memory) を使用します。
- **実行**: `go test ./...` で全テストを実行します。
- **作成**: テスト対象と同じディレクトリに `*_test.go` を作成します。
- **検証**: `gotest.tools/v3/assert` を使用してアサーションを行います。
- **モック**: 可能な限り `setupTest` 関数（`controllers/auth_controller_test.go` 等を参照）を作成し、実際の下位レイヤーを組み合わせて結合テストに近い形で行います。

## 5. 安全ルール（重要）
- **環境変数**: `.env` ファイルを直接編集する際は、必ずバックアップを作成し、秘密情報の露出に注意します。
- **破壊的変更**: 既存の `models` や `migrations` を変更する場合は、影響範囲を慎重に調査し、事前にユーザーへ報告します。
- **削除操作**: ファイルの削除（`rm`）やデータベースのドロップは原則禁止です。必要な場合は必ず承認を得ます。
- **依存関係の追加**: `go get` 等でパッケージを追加する場合は、目的と影響を説明し、承認を得てから実行します。

## 6. コミュニケーション
- **言語**: 全ての応答、プランニング、コードコメント（必要な場合）は **日本語** で行います。
- **透明性**: 実行するコマンドとその意図を事前に説明します。
