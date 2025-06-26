# GitHub ワークフロー設定

## 📁 構成

この `.github` ディレクトリには、GitHub Actions ワークフローやプロジェクトテンプレートを配置します。

### 🔄 推奨ワークフロー

```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.24.2
      - run: go test ./...
      
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: docker build -t tru-s3:latest .
```

### 📝 Issue テンプレート

- バグレポート
- 機能要求
- セキュリティ問題

### 🔀 PR テンプレート

- 変更内容の説明
- テスト方法
- チェックリスト