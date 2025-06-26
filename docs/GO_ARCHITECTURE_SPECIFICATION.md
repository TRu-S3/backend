# Go アプリケーション アーキテクチャ仕様書

## 概要

TRu-S3 Backend は、Clean Architecture（クリーンアーキテクチャ）の原則に基づいて設計されたGoアプリケーションです。レイヤー分離により、保守性、テスト容易性、拡張性を重視した設計となっています。

## 最新の実装状況

**完全実装済みハンドラー**:
- UserHandler (ユーザー管理)
- TagHandler (タグ管理)
- ProfileHandler (プロフィール管理)
- MatchingHandler (マッチング機能)
- FileHandler (ファイル管理)
- ContestHandler (コンテスト管理)
- BookmarkHandler (ブックマーク機能)
- HackathonHandler (ハッカソン管理)

**全テーブル対応**: 完全なCRUD操作実装済み  
**GCP Cloud SQL**: PostgreSQL 17 本番環境対応  
**API エンドポイント**: 45+ のREST API提供

## アーキテクチャ概要

### 1. レイヤー構成

```
┌─────────────────────────────────────┐
│            Interfaces               │  ← HTTP Handlers, Routes
├─────────────────────────────────────┤
│           Application               │  ← Business Logic, Use Cases
├─────────────────────────────────────┤
│             Domain                  │  ← Entities, Repositories
├─────────────────────────────────────┤
│          Infrastructure             │  ← External Services (GCS, DB)
└─────────────────────────────────────┘
```

### 2. 依存関係のルール

- **内側のレイヤーは外側のレイヤーを知らない**
- **外側のレイヤーは内側のレイヤーに依存する**
- **依存関係の逆転原理**: インターフェースを使用して依存関係を抽象化

## プロジェクト構造

```
TRu-S3/
├── main.go                 # エントリーポイント
├── internal/               # 内部パッケージ
│   ├── interfaces/         # HTTPハンドラー、ルーティング
│   ├── application/        # ビジネスロジック、ユースケース
│   ├── domain/            # エンティティ、リポジトリインターフェース
│   ├── infrastructure/     # 外部サービス実装
│   ├── database/          # データベースモデル、接続
│   ├── config/            # 設定管理
│   └── utils/             # ユーティリティ関数
├── migrations/            # データベースマイグレーション
├── scripts/              # デプロイメントスクリプト
├── test/                 # テストファイル
└── docs/                # ドキュメント
```

## 各レイヤーの詳細

### 1. Interfaces Layer（インターフェース層）

#### 1.1 概要
HTTPリクエストの受付とレスポンスの返却を担当。外部からのアクセスポイント。

#### 1.2 主要コンポーネント

**ハンドラー（Handlers）**
- `UserHandler`: ユーザー管理API (完全CRUD)
- `TagHandler`: タグ管理API (完全CRUD)
- `ProfileHandler`: プロフィール管理API (完全CRUD + 特別エンドポイント)
- `MatchingHandler`: マッチング機能API (完全CRUD + ユーザーマッチ取得)
- `FileHandler`: ファイル管理API (完全CRUD + ダウンロード)
- `ContestHandler`: コンテスト管理API (完全CRUD)
- `BookmarkHandler`: ブックマーク管理API (完全CRUD)  
- `HackathonHandler`: ハッカソン管理API (完全CRUD + 参加者管理)
- `BaseHandler`: 共通ハンドラー機能

**ルーティング（Routes）**
- `routes.go`: APIエンドポイントの定義
- RESTful設計に基づくURL構造

#### 1.3 主要機能
- HTTPリクエストの受信・パース
- バリデーション（基本的な形式チェック）
- JSON レスポンスの生成
- エラーハンドリング
- ページネーション
- CORS設定

#### 1.4 設計パターン
```go
type FileHandler struct {
    fileService *application.FileService
}

func (h *FileHandler) CreateFile(c *gin.Context) {
    // 1. リクエストパース
    // 2. バリデーション
    // 3. アプリケーション層への委譲
    // 4. レスポンス生成
}
```

### 2. Application Layer（アプリケーション層）

#### 2.1 概要
ビジネスロジックとユースケースを実装。ドメインオブジェクトを使用してビジネスルールを実行。

#### 2.2 主要コンポーネント

**サービス（Services）**
- `FileService`: ファイル操作の業務ロジック

#### 2.3 主要機能
- ビジネスルールの実装
- トランザクション管理
- ドメインサービスの調整
- データ変換・検証
- エラーハンドリング

#### 2.4 設計パターン
```go
type FileService struct {
    fileRepo domain.FileRepository
}

func (s *FileService) CreateFile(ctx context.Context, req *domain.CreateFileRequest) (*domain.File, error) {
    // 1. ビジネスルール検証
    // 2. ドメインロジック実行
    // 3. リポジトリへの委譲
    // 4. 結果の返却
}
```

### 3. Domain Layer（ドメイン層）

