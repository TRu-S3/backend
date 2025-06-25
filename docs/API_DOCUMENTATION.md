# TRu-S3 API ドキュメント

## 概要

TRu-S3は、ファイル管理、コンテスト管理、ブックマーク機能を提供するRESTful APIサーバーです。Go言語で実装され、PostgreSQLデータベースとGoogle Cloud Storageを使用しています。

## 基本情報

- **ベースURL**: `/api/v1`
- **データ形式**: JSON
- **認証**: Google Cloud Platform認証
- **文字エンコーディング**: UTF-8

## 共通仕様

### HTTPステータスコード

| ステータスコード | 説明 |
|---|---|
| 200 | リクエスト成功 |
| 201 | リソース作成成功 |
| 400 | リクエストエラー（バリデーションエラー） |
| 404 | リソースが見つからない |
| 409 | リソースの競合（重複など） |
| 500 | サーバーエラー |

### エラーレスポンス形式

```json
{
  "error": "エラーメッセージ"
}
```

### ページネーション

一覧取得APIでは、以下のクエリパラメータでページネーション可能です：

- `page`: ページ番号（デフォルト: 1）
- `limit`: 1ページあたりの件数（デフォルト: 10、最大: 100）

レスポンス形式：

```json
{
  "データ": [...],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 50
  }
}
```

---

## 1. ファイル管理API

### 1.1 ファイルアップロード

**エンドポイント**: `POST /api/v1/files`

**説明**: 新しいファイルをアップロードします。

**リクエスト形式**: `multipart/form-data`

**リクエストパラメータ**:
- `file` (必須): アップロードするファイル

**レスポンス例**:
```json
{
  "id": "example.txt",
  "name": "example.txt",
  "path": "files/example.txt",
  "size": 1024,
  "content_type": "text/plain",
  "created_at": 1625097600,
  "updated_at": 1625097600
}
```

**エラーケース**:
- `400`: ファイルが提供されていない、マルチパートフォームの解析に失敗
- `409`: 同名のファイルが既に存在
- `500`: ファイルの作成に失敗

### 1.2 ファイル一覧取得

**エンドポイント**: `GET /api/v1/files`

**説明**: ファイルの一覧を取得します。

**クエリパラメータ**:
- `prefix` (オプション): ファイル名のプレフィックスでフィルタリング
- `limit` (オプション): 取得件数（デフォルト: 100）
- `offset` (オプション): 開始位置（デフォルト: 0）

**レスポンス例**:
```json
{
  "files": [
    {
      "id": "example1.txt",
      "name": "example1.txt",
      "path": "files/example1.txt",
      "size": 1024,
      "content_type": "text/plain",
      "created_at": 1625097600,
      "updated_at": 1625097600
    }
  ],
  "count": 1
}
```

### 1.3 ファイル情報取得

**エンドポイント**: `GET /api/v1/files/:id`

**説明**: 指定されたファイルの情報を取得します。

**パスパラメータ**:
- `id` (必須): ファイルID

**レスポンス例**:
```json
{
  "id": "example.txt",
  "name": "example.txt",
  "path": "files/example.txt",
  "size": 1024,
  "content_type": "text/plain",
  "created_at": 1625097600,
  "updated_at": 1625097600
}
```

**エラーケース**:
- `404`: ファイルが見つからない

### 1.4 ファイルダウンロード

**エンドポイント**: `GET /api/v1/files/:id/download`

**説明**: 指定されたファイルをダウンロードします。

**パスパラメータ**:
- `id` (必須): ファイルID

**レスポンス**: ファイルのバイナリデータ

**レスポンスヘッダー**:
- `Content-Disposition`: attachment; filename=ファイル名
- `Content-Type`: ファイルのMIMEタイプ
- `Content-Length`: ファイルサイズ

### 1.5 ファイル更新

**エンドポイント**: `PUT /api/v1/files/:id`

**説明**: 既存のファイルを更新します。

**パスパラメータ**:
- `id` (必須): ファイルID

**リクエスト形式**: `multipart/form-data`

**リクエストパラメータ**:
- `name` (オプション): 新しいファイル名
- `file` (オプション): 新しいファイル

**レスポンス例**:
```json
{
  "id": "example.txt",
  "name": "new-example.txt",
  "path": "files/new-example.txt",
  "size": 2048,
  "content_type": "text/plain",
  "created_at": 1625097600,
  "updated_at": 1625184000
}
```

