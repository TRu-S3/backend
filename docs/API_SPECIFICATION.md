# API仕様書

## 概要

TRu-S3 Backend APIは、ユーザー管理、マッチング機能、ファイル管理、コンテスト管理、ブックマーク機能、ハッカソン管理を提供する包括的なRESTful APIです。完全なCRUD操作をサポートし、Clean Architectureに基づいて設計されています。

## 最新の実装状況

**完全実装済み**:
- ユーザー管理 (Users)
- タグ管理 (Tags)  
- プロフィール管理 (Profiles)
- マッチング機能 (Matchings)
- ファイル管理 (Files)
- コンテスト管理 (Contests)
- ブックマーク機能 (Bookmarks)
- ハッカソン管理 (Hackathons)
- ハッカソン参加者管理 (Participants)

## 基本情報

- **ベースURL**: `/api/v1`
- **認証**: 現在未実装（今後JWT認証を予定）
- **レスポンス形式**: JSON
- **HTTPメソッド**: GET, POST, PUT, DELETE, OPTIONS
- **CORS**: 全オリジンからのアクセスを許可

## 共通仕様

### レスポンス形式

#### 成功レスポンス
```json
{
  "status": "success",
  "data": { ... }
}
```

#### エラーレスポンス
```json
{
  "error": "エラーメッセージ"
}
```

### ページネーション

リスト取得APIでは以下のクエリパラメータでページネーションが可能です：

- `page`: ページ番号（デフォルト: 1）
- `limit`: 取得件数（デフォルト: 10、最大: 100）

レスポンスには以下の形式でページネーション情報が含まれます：

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

### 日付形式

全ての日付はRFC3339形式（ISO 8601）を使用します：
`2023-12-31T23:59:59Z`

## ヘルスチェック

### 基本ヘルスチェック
- **エンドポイント**: `GET /`
- **レスポンス**:
```json
{
  "message": "TRu-S3 Backend is running!",
  "status": "healthy"
}
```

### ヘルスチェック
- **エンドポイント**: `GET /health`
- **レスポンス**:
```json
{
  "status": "ok"
}
```

## 1. ユーザー管理API

### 1.1 ユーザー作成
- **エンドポイント**: `POST /api/v1/users`
- **リクエストボディ**:
```json
{
  "name": "田中太郎",
  "gmail": "tanaka@example.com"
}
```
- **レスポンス**:
```json
{
  "id": 1,
  "name": "田中太郎",
  "gmail": "tanaka@example.com",
  "created_at": "2023-12-31T23:59:59Z",
  "updated_at": "2023-12-31T23:59:59Z"
}
```
- **エラー**:
  - `400`: 無効なメールアドレス
  - `409`: メールアドレスが既に存在

### 1.2 ユーザー一覧取得
- **エンドポイント**: `GET /api/v1/users`
- **クエリパラメータ**:
  - `page`: ページ番号（デフォルト: 1）
  - `limit`: 取得件数（デフォルト: 10、最大: 100）
  - `name`: 名前での部分検索
  - `gmail`: Gmailでの部分検索
- **レスポンス**:
```json
{
  "users": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25
  }
}
```

### 1.3 ユーザー詳細取得
- **エンドポイント**: `GET /api/v1/users/{id}`
- **レスポンス**: ユーザー詳細（プロフィール、ブックマーク含む）

### 1.4 ユーザー更新
- **エンドポイント**: `PUT /api/v1/users/{id}`
- **リクエストボディ**: 更新したいフィールドのみ送信

### 1.5 ユーザー削除
- **エンドポイント**: `DELETE /api/v1/users/{id}`

### 1.6 ユーザーのマッチ一覧
- **エンドポイント**: `GET /api/v1/users/{user_id}/matches`
- **クエリパラメータ**:
  - `status`: マッチングステータス（デフォルト: accepted）

## 2. タグ管理API