#### 3.1 概要
ビジネスの核となるエンティティとルールを定義。外部の技術的詳細に依存しない純粋なビジネスロジック。

#### 3.2 主要コンポーネント

**エンティティ（Entities）**
- `File`: ファイルドメインエンティティ
- `FileData`: ファイルコンテンツ付きエンティティ

**リクエスト/レスポンス**
- `CreateFileRequest`: ファイル作成リクエスト
- `UpdateFileRequest`: ファイル更新リクエスト
- `FileQuery`: ファイル検索クエリ

**リポジトリインターフェース**
- `FileRepository`: ファイル永続化インターフェース

**エラー定義**
- `ErrFileNotFound`: ファイル未発見エラー
- `ErrFileAlreadyExists`: ファイル重複エラー
- `ErrInvalidFileName`: 無効ファイル名エラー

#### 3.3 設計原則
- **Pure Go**: 外部ライブラリに依存しない
- **Immutable**: エンティティは不変
- **Self-validating**: ドメインオブジェクトは自己検証

#### 3.4 エンティティ設計
```go
type File struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Path        string    `json:"path"`
    Size        int64     `json:"size"`
    ContentType string    `json:"content_type"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

### 4. Infrastructure Layer（インフラストラクチャー層）

#### 4.1 概要
外部システムとの接続を担当。データベース、クラウドサービス、外部APIとの通信を実装。

#### 4.2 主要コンポーネント

**Google Cloud Storage**
- `gcs_client.go`: GCSクライアントの作成・管理
- `gcs_file_repository.go`: ファイルストレージ実装

**データベース**
- `connection.go`: データベース接続管理
- `models.go`: GORM用モデル定義

#### 4.3 主要機能
- 外部サービス接続
- データ永続化
- 外部API通信
- 設定管理

#### 4.4 設計パターン
```go
type GCSFileRepository struct {
    client     *storage.Client
    bucketName string
    folder     string
}

func (r *GCSFileRepository) Create(ctx context.Context, req *domain.CreateFileRequest) (*domain.File, error) {
    // GCS への具体的な実装
}
```

## 技術スタック

### 1. Webフレームワーク
- **Gin**: `github.com/gin-gonic/gin v1.10.1`
  - 高性能HTTPルーター
  - ミドルウェアサポート
  - JSON バインディング

### 2. データベース
- **GORM**: `gorm.io/gorm v1.30.0`
  - ORM機能
  - マイグレーション
  - 接続プール管理
- **PostgreSQL Driver**: `gorm.io/driver/postgres v1.6.0`

### 3. クラウドサービス
- **Google Cloud Storage**: `cloud.google.com/go/storage v1.55.0`
  - ファイルストレージ
  - オブジェクト管理
- **Google Cloud APIs**: `google.golang.org/api v0.235.0`

### 4. 設定管理
- **godotenv**: `github.com/joho/godotenv v1.5.1`
  - 環境変数管理
  - 設定ファイル読み込み

### 5. テスティング
- **Testify**: `github.com/stretchr/testify v1.10.0`
  - アサーション
  - モック
  - テストスイート

## 設定管理

### 1. 設定構造体
```go
type Config struct {
    // Server Configuration
    Port    string
    GinMode string

    // GCP Configuration
    GCSBucketName                string
    GCSFolder                    string
    GoogleCloudProject           string
    GoogleApplicationCredentials string

    // Database Configuration
    DBHost                 string
    DBPort                 string
    DBName                 string
    DBUser                 string
    DBPassword             string
    DBSSLMode              string
    DBMaxOpenConns         int
    DBMaxIdleConns         int
    CloudSQLConnectionName string
    UseCloudSQLProxy       bool
}
```

### 2. 環境変数
| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| PORT | 8080 | サーバーポート |
| GIN_MODE | debug | Ginの実行モード |
| GCS_BUCKET_NAME | 202506-zenn-ai-agent-hackathon | GCSバケット名 |
| GCS_FOLDER | test | GCSフォルダ名 |
| DB_HOST | localhost | データベースホスト |
| DB_PORT | 5432 | データベースポート |
| DB_NAME | tru_s3 | データベース名 |
| DB_USER | postgres | データベースユーザー |
| DB_PASSWORD | postgres123 | データベースパスワード |

### 3. 設定検証
アプリケーション起動時に設定値の妥当性を検証：
- 必須項目の存在確認
- 数値型の妥当性チェック
- SSL設定の検証
- 接続プール設定の妥当性

## エラーハンドリング

### 1. エラー分類

**ドメインエラー**
```go
var (
    ErrFileNotFound      = errors.New("file not found")
    ErrFileAlreadyExists = errors.New("file already exists")
    ErrInvalidFileName   = errors.New("invalid file name")
    ErrInvalidFileSize   = errors.New("invalid file size")
)
```

**HTTPエラー**
- 400 Bad Request: 不正なリクエスト
- 404 Not Found: リソース未発見
- 409 Conflict: リソース競合
- 500 Internal Server Error: サーバーエラー