**エラーケース**:
- `400`: 無効なファイル名
- `404`: ファイルが見つからない

### 1.6 ファイル削除

**エンドポイント**: `DELETE /api/v1/files/:id`

**説明**: 指定されたファイルを削除します。

**パスパラメータ**:
- `id` (必須): ファイルID

**レスポンス例**:
```json
{
  "message": "File deleted successfully"
}
```

**エラーケース**:
- `404`: ファイルが見つからない

---

## 2. コンテスト管理API

### 2.1 コンテスト作成

**エンドポイント**: `POST /api/v1/contests`

**説明**: 新しいコンテストを作成します。

**リクエスト例**:
```json
{
  "backend_quota": 2,
  "frontend_quota": 1,
  "ai_quota": 1,
  "application_deadline": "2024-12-31T23:59:59Z",
  "purpose": "Webアプリケーション開発コンテスト",
  "message": "一緒に素晴らしいWebアプリを作りましょう！",
  "author_id": 1
}
```

**リクエストフィールド**:
- `backend_quota` (必須): バックエンド募集人数（0以上）
- `frontend_quota` (必須): フロントエンド募集人数（0以上）
- `ai_quota` (必須): AI募集人数（0以上）
- `application_deadline` (必須): 募集期限（RFC3339形式）
- `purpose` (必須): コンテストの目的
- `message` (必須): メッセージ
- `author_id` (必須): 投稿者ID

**レスポンス例**:
```json
{
  "id": 1,
  "backend_quota": 2,
  "frontend_quota": 1,
  "ai_quota": 1,
  "application_deadline": "2024-12-31T23:59:59Z",
  "purpose": "Webアプリケーション開発コンテスト",
  "message": "一緒に素晴らしいWebアプリを作りましょう！",
  "author_id": 1,
  "created_at": "2024-06-25T10:00:00Z",
  "updated_at": "2024-06-25T10:00:00Z"
}
```

**エラーケース**:
- `400`: バリデーションエラー（必須フィールドの不足、形式エラー）

### 2.2 コンテスト一覧取得

**エンドポイント**: `GET /api/v1/contests`

**説明**: コンテストの一覧を取得します。

**クエリパラメータ**:
- `page` (オプション): ページ番号（デフォルト: 1）
- `limit` (オプション): 1ページあたりの件数（デフォルト: 10、最大: 100）
- `author_id` (オプション): 投稿者IDでフィルタリング
- `active` (オプション): "true"の場合、締切前のコンテストのみ取得

**レスポンス例**:
```json
{
  "contests": [
    {
      "id": 1,
      "backend_quota": 2,
      "frontend_quota": 1,
      "ai_quota": 1,
      "application_deadline": "2024-12-31T23:59:59Z",
      "purpose": "Webアプリケーション開発コンテスト",
      "message": "一緒に素晴らしいWebアプリを作りましょう！",
      "author_id": 1,
      "created_at": "2024-06-25T10:00:00Z",
      "updated_at": "2024-06-25T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 15
  }
}
```

### 2.3 コンテスト詳細取得

**エンドポイント**: `GET /api/v1/contests/:id`

**説明**: 指定されたコンテストの詳細を取得します。

**パスパラメータ**:
- `id` (必須): コンテストID

**レスポンス例**:
```json
{
  "id": 1,
  "backend_quota": 2,
  "frontend_quota": 1,
  "ai_quota": 1,
  "application_deadline": "2024-12-31T23:59:59Z",
  "purpose": "Webアプリケーション開発コンテスト",
  "message": "一緒に素晴らしいWebアプリを作りましょう！",
  "author_id": 1,
  "created_at": "2024-06-25T10:00:00Z",
  "updated_at": "2024-06-25T10:00:00Z"
}
```

**エラーケース**:
- `404`: コンテストが見つからない

### 2.4 コンテスト更新

**エンドポイント**: `PUT /api/v1/contests/:id`

**説明**: 既存のコンテストを更新します。

**パスパラメータ**:
- `id` (必須): コンテストID

**リクエスト例**:
```json
{
  "backend_quota": 3,
  "purpose": "Webアプリケーション開発コンテスト（更新版）",
  "application_deadline": "2024-12-31T23:59:59Z"
}
```

