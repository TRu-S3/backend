# 🔐 TRu-S3 Cloud SQL セキュリティガイド

## 📋 概要

このドキュメントでは、TRu-S3プロジェクトでCloud SQLに安全に接続するための3つのアプローチを説明します。

## 🛠️ セットアップオプション

### 1. 基本セットアップ（開発用）
```bash
# 現在の設定（パスワード認証）
cp .env.gcp.backup .env
```

### 2. Cloud SQL Proxy セットアップ
```bash
# プロキシ経由での安全な接続
./scripts/setup-cloud-sql-proxy.sh
./start-proxy.sh
cp .env.proxy .env
```

### 3. IAM認証セットアップ（推奨）
```bash
# IAM認証による最高レベルのセキュリティ
./scripts/setup-iam-auth.sh
cp .env.iam .env
```

### 4. フル・セキュリティセットアップ（本番用）
```bash
# SSL/TLS + プライベートIP + VPCコネクタ
./scripts/setup-network-security.sh
cp .env.secure .env
```

## 🔍 接続テスト

各設定をテストするには：

```bash
# 全ての設定をテスト
go run test_secure_connection.go

# 特定の設定のみをテスト
godotenv -f .env.secure go run test_secure_connection.go
```

## 📊 セキュリティレベル比較

| 方法 | セキュリティ | 複雑さ | 用途 |
|------|-------------|--------|------|
| 直接接続 | ⭐⭐ | ⭐ | 開発 |
| Cloud SQL Proxy | ⭐⭐⭐⭐ | ⭐⭐ | ステージング |
| IAM認証 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | 本番 |
| フルセキュリティ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | エンタープライズ |

## 🚀 本番デプロイメント

本番環境では以下の手順を推奨：

1. **IAM認証の有効化**
2. **SSL/TLS証明書の設定**
3. **プライベートIPの使用**
4. **VPCコネクタの設定**
5. **Docker Composeによるデプロイ**

```bash
# 本番デプロイ
docker-compose -f docker-compose.prod.yml up -d
```

## ⚠️ 重要な注意事項

- **service-account.json** は絶対にGitにコミットしない
- **SSL証明書** は安全な場所に保管
- **パスワード** は環境変数で管理
- **プライベートIP** を可能な限り使用

## 🔧 トラブルシューティング

### 接続タイムアウト
- Cloud SQL Proxyの使用を検討
- ファイアウォール設定の確認

### 認証エラー
- サービスアカウントの権限確認
- IAMポリシーの見直し

### SSL証明書エラー
- 証明書の有効期限確認
- SSL設定の見直し