### 2.1 タグ作成
- **エンドポイント**: `POST /api/v1/tags`
- **リクエストボディ**:
```json
{
  "name": "フロントエンド"
}
```

### 2.2 タグ一覧取得
- **エンドポイント**: `GET /api/v1/tags`
- **クエリパラメータ**:
  - `name`: タグ名での部分検索

### 2.3 タグ詳細取得
- **エンドポイント**: `GET /api/v1/tags/{id}`

### 2.4 タグ更新
- **エンドポイント**: `PUT /api/v1/tags/{id}`

### 2.5 タグ削除
- **エンドポイント**: `DELETE /api/v1/tags/{id}`
- **注意**: プロフィールで使用中のタグは削除できません

## 3. プロフィール管理API

### 3.1 プロフィール作成
- **エンドポイント**: `POST /api/v1/profiles`
- **リクエストボディ**:
```json
{
  "user_id": 1,
  "tag_id": 2,
  "bio": "Webエンジニアです",
  "age": 25,
  "location": "東京"
}
```

### 3.2 プロフィール一覧取得
- **エンドポイント**: `GET /api/v1/profiles`
- **クエリパラメータ**:
  - `user_id`: ユーザーIDでフィルタ
  - `tag_id`: タグIDでフィルタ
  - `location`: 場所での部分検索
  - `min_age`, `max_age`: 年齢範囲フィルタ

### 3.3 プロフィール詳細取得
- **エンドポイント**: `GET /api/v1/profiles/{id}`

### 3.4 ユーザーIDでプロフィール取得
- **エンドポイント**: `GET /api/v1/profiles/user/{user_id}`

### 3.5 プロフィール更新
- **エンドポイント**: `PUT /api/v1/profiles/{id}`

### 3.6 プロフィール削除
- **エンドポイント**: `DELETE /api/v1/profiles/{id}`

## 4. マッチング管理API

### 4.1 マッチング作成
- **エンドポイント**: `POST /api/v1/matchings`
- **リクエストボディ**:
```json
{
  "user1_id": 1,
  "user2_id": 2,
  "status": "pending"
}
```
- **ステータス**: `pending`, `accepted`, `rejected`, `blocked`

### 4.2 マッチング一覧取得
- **エンドポイント**: `GET /api/v1/matchings`
- **クエリパラメータ**:
  - `user_id`: 指定ユーザーが関わるマッチング
  - `user1_id`, `user2_id`: 特定ユーザーでフィルタ
  - `status`: ステータスでフィルタ

### 4.3 マッチング詳細取得
- **エンドポイント**: `GET /api/v1/matchings/{id}`

### 4.4 マッチング更新
- **エンドポイント**: `PUT /api/v1/matchings/{id}`
- **主な用途**: ステータス変更（承認・拒否など）

### 4.5 マッチング削除
- **エンドポイント**: `DELETE /api/v1/matchings/{id}`

## 5. ファイル管理API

### 5.1 ファイルアップロード
- **エンドポイント**: `POST /api/v1/files`
- **Content-Type**: `multipart/form-data`
- **リクエストパラメータ**:
  - `file`: ファイル（必須）
- **レスポンス**:
```json
{
  "id": "file-uuid",
  "name": "example.txt",
  "size": 1024,
  "content_type": "text/plain",
  "path": "files/example.txt",
  "created_at": "2023-12-31T23:59:59Z",
  "updated_at": "2023-12-31T23:59:59Z"
}
```
- **エラー**:
  - `400`: ファイルが提供されていない、無効なファイル名
  - `409`: ファイルが既に存在する
  - `500`: サーバーエラー

### 5.2 ファイル一覧取得
- **エンドポイント**: `GET /api/v1/files`
- **クエリパラメータ**:
  - `prefix`: ファイル名のプレフィックス（任意）
  - `page`: ページ番号（デフォルト: 1）
  - `limit`: 取得件数（デフォルト: 100、最大: 100）
  - `offset`: オフセット（デフォルト: 0）