**リクエストフィールド**（すべてオプション）:
- `backend_quota`: バックエンド募集人数（0以上）
- `frontend_quota`: フロントエンド募集人数（0以上）
- `ai_quota`: AI募集人数（0以上）
- `application_deadline`: 募集期限（RFC3339形式）
- `purpose`: コンテストの目的
- `message`: メッセージ

**レスポンス例**:
```json
{
  "id": 1,
  "backend_quota": 3,
  "frontend_quota": 1,
  "ai_quota": 1,
  "application_deadline": "2024-12-31T23:59:59Z",
  "purpose": "Webアプリケーション開発コンテスト（更新版）",
  "message": "一緒に素晴らしいWebアプリを作りましょう！",
  "author_id": 1,
  "created_at": "2024-06-25T10:00:00Z",
  "updated_at": "2024-06-25T12:00:00Z"
}
```

**エラーケース**:
- `400`: バリデーションエラー
- `404`: コンテストが見つからない

### 2.5 コンテスト削除

**エンドポイント**: `DELETE /api/v1/contests/:id`

**説明**: 指定されたコンテストを削除します。

**パスパラメータ**:
- `id` (必須): コンテストID

**レスポンス例**:
```json
{
  "message": "Contest deleted successfully"
}
```

**エラーケース**:
- `404`: コンテストが見つからない

---

## 3. ブックマーク管理API

### 3.1 ブックマーク作成

**エンドポイント**: `POST /api/v1/bookmarks`

**説明**: 新しいブックマークを作成します。

**リクエスト例**:
```json
{
  "user_id": 1,
  "bookmarked_user_id": 2
}
```

**リクエストフィールド**:
- `user_id` (必須): ブックマークするユーザーID
- `bookmarked_user_id` (必須): ブックマークされるユーザーID

**レスポンス例**:
```json
{
  "id": 1,
  "user_id": 1,
  "bookmarked_user_id": 2,
  "created_at": "2024-06-25T10:00:00Z"
}
```

**エラーケース**:
- `400`: 自分自身をブックマークしようとした場合
- `500`: ブックマークの作成に失敗

### 3.2 ブックマーク一覧取得

**エンドポイント**: `GET /api/v1/bookmarks`

**説明**: ブックマークの一覧を取得します。

**クエリパラメータ**:
- `page` (オプション): ページ番号（デフォルト: 1）
- `limit` (オプション): 1ページあたりの件数（デフォルト: 10、最大: 100）
- `user_id` (オプション): 特定のユーザーのブックマークのみ取得

**レスポンス例**:
```json
{
  "bookmarks": [
    {
      "id": 1,
      "user_id": 1,
      "bookmarked_user_id": 2,
      "created_at": "2024-06-25T10:00:00Z",
      "User": {
        "id": 1,
        "name": "山田太郎",
        "gmail": "yamada@example.com"
      },
      "BookmarkedUser": {
        "id": 2,
        "name": "佐藤花子",
        "gmail": "sato@example.com"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 5
  }
}
```

### 3.3 ブックマーク削除

**エンドポイント**: `DELETE /api/v1/bookmarks/:id`

**説明**: 指定されたブックマークを削除します。

**パスパラメータ**:
- `id` (必須): ブックマークID

**レスポンス例**:
```json
{
  "message": "Bookmark deleted successfully"
}
```

**エラーケース**:
- `404`: ブックマークが見つからない

---

## 4. ハッカソン管理API

### 4.1 ハッカソン作成

**エンドポイント**: `POST /api/v1/hackathons`

**説明**: 新しいハッカソンを作成します。

**リクエスト例**:
```json
{
  "name": "AI Hackathon with Google Cloud",
  "description": "Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。",
  "start_date": "2024-08-15T09:00:00+09:00",
  "end_date": "2024-08-17T18:00:00+09:00",
  "registration_start": "2024-07-01T00:00:00+09:00",
  "registration_deadline": "2024-08-10T23:59:59+09:00",
  "max_participants": 100,
  "location": "Google Cloud Tokyo オフィス",
  "organizer": "Google Cloud Japan",
  "contact_email": "hackathon-ai@googlecloud.com",
  "prize_info": "優勝: Google Cloud クレジット $5,000 + Google Pixel",
  "rules": "・Google Cloud技術の使用必須\n・チーム人数は2-5人\n・オリジナル作品であること",
  "tech_stack": "[\"Google Cloud Vertex AI\", \"Gemini API\", \"Python\", \"JavaScript\"]",
  "is_public": true,
  "banner_url": "https://example.com/banner.jpg",
  "website_url": "https://cloud.google.com/events/hackathon-ai"
}
```

