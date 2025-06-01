# Gin Fleamarket プロジェクトガイドラインおよびJunie利用のガイドライン

このドキュメントは、Gin Fleamarketプロジェクトに取り組む開発者のための重要な情報を提供します。
## ルール

- 回答は全て日本語で行うこと
- プランニングの過程も日本語で出力すること

## ビルド/設定手順

### 環境設定

1. **環境変数**:
   - `.env.sample`を`.env`にコピーして、変数を設定します:
     ```
     ENV=prod
     DB_HOST=localhost
     DB_USER=ginuser
     DB_PASSWORD=ginpassword
     DB_NAME=fleamarket
     DB_PORT=15432
     SECRET_KEY=`openssl rand -hex 32で生成`
     ```

2. **データベース設定**:
   - プロジェクトはDocker Composeを使用してPostgreSQLとpgAdminをセットアップします:
     ```bash
     docker-compose up -d
     ```
   - PostgreSQLはポート15432で利用可能になります
   - pgAdminはhttp://localhost:81で利用可能になります（ログイン: gin@example.com、パスワード: ginpassword）

3. **アプリケーションの実行**:
   - 以下のコマンドでアプリケーションを起動します:
     ```bash
     go run main.go
     ```
   - サーバーはhttp://localhost:8080で実行されます

### プロジェクト構造

- **controllers/**: HTTPリクエストハンドラー
- **dto/**: リクエスト/レスポンス用のデータ転送オブジェクト
- **infra/**: インフラストラクチャコード（DB設定、初期化処理）
- **middlewares/**: HTTPミドルウェア（認証など）
- **migrations/**: データベースマイグレーションスクリプト
- **models/**: データモデル
- **repositories/**: データアクセス層
- **services/**: ビジネスロジック層

## テスト情報

### テスト設定

- テストはPostgreSQLの代わりにインメモリSQLiteデータベースを使用します
- テスト環境は`.env.test`で設定されています（`ENV=test`のみ設定）
- `ENV=test`の場合、アプリケーションは自動的にSQLiteを使用します

### テストの実行

- すべてのテストを実行:
  ```bash
  go test ./...
  ```

- 特定のパッケージのテストを実行:
  ```bash
  go test ./controllers
  ```

- 詳細出力で特定のテストを実行:
  ```bash
  go test -v ./controllers -run TestSignup
  ```

### テストの作成

1. **テスト構造**:
   - テストはテスト対象のコードに隣接する`*_test.go`という名前のファイルに配置する
   - 標準のGoテストパッケージとhttptestを使用してHTTPテストを行う
   - アサーションにはgotest.tools/v3/assertを使用する

2. **テスト例**:
   ```
   func TestSignup(t *testing.T) {
       // テスト環境のセットアップ
       r, _ := setupAuthTest()

       // テストリクエストの作成
       signupInput := dto.SignupInput{
           Email:    "test@example.com",
           Password: "password123",
       }
       reqBody, _ := json.Marshal(signupInput)

       // リクエストの実行
       w := httptest.NewRecorder()
       req, _ := http.NewRequest("POST", "/auth/signup", bytes.NewBuffer(reqBody))
       req.Header.Set("Content-Type", "application/json")
       r.ServeHTTP(w, req)

       // レスポンスの検証
       assert.Equal(t, http.StatusCreated, w.Code)
   }
   ```

3. **テストセットアップ**:
   - テスト環境を初期化するセットアップ関数を作成します:
     ```
     func setupAuthTest() (*gin.Engine, *gorm.DB) {
         // テスト環境の読み込み
         if err := godotenv.Load("../.env.test"); err != nil {
             os.Setenv("ENV", "test") // フォールバック
         }

         // データベースのセットアップ
         db := infra.SetupDB()
         db.AutoMigrate(&models.User{})
         db.Exec("DELETE FROM users") // 既存データのクリア

         // ルーターのセットアップ
         r := gin.Default()
         authRepository := repositories.NewAuthRepository(db)
         authService := services.NewAuthService(authRepository)
         authController := NewAuthController(authService)

         // ルートのセットアップ
         authGroup := r.Group("/auth")
         authGroup.POST("/signup", authController.Signup)
         authGroup.POST("/login", authController.Login)

         return r, db
     }
     ```

## 追加の開発情報

### コードスタイルと規約

1. **プロジェクトアーキテクチャ**:
   - プロジェクトは明確な関心の分離を持つクリーンアーキテクチャパターンに従っています:
     - Controllersは HTTPリクエストとレスポンスを処理
     - Servicesはビジネスロジックを含む
     - Repositoriesはデータアクセスを処理
     - Modelsはデータ構造を定義
     - DTOsはリクエスト/レスポンス構造を定義

2. **依存性注入**:
   - プロジェクトは手動の依存性注入を使用しています:
     ```
     // リポジトリ層
     itemRepository := repositories.NewItemRepository(db)

     // サービス層（リポジトリを使用）
     itemService := services.NewItemService(itemRepository)

     // コントローラ層（サービスを使用）
     itemController := controllers.NewItemController(itemService)
     ```

3. **認証**:
   - JWT認証はauthサービスとミドルウェアに実装されています
   - 保護されたルートはAuthMiddlewareを使用する必要があります:
     ```
     itemRouterWithAuth := r.Group("/items", middlewares.AuthMiddleware(authService))
     ```

### デバッグ

1. **環境モード**:
   - `ENV=prod`: PostgreSQLデータベースを使用
   - `ENV=test`: インメモリSQLiteデータベースを使用
   - Ginはデフォルトでデバッグモードで実行され、詳細なログを提供します

2. **データベースアクセス**:
   - pgAdmin（http://localhost:81）を使用してPostgreSQLデータベースを検査・クエリ
   - SQLite（テストモード）の場合、データはメモリ内に保存され、永続化されません

3. **APIテスト**:
   - サンプルのInsomniaコレクションが`sample/Insomnia.json`で利用可能
   - これをInsomniaにインポートしてAPIエンドポイントをテスト