- **レスポンス**:
```json
{
  "files": [
    {
      "id": "file-uuid",
      "name": "example.txt",
      "size": 1024,
      "content_type": "text/plain",
      "path": "files/example.txt",
      "created_at": "2023-12-31T23:59:59Z",
      "updated_at": "2023-12-31T23:59:59Z"
    }
  ],
  "count": 1
}
```

### 5.3 ファイル詳細取得
- **エンドポイント**: `GET /api/v1/files/{id}`
- **パスパラメータ**:
  - `id`: ファイルID（必須）
- **レスポンス**: ファイルメタデータ
- **エラー**:
  - `400`: ファイルIDが無効
  - `404`: ファイルが見つからない

### 5.4 ファイルダウンロード
- **エンドポイント**: `GET /api/v1/files/{id}/download`
- **パスパラメータ**:
  - `id`: ファイルID（必須）
- **レスポンス**: ファイルの内容（バイナリ）
- **ヘッダー**:
  - `Content-Disposition`: attachment; filename=ファイル名
  - `Content-Type`: ファイルのMIMEタイプ
  - `Content-Length`: ファイルサイズ
- **エラー**:
  - `400`: ファイルIDが無効
  - `404`: ファイルが見つからない

### 5.5 ファイル更新
- **エンドポイント**: `PUT /api/v1/files/{id}`
- **Content-Type**: `multipart/form-data`
- **パスパラメータ**:
  - `id`: ファイルID（必須）
- **リクエストパラメータ**:
  - `name`: 新しいファイル名（任意）
  - `file`: 新しいファイル（任意）
- **レスポンス**: 更新されたファイルメタデータ
- **エラー**:
  - `400`: 無効なパラメータ
  - `404`: ファイルが見つからない

### 5.6 ファイル削除
- **エンドポイント**: `DELETE /api/v1/files/{id}`
- **パスパラメータ**:
  - `id`: ファイルID（必須）
- **レスポンス**:
```json
{
  "message": "File deleted successfully"
}
```
- **エラー**:
  - `400`: ファイルIDが無効
  - `404`: ファイルが見つからない

## 6. コンテスト管理API

### 6.1 コンテスト作成
- **エンドポイント**: `POST /api/v1/contests`
- **リクエストボディ**:
```json
{
  "backend_quota": 5,
  "frontend_quota": 3,
  "ai_quota": 2,
  "application_deadline": "2023-12-31T23:59:59Z",
  "purpose": "ウェブアプリケーション開発コンテスト",
  "message": "革新的なウェブアプリケーションを作成してください",
  "author_id": 1
}
```
- **レスポンス**: 作成されたコンテスト情報
- **エラー**:
  - `400`: 無効な入力データ
  - `500`: サーバーエラー

### 6.2 コンテスト一覧取得
- **エンドポイント**: `GET /api/v1/contests`
- **クエリパラメータ**:
  - `page`: ページ番号（デフォルト: 1）
  - `limit`: 取得件数（デフォルト: 10、最大: 100）
  - `author_id`: 作成者IDでフィルタ
  - `active`: `true`の場合、締切前のコンテストのみ取得
- **レスポンス**:
```json
{
  "contests": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25
  }
}
```

### 6.3 コンテスト詳細取得
- **エンドポイント**: `GET /api/v1/contests/{id}`
- **パスパラメータ**:
  - `id`: コンテストID
- **レスポンス**: コンテスト詳細情報
- **エラー**:
  - `404`: コンテストが見つからない

### 6.4 コンテスト更新
- **エンドポイント**: `PUT /api/v1/contests/{id}`
- **リクエストボディ**: 更新したいフィールドのみ送信
```json
{
  "backend_quota": 10,
  "purpose": "新しいコンテストの目的"
}
```
- **レスポンス**: 更新されたコンテスト情報
- **エラー**:
  - `400`: 無効な入力データ
  - `404`: コンテストが見つからない