**リクエストフィールド**:
- `name` (必須): ハッカソン名
- `description` (オプション): 説明
- `start_date` (必須): 開始日時（RFC3339形式）
- `end_date` (必須): 終了日時（RFC3339形式）
- `registration_start` (必須): 参加登録開始日時（RFC3339形式）
- `registration_deadline` (必須): 参加登録締切日時（RFC3339形式）
- `max_participants` (オプション): 最大参加者数（0は無制限、デフォルト: 0）
- `location` (オプション): 開催場所
- `organizer` (必須): 主催者
- `contact_email` (オプション): 連絡先メールアドレス
- `prize_info` (オプション): 賞品・賞金情報
- `rules` (オプション): ルール・規則
- `tech_stack` (オプション): 推奨技術スタック（JSON文字列形式）
- `is_public` (オプション): 公開/非公開（デフォルト: true）
- `banner_url` (オプション): バナー画像URL
- `website_url` (オプション): 公式ウェブサイトURL

**レスポンス例**:
```json
{
  "id": 1,
  "name": "AI Hackathon with Google Cloud",
  "description": "Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。",
  "start_date": "2024-08-15T09:00:00+09:00",
  "end_date": "2024-08-17T18:00:00+09:00",
  "registration_start": "2024-07-01T00:00:00+09:00",
  "registration_deadline": "2024-08-10T23:59:59+09:00",
  "max_participants": 100,
  "location": "Google Cloud Tokyo オフィス",
  "organizer": "Google Cloud Japan",
  "contact_email": "hackathon-ai@googlecloud.com",
  "prize_info": "優勝: Google Cloud クレジット $5,000 + Google Pixel",
  "rules": "・Google Cloud技術の使用必須\n・チーム人数は2-5人\n・オリジナル作品であること",
  "tech_stack": "[\"Google Cloud Vertex AI\", \"Gemini API\", \"Python\", \"JavaScript\"]",
  "status": "upcoming",
  "is_public": true,
  "banner_url": "https://example.com/banner.jpg",
  "website_url": "https://cloud.google.com/events/hackathon-ai",
  "created_at": "2024-06-25T10:00:00Z",
  "updated_at": "2024-06-25T10:00:00Z"
}
```

**エラーケース**:
- `400`: バリデーションエラー（必須フィールドの不足、日付形式エラー、日付の論理エラー）

### 4.2 ハッカソン一覧取得

**エンドポイント**: `GET /api/v1/hackathons`

**説明**: ハッカソンの一覧を取得します。

**クエリパラメータ**:
- `page` (オプション): ページ番号（デフォルト: 1）
- `limit` (オプション): 1ページあたりの件数（デフォルト: 10、最大: 100）
- `status` (オプション): ステータスでフィルタリング（upcoming, ongoing, completed, cancelled）
- `organizer` (オプション): 主催者名でフィルタリング（部分一致）
- `is_public` (オプション): 公開/非公開でフィルタリング（true/false）
- `upcoming` (オプション): "true"の場合、開始前のハッカソンのみ取得
- `ongoing` (オプション): "true"の場合、開催中のハッカソンのみ取得
- `registration_open` (オプション): "true"の場合、参加登録受付中のハッカソンのみ取得

**レスポンス例**:
```json
{
  "hackathons": [
    {
      "id": 1,
      "name": "AI Hackathon with Google Cloud",
      "description": "Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。",
      "start_date": "2024-08-15T09:00:00+09:00",
      "end_date": "2024-08-17T18:00:00+09:00",
      "registration_start": "2024-07-01T00:00:00+09:00",
      "registration_deadline": "2024-08-10T23:59:59+09:00",
      "max_participants": 100,
      "location": "Google Cloud Tokyo オフィス",
      "organizer": "Google Cloud Japan",
      "status": "upcoming",
      "is_public": true,
      "created_at": "2024-06-25T10:00:00Z",
      "updated_at": "2024-06-25T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 2
  }
}
```