### 2. エラーレスポンス形式
```json
{
    "error": "Error message description"
}
```

## ユーティリティ

### 1. HTTP ユーティリティ（`utils/http.go`）

**ページネーション**
```go
type PaginationParams struct {
    Page   int
    Limit  int
    Offset int
}
```

**レスポンスヘルパー**
- `StandardResponse()`: 標準レスポンス
- `ErrorResponse()`: エラーレスポンス
- `CreatedResponse()`: 作成レスポンス

### 2. バリデーション ユーティリティ（`utils/validation.go`）

**バリデーション関数**
- `ValidateRequired()`: 必須フィールドチェック
- `ValidatePositiveInt()`: 正の整数チェック
- `ValidateEmail()`: メール形式チェック
- `ValidateDateRange()`: 日付範囲チェック

## セキュリティ

### 1. CORS設定
```go
r.Use(func(c *gin.Context) {
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
    
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(204)
        return
    }
    
    c.Next()
})
```

### 2. データベースセキュリティ
- **接続プール**: リソース枯渇防止
- **準備済みステートメント**: SQLインジェクション防止
- **SSL接続**: データ暗号化（本番環境）

### 3. ファイルセキュリティ
- **ファイル名検証**: パストラバーサル攻撃防止
- **Content-Type検出**: 適切なMIMEタイプ設定
- **サイズ制限**: リソース枯渇防止

## パフォーマンス最適化

### 1. データベース最適化
- **接続プール**: 最大25接続、アイドル5接続
- **プリロード**: N+1問題解決
- **インデックス**: 検索性能向上

### 2. HTTPパフォーマンス
- **コンテキストタイムアウト**: 長時間リクエスト防止
- **ページネーション**: 大量データ処理効率化
- **圧縮**: Ginの自動GZIP圧縮

### 3. ファイル処理最適化
- **ストリーミング**: 大容量ファイル対応
- **コンテキストキャンセル**: リクエスト中断対応

## テスト戦略

### 1. テスト種別
- **単体テスト**: 各レイヤーの個別テスト
- **統合テスト**: レイヤー間の連携テスト
- **APIテスト**: HTTPエンドポイントテスト

### 2. モック戦略
```go
// リポジトリモック例
type MockFileRepository struct {
    mock.Mock
}

func (m *MockFileRepository) Create(ctx context.Context, req *domain.CreateFileRequest) (*domain.File, error) {
    args := m.Called(ctx, req)
    return args.Get(0).(*domain.File), args.Error(1)
}
```

### 3. テストデータ管理
- **テストフィクスチャ**: 再現可能なテストデータ
- **データベースロールバック**: テスト間の独立性確保

## デプロイメント

### 1. Docker化
```dockerfile
FROM golang:1.24.2-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### 2. Google Cloud Run
- **自動スケーリング**: トラフィックに応じた自動拡張
- **HTTPS終端**: 自動SSL証明書
- **ヘルスチェック**: `/health` エンドポイント

### 3. CI/CD
- **Cloud Build**: 自動ビルド・デプロイ
- **環境変数管理**: Secret Manager連携

## 監視・ログ

### 1. ログ出力
```go
log.Printf("Configuration loaded:")
log.Printf("  Port: %s", cfg.Port)
log.Printf("  Gin Mode: %s", cfg.GinMode)
```

### 2. ヘルスチェック
- `GET /`: 基本ヘルスチェック
- `GET /health`: 詳細ヘルスチェック

### 3. エラー監視
- **構造化ログ**: JSON形式でのログ出力
- **エラー集約**: 集中エラー管理

## 今後の拡張予定

### 1. 認証・認可
- **JWT認証**: ユーザー認証機能
- **RBAC**: ロールベースアクセス制御
- **OAuth2**: 外部認証プロバイダー連携

### 2. 非同期処理
- **メッセージキュー**: Cloud Tasks, Pub/Sub
- **バックグラウンドジョブ**: ファイル処理、通知

### 3. キャッシュ
- **Redis**: セッション管理、キャッシュ
- **CDN**: 静的ファイル配信

### 4. 監視強化
- **メトリクス収集**: Prometheus連携
- **トレーシング**: OpenTelemetry導入
- **アラート**: 異常検知・通知

## ベストプラクティス

### 1. コード品質
- **Linting**: golangci-lint使用
- **フォーマッティング**: gofmt統一
- **テストカバレッジ**: 80%以上維持

### 2. Git管理
- **ブランチ戦略**: Git Flow
- **コミットメッセージ**: Conventional Commits
- **PR レビュー**: 必須レビュープロセス

### 3. ドキュメント
- **API仕様**: OpenAPI 3.0準拠
- **アーキテクチャ図**: 定期更新
- **運用手順書**: 障害対応マニュアル

このアーキテクチャにより、保守性が高く、テスタブルで、拡張可能なGoアプリケーションが実現されています。