### 6.5 コンテスト削除
- **エンドポイント**: `DELETE /api/v1/contests/{id}`
- **レスポンス**:
```json
{
  "message": "Contest deleted successfully"
}
```
- **エラー**:
  - `404`: コンテストが見つからない

## 7. ブックマーク管理API

### 7.1 ブックマーク作成
- **エンドポイント**: `POST /api/v1/bookmarks`
- **リクエストボディ**:
```json
{
  "user_id": 1,
  "bookmarked_user_id": 2
}
```
- **レスポンス**: 作成されたブックマーク情報
- **エラー**:
  - `400`: 自分自身をブックマークしようとした場合
  - `500`: サーバーエラー

### 7.2 ブックマーク一覧取得
- **エンドポイント**: `GET /api/v1/bookmarks`
- **クエリパラメータ**:
  - `page`: ページ番号（デフォルト: 1）
  - `limit`: 取得件数（デフォルト: 10、最大: 100）
  - `user_id`: ユーザーIDでフィルタ
- **レスポンス**:
```json
{
  "bookmarks": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 15
  }
}
```

### 7.3 ブックマーク更新
- **エンドポイント**: `PUT /api/v1/bookmarks/{id}`
- **リクエストボディ**: 更新したいフィールドのみ送信
- **注意**: 自分自身のブックマークは作成できません

### 7.4 ブックマーク削除
- **エンドポイント**: `DELETE /api/v1/bookmarks/{id}`
- **レスポンス**:
```json
{
  "message": "Bookmark deleted successfully"
}
```
- **エラー**:
  - `404`: ブックマークが見つからない

## 8. ハッカソン管理API

### 8.1 ハッカソン作成
- **エンドポイント**: `POST /api/v1/hackathons`
- **リクエストボディ**:
```json
{
  "name": "Tech Innovation Hackathon 2023",
  "description": "革新的な技術ソリューションを開発するハッカソン",
  "start_date": "2023-06-15T09:00:00Z",
  "end_date": "2023-06-17T18:00:00Z",
  "registration_start": "2023-05-01T00:00:00Z",
  "registration_deadline": "2023-06-10T23:59:59Z",
  "max_participants": 100,
  "location": "東京国際フォーラム",
  "organizer": "Tech Community Japan",
  "contact_email": "info@techcommunity.jp",
  "prize_info": "1位: 100万円, 2位: 50万円, 3位: 25万円",
  "rules": "チーム人数は最大4名まで",
  "tech_stack": "React, Node.js, Python, AWS",
  "is_public": true,
  "banner_url": "https://example.com/banner.jpg",
  "website_url": "https://hackathon.example.com"
}
```
- **レスポンス**: 作成されたハッカソン情報（status: "upcoming"で作成）
- **エラー**:
  - `400`: 無効な日付形式、日付の論理エラー

### 8.2 ハッカソン一覧取得
- **エンドポイント**: `GET /api/v1/hackathons`
- **クエリパラメータ**:
  - `page`: ページ番号（デフォルト: 1）
  - `limit`: 取得件数（デフォルト: 10、最大: 100）
  - `status`: ステータスでフィルタ（upcoming, ongoing, completed, cancelled）
  - `organizer`: 主催者名でフィルタ（部分一致）
  - `is_public`: `true`/`false`で公開状態フィルタ
  - `upcoming`: `true`の場合、今後開催のハッカソンのみ
  - `ongoing`: `true`の場合、開催中のハッカソンのみ
  - `registration_open`: `true`の場合、参加登録受付中のハッカソンのみ
- **レスポンス**:
```json
{
  "hackathons": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 30
  }
}
```

### 8.3 ハッカソン詳細取得
- **エンドポイント**: `GET /api/v1/hackathons/{id}`
- **レスポンス**: ハッカソン詳細情報（参加者情報を含む）
- **エラー**:
  - `404`: ハッカソンが見つからない