### 4.3 ハッカソン詳細取得

**エンドポイント**: `GET /api/v1/hackathons/:id`

**説明**: 指定されたハッカソンの詳細情報と参加者一覧を取得します。

**パスパラメータ**:
- `id` (必須): ハッカソンID

**レスポンス例**:
```json
{
  "id": 1,
  "name": "AI Hackathon with Google Cloud",
  "description": "Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。",
  "start_date": "2024-08-15T09:00:00+09:00",
  "end_date": "2024-08-17T18:00:00+09:00",
  "registration_start": "2024-07-01T00:00:00+09:00",
  "registration_deadline": "2024-08-10T23:59:59+09:00",
  "max_participants": 100,
  "location": "Google Cloud Tokyo オフィス",
  "organizer": "Google Cloud Japan",
  "contact_email": "hackathon-ai@googlecloud.com",
  "prize_info": "優勝: Google Cloud クレジット $5,000 + Google Pixel",
  "rules": "・Google Cloud技術の使用必須\n・チーム人数は2-5人\n・オリジナル作品であること",
  "tech_stack": "[\"Google Cloud Vertex AI\", \"Gemini API\", \"Python\", \"JavaScript\"]",
  "status": "upcoming",
  "is_public": true,
  "banner_url": "https://example.com/banner.jpg",
  "website_url": "https://cloud.google.com/events/hackathon-ai",
  "created_at": "2024-06-25T10:00:00Z",
  "updated_at": "2024-06-25T10:00:00Z",
  "participants": [
    {
      "id": 1,
      "hackathon_id": 1,
      "user_id": 1,
      "team_name": "AI Innovators",
      "role": "developer",
      "registration_date": "2024-07-05T10:00:00Z",
      "status": "registered",
      "notes": "バックエンド開発経験あり",
      "user": {
        "id": 1,
        "name": "山田太郎",
        "gmail": "yamada@example.com"
      }
    }
  ]
}
```

**エラーケース**:
- `404`: ハッカソンが見つからない

### 4.4 ハッカソン更新

**エンドポイント**: `PUT /api/v1/hackathons/:id`

**説明**: 既存のハッカソンを更新します。

**パスパラメータ**:
- `id` (必須): ハッカソンID

**リクエスト例**:
```json
{
  "max_participants": 150,
  "status": "ongoing",
  "description": "Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。（更新版）"
}
```

**リクエストフィールド**（すべてオプション）:
- 作成時と同じフィールドがすべて更新可能
- `status`: ステータス（upcoming, ongoing, completed, cancelled）

**レスポンス例**:
```json
{
  "id": 1,
  "name": "AI Hackathon with Google Cloud",
  "description": "Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。（更新版）",
  "max_participants": 150,
  "status": "ongoing",
  "created_at": "2024-06-25T10:00:00Z",
  "updated_at": "2024-08-15T09:30:00Z"
}
```

**エラーケース**:
- `400`: バリデーションエラー（無効なステータス、日付の論理エラー）
- `404`: ハッカソンが見つからない

### 4.5 ハッカソン削除

**エンドポイント**: `DELETE /api/v1/hackathons/:id`

**説明**: 指定されたハッカソンを削除します。

**パスパラメータ**:
- `id` (必須): ハッカソンID

**レスポンス例**:
```json
{
  "message": "Hackathon deleted successfully"
}
```

**エラーケース**:
- `404`: ハッカソンが見つからない

### 4.6 ハッカソン参加登録

**エンドポイント**: `POST /api/v1/hackathons/:id/participants`

**説明**: 指定されたハッカソンに参加登録します。

**パスパラメータ**:
- `id` (必須): ハッカソンID

**リクエスト例**:
```json
{
  "user_id": 1,
  "team_name": "AI Innovators",
  "role": "developer",
  "notes": "バックエンド開発経験あり、PythonとTensorFlowが得意です。"
}
```

**リクエストフィールド**:
- `user_id` (必須): ユーザーID
- `team_name` (オプション): チーム名
- `role` (オプション): 役割（developer, designer, pm等）
- `notes` (オプション): 参加時のメモ・自己紹介

