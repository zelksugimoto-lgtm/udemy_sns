#!/bin/bash

# ====================
# Firebase Hosting デプロイスクリプト
# ====================

set -e  # エラーが発生したら即座に終了

echo "=========================================="
echo " Firebase Hosting デプロイスクリプト"
echo "=========================================="
echo ""

# カレントディレクトリをfrontendに移動
cd "$(dirname "$0")/.."

# Firebase CLIがインストールされているか確認
if ! command -v firebase &> /dev/null; then
    echo "❌ Firebase CLIがインストールされていません"
    echo ""
    echo "以下のコマンドでインストールしてください:"
    echo "  npm install -g firebase-tools"
    echo ""
    exit 1
fi

# Firebase ログイン確認
echo "🔐 Firebaseへのログイン状態を確認中..."
if ! firebase projects:list &> /dev/null; then
    echo "❌ Firebaseにログインしていません"
    echo ""
    echo "以下のコマンドでログインしてください:"
    echo "  firebase login"
    echo ""
    exit 1
fi
echo "✅ Firebaseにログイン済み"
echo ""

# .env.productionファイルの存在確認
if [ ! -f ".env.production" ]; then
    echo "❌ .env.productionファイルが見つかりません"
    echo ""
    echo "frontend/.env.productionを作成してください"
    exit 1
fi

# .env.productionの内容を確認
echo "📄 本番環境設定を確認中..."
cat .env.production
echo ""

# APIのURLが設定されているか確認
if grep -q "your-backend-service-xxxxx" .env.production; then
    echo "⚠️  警告: VITE_API_BASE_URLがプレースホルダーのままです"
    echo ""
    echo "デプロイ前に.env.productionを編集して、"
    echo "実際のCloud Run URLに変更してください。"
    echo ""
    read -p "このまま続行しますか？ (y/N): " confirm
    if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
        echo "デプロイをキャンセルしました"
        exit 0
    fi
fi

# ビルド実行
echo "🔨 本番用ビルドを実行中..."
npm run build
echo "✅ ビルドが完了しました"
echo ""

# デプロイ確認
echo "🚀 Firebase Hostingへのデプロイを開始します"
echo ""
read -p "デプロイを続行しますか？ (y/N): " confirm
if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
    echo "デプロイをキャンセルしました"
    exit 0
fi

# Firebase Hostingへデプロイ
echo ""
echo "📤 Firebase Hostingへデプロイ中..."
firebase deploy --only hosting

echo ""
echo "=========================================="
echo "✅ デプロイが完了しました！"
echo "=========================================="
echo ""
echo "デプロイされたサイト:"
firebase hosting:sites:list