### 8.4 ハッカソン更新
- **エンドポイント**: `PUT /api/v1/hackathons/{id}`
- **リクエストボディ**: 更新したいフィールドのみ送信
```json
{
  "name": "Updated Hackathon Name",
  "status": "ongoing",
  "max_participants": 150
}
```
- **レスポンス**: 更新されたハッカソン情報
- **エラー**:
  - `400`: 無効なステータス、日付の論理エラー
  - `404`: ハッカソンが見つからない

### 8.5 ハッカソン削除
- **エンドポイント**: `DELETE /api/v1/hackathons/{id}`
- **レスポンス**:
```json
{
  "message": "Hackathon deleted successfully"
}
```
- **エラー**:
  - `404`: ハッカソンが見つからない

### 8.6 ハッカソン参加登録
- **エンドポイント**: `POST /api/v1/hackathons/{id}/participants`
- **リクエストボディ**:
```json
{
  "user_id": 1,
  "team_name": "Innovative Coders",
  "role": "Backend Developer",
  "notes": "React Native経験豊富"
}
```
- **レスポンス**: 作成された参加者情報（status: "registered"で作成）
- **エラー**:
  - `400`: 登録期間外、定員超過
  - `404`: ハッカソンが見つからない

### 8.7 ハッカソン参加者一覧取得
- **エンドポイント**: `GET /api/v1/hackathons/{id}/participants`
- **クエリパラメータ**:
  - `page`: ページ番号（デフォルト: 1）
  - `limit`: 取得件数（デフォルト: 10、最大: 100）
  - `status`: 参加ステータスでフィルタ
  - `team_name`: チーム名でフィルタ（部分一致）
- **レスポンス**:
```json
{
  "participants": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25
  }
}
```

### 8.8 ハッカソン参加者更新
- **エンドポイント**: `PUT /api/v1/hackathons/{id}/participants/{participant_id}`
- **リクエストボディ**: 更新したいフィールドのみ送信
```json
{
  "team_name": "新しいチーム名",
  "role": "Team Leader", 
  "status": "confirmed",
  "notes": "更新されたメモ"
}
```
- **ステータス**: `registered`, `confirmed`, `cancelled`, `disqualified`

### 8.9 ハッカソン参加者削除
- **エンドポイント**: `DELETE /api/v1/hackathons/{id}/participants/{participant_id}`
- **レスポンス**:
```json
{
  "message": "Participant deleted successfully"
}
```
- **エラー**:
  - `404`: 参加者が見つからない

## エラーコード一覧

| HTTPステータス | 説明 |
|---------------|------|
| 200 | OK - 成功 |
| 201 | Created - リソース作成成功 |
| 204 | No Content - OPTIONSリクエスト |
| 400 | Bad Request - 無効なリクエスト |
| 404 | Not Found - リソースが見つからない |
| 409 | Conflict - リソースの競合 |
| 500 | Internal Server Error - サーバーエラー |

## バリデーション

### ファイル
- ファイル名は空文字禁止
- ファイルサイズ制限あり（設定による）

### コンテスト
- quota値は0以上の整数
- 締切日は有効なRFC3339形式
- author_idは必須

### ハッカソン
- 日付の論理チェック:
  - 終了日 > 開始日
  - 登録締切 < 開始日
  - 登録開始 < 登録締切
- max_participantsは0以上
- statusは指定された値のみ（upcoming, ongoing, completed, cancelled）

### ブックマーク
- 自分自身のブックマーク禁止
- user_idとbookmarked_user_idは必須

## セキュリティ

### CORS設定
- すべてのオリジンを許可: `*`
- 許可メソッド: `GET, POST, PUT, DELETE, OPTIONS`
- 許可ヘッダー: `Origin, Content-Type, Accept, Authorization`

### 今後の拡張予定
- JWT認証の実装
- APIキーベース認証
- レート制限
- ファイルアップロードサイズ制限
- 詳細なアクセス権限制御