**レスポンス例**:
```json
{
  "id": 1,
  "hackathon_id": 1,
  "user_id": 1,
  "team_name": "AI Innovators",
  "role": "developer",
  "registration_date": "2024-07-05T10:00:00Z",
  "status": "registered",
  "notes": "バックエンド開発経験あり、PythonとTensorFlowが得意です。",
  "user": {
    "id": 1,
    "name": "山田太郎",
    "gmail": "yamada@example.com"
  },
  "hackathon": {
    "id": 1,
    "name": "AI Hackathon with Google Cloud"
  }
}
```

**エラーケース**:
- `400`: 参加登録期間外、定員超過、重複登録
- `404`: ハッカソンが見つからない

### 4.7 参加者一覧取得

**エンドポイント**: `GET /api/v1/hackathons/:id/participants`

**説明**: 指定されたハッカソンの参加者一覧を取得します。

**パスパラメータ**:
- `id` (必須): ハッカソンID

**クエリパラメータ**:
- `page` (オプション): ページ番号（デフォルト: 1）
- `limit` (オプション): 1ページあたりの件数（デフォルト: 10、最大: 100）
- `status` (オプション): ステータスでフィルタリング（registered, confirmed, cancelled, disqualified）
- `team_name` (オプション): チーム名でフィルタリング（部分一致）

**レスポンス例**:
```json
{
  "participants": [
    {
      "id": 1,
      "hackathon_id": 1,
      "user_id": 1,
      "team_name": "AI Innovators",
      "role": "developer",
      "registration_date": "2024-07-05T10:00:00Z",
      "status": "registered",
      "notes": "バックエンド開発経験あり",
      "user": {
        "id": 1,
        "name": "山田太郎",
        "gmail": "yamada@example.com"
      }
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 10,
    "total": 25
  }
}
```

### 4.8 参加者削除

**エンドポイント**: `DELETE /api/v1/hackathons/:id/participants/:participant_id`

**説明**: 指定された参加者をハッカソンから削除します。

**パスパラメータ**:
- `id` (必須): ハッカソンID
- `participant_id` (必須): 参加者ID

**レスポンス例**:
```json
{
  "message": "Participant deleted successfully"
}
```

**エラーケース**:
- `404`: 参加者が見つからない

---

## 5. データベーススキーマ

### 5.1 ファイルメタデータテーブル (file_metadata)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | VARCHAR(255) | PRIMARY KEY | ファイルID |
| name | VARCHAR(255) | NOT NULL | ファイル名 |
| path | VARCHAR(500) | NOT NULL | ファイルパス |
| size | BIGINT | NOT NULL | ファイルサイズ（バイト） |
| content_type | VARCHAR(100) | | MIMEタイプ |
| checksum | VARCHAR(64) | | チェックサム |
| tags | TEXT | | タグ（JSON文字列） |
| created_at | BIGINT | AUTO | 作成日時（Unix timestamp） |
| updated_at | BIGINT | AUTO | 更新日時（Unix timestamp） |

### 5.2 ユーザーテーブル (users)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | ユーザーID |
| gmail | VARCHAR(255) | UNIQUE, NOT NULL | Gmailアドレス |
| name | VARCHAR(100) | NOT NULL | ユーザー名 |
| icon_url | TEXT | | プロフィール画像URL |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 作成日時 |

### 5.3 コンテストテーブル (contests)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | コンテストID |
| backend_quota | INTEGER | NOT NULL, DEFAULT 0 | バックエンド募集人数 |
| frontend_quota | INTEGER | NOT NULL, DEFAULT 0 | フロントエンド募集人数 |
| ai_quota | INTEGER | NOT NULL, DEFAULT 0 | AI募集人数 |
| application_deadline | TIMESTAMP | NOT NULL | 募集期限 |
| purpose | TEXT | NOT NULL | 目的 |
| message | TEXT | NOT NULL | メッセージ |
| author_id | INTEGER | NOT NULL, FK | 投稿者ID |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

### 5.4 ハッカソンテーブル (hackathons)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | ハッカソンID |
| name | VARCHAR(255) | NOT NULL | ハッカソン名 |
| description | TEXT | | 説明 |
| start_date | TIMESTAMP | NOT NULL | 開始日時 |
| end_date | TIMESTAMP | NOT NULL | 終了日時 |
| registration_start | TIMESTAMP | NOT NULL | 参加登録開始日時 |
| registration_deadline | TIMESTAMP | NOT NULL | 参加登録締切日時 |
| max_participants | INTEGER | DEFAULT 0 | 最大参加者数（0は無制限） |
| location | VARCHAR(255) | | 開催場所 |
| organizer | VARCHAR(255) | NOT NULL | 主催者 |
| contact_email | VARCHAR(255) | | 連絡先メールアドレス |
| prize_info | TEXT | | 賞品・賞金情報 |
| rules | TEXT | | ルール・規則 |
| tech_stack | TEXT | | 推奨技術スタック（JSON形式） |
| status | VARCHAR(50) | DEFAULT upcoming | ステータス |
| is_public | BOOLEAN | DEFAULT true | 公開/非公開フラグ |
| banner_url | TEXT | | バナー画像URL |
| website_url | TEXT | | 公式ウェブサイトURL |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

### 5.5 ハッカソン参加者テーブル (hackathon_participants)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | 参加者ID |
| hackathon_id | INTEGER | NOT NULL, FK | ハッカソンID |
| user_id | INTEGER | NOT NULL, FK | ユーザーID |
| team_name | VARCHAR(255) | | チーム名 |
| role | VARCHAR(100) | | 役割 |
| registration_date | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 参加登録日時 |
| status | VARCHAR(50) | DEFAULT registered | ステータス |
| notes | TEXT | | 参加時のメモ・自己紹介 |

### 5.6 ブックマークテーブル (bookmarks)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | ブックマークID |
| user_id | INTEGER | NOT NULL, FK | ブックマークするユーザーID |
| bookmarked_user_id | INTEGER | NOT NULL, FK | ブックマークされるユーザーID |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 作成日時 |

### 5.7 プロフィールテーブル (profiles)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | プロフィールID |
| user_id | INTEGER | NOT NULL, FK, UNIQUE | ユーザーID |
| bio | TEXT | | 自己紹介文 |
| tag_id | INTEGER | FK | タグID |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

### 5.8 タグテーブル (tags)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | タグID |
| name | VARCHAR(50) | UNIQUE, NOT NULL | タグ名 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

### 5.9 マッチングテーブル (matchings)

| フィールド名 | 型 | 制約 | 説明 |
|---|---|---|---|
| id | SERIAL | PRIMARY KEY | マッチングID |
| user_id | INTEGER | NOT NULL, FK | ユーザーID |
| notify_id | INTEGER | | 通知ID |
| content | TEXT | | マッチング内容 |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 作成日時 |
| updated_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 更新日時 |

---

## 6. セキュリティ仕様

### 6.1 認证方式

- **Google Cloud Platform認証**: GCP IAMを使用
- **データベース接続**: Cloud SQL Auth Proxyによる暗号化通信

### 6.2 CORS設定

現在は開発用設定で全オリジンを許可しています。本番環境では適切なオリジンの制限が必要です。

### 6.3 入力検証

- **Ginフレームワーク**: `ShouldBindJSON`による自動バリデーション
- **データベース制約**: 外部キー制約、ユニーク制約、チェック制約
- **ファイル検証**: ファイルサイズ、ファイル名の妥当性チェック

---

## 7. 技術仕様

### 7.1 使用技術

- **言語**: Go 1.24.2
- **フレームワーク**: Gin（HTTPフレームワーク）
- **ORM**: GORM
- **データベース**: PostgreSQL
- **ストレージ**: Google Cloud Storage
- **デプロイ**: Google Cloud Run
- **コンテナ**: Docker

### 7.2 アーキテクチャ

オニオンアーキテクチャを採用し、以下の層で構成されています：

- **Domain層**: ビジネスロジック、エンティティ
- **Application層**: アプリケーションサービス
- **Infrastructure層**: 外部サービス（GCS、データベース）との連携
- **Interfaces層**: HTTPハンドラ、ルーティング

### 7.3 パフォーマンス最適化

- **データベースインデックス**: 検索頻度の高いフィールドにインデックス設定
- **ページネーション**: 大量データの効率的な取得
- **Connection Pooling**: データベース接続の最適化

---

## 8. 運用情報

### 8.1 ログ仕様

- **アプリケーションログ**: 標準出力への構造化ログ
- **エラーログ**: エラー発生時の詳細情報
- **アクセスログ**: HTTPリクエスト/レスポンス情報

### 8.2 監視項目

- **HTTPステータス**: 4xx/5xxエラーの監視
- **レスポンス時間**: API応答時間の監視
- **データベース**: 接続数、クエリ性能
- **ストレージ**: GCSの使用量、転送量

### 8.3 バックアップ

- **データベース**: Cloud SQLの自動バックアップ
- **ファイル**: Google Cloud Storageの多重化

---

## 9. 開発者向け情報

### 9.1 ローカル開発環境

```bash
# Docker Composeを使用した開発環境の起動
docker-compose up -d

# アプリケーションの起動
go run main.go
```

### 9.2 テスト実行

```bash
# ユニットテストの実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...
```

### 9.3 データベースマイグレーション

```bash
# マイグレーションの実行
psql -h localhost -U postgres -d trusts3 -f migrations/001_create_matching_app_tables.sql
psql -h localhost -U postgres -d trusts3 -f migrations/002_create_contests_table.sql
psql -h localhost -U postgres -d trusts3 -f migrations/003_create_hackathons_table.sql
```

---

## 10. API使用例

### 10.1 cURLでのサンプルリクエスト

#### ファイルアップロード
```bash
curl -X POST \
  http://localhost:8080/api/v1/files \
  -H 'Content-Type: multipart/form-data' \
  -F 'file=@example.txt'
```

#### コンテスト作成
```bash
curl -X POST \
  http://localhost:8080/api/v1/contests \
  -H 'Content-Type: application/json' \
  -d '{
    "backend_quota": 2,
    "frontend_quota": 1,
    "ai_quota": 1,
    "application_deadline": "2024-12-31T23:59:59Z",
    "purpose": "Webアプリケーション開発コンテスト",
    "message": "一緒に素晴らしいWebアプリを作りましょう！",
    "author_id": 1
  }'
```

#### ブックマーク作成
```bash
curl -X POST \
  http://localhost:8080/api/v1/bookmarks \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": 1,
    "bookmarked_user_id": 2
  }'
```

#### ハッカソン作成
```bash
curl -X POST \
  http://localhost:8080/api/v1/hackathons \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "AI Hackathon with Google Cloud",
    "description": "Google CloudのAI/ML技術を活用した革新的なアプリケーションを開発するハッカソンです。",
    "start_date": "2024-08-15T09:00:00+09:00",
    "end_date": "2024-08-17T18:00:00+09:00",
    "registration_start": "2024-07-01T00:00:00+09:00",
    "registration_deadline": "2024-08-10T23:59:59+09:00",
    "max_participants": 100,
    "location": "Google Cloud Tokyo オフィス",
    "organizer": "Google Cloud Japan"
  }'
```

#### ハッカソン参加登録
```bash
curl -X POST \
  http://localhost:8080/api/v1/hackathons/1/participants \
  -H 'Content-Type: application/json' \
  -d '{
    "user_id": 1,
    "team_name": "AI Innovators",
    "role": "developer",
    "notes": "バックエンド開発経験あり"
  }'
```

### 10.2 JavaScriptでの使用例

```javascript
// ファイルアップロード
const formData = new FormData();
formData.append('file', fileInput.files[0]);

fetch('/api/v1/files', {
  method: 'POST',
  body: formData
})
.then(response => response.json())
.then(data => console.log(data));

// コンテスト一覧取得
fetch('/api/v1/contests?page=1&limit=10')
.then(response => response.json())
.then(data => console.log(data));

// ハッカソン一覧取得
fetch('/api/v1/hackathons?status=upcoming&registration_open=true')
.then(response => response.json())
.then(data => console.log(data));

// ハッカソン参加登録
fetch('/api/v1/hackathons/1/participants', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    user_id: 1,
    team_name: 'AI Innovators',
    role: 'developer',
    notes: 'バックエンド開発経験あり'
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

---

## 11. 変更履歴

| 日付 | バージョン | 変更内容 |
|---|---|---|
| 2024-06-25 | 1.0.0 | 初版作成 |
| 2024-06-25 | 1.1.0 | ハッカソン管理API追加 |

---

## 12. 問い合わせ先

技術的な質問や不具合報告については、プロジェクトのGitHubリポジトリのIssuesをご利用ください。

---

*最終更新日: 2024年6月25日*