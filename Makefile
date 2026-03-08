.PHONY: help dev-up dev-down dev-logs test-setup test-teardown test-backend test-e2e test deploy-frontend

# デフォルトターゲット
help: ## 全コマンドの一覧を表示
	@echo "=========================================="
	@echo " SNS Application - Makefile Commands"
	@echo "=========================================="
	@echo ""
	@echo "開発環境:"
	@echo "  make dev-up        - 開発環境を起動"
	@echo "  make dev-down      - 開発環境を停止"
	@echo "  make dev-logs      - 開発環境のログを表示"
	@echo ""
	@echo "テスト環境:"
	@echo "  make test-setup    - テスト環境を起動"
	@echo "  make test-teardown - テスト環境を停止・削除"
	@echo ""
	@echo "テスト実行:"
	@echo "  make test-backend  - バックエンドの単体テストを実行"
	@echo "  make test-e2e      - フロントエンドのE2Eテストを実行"
	@echo "  make test          - 全テストを実行（backend + e2e）"
	@echo ""
	@echo "デプロイ:"
	@echo "  make deploy-frontend - フロントエンドをFirebase Hostingにデプロイ"
	@echo ""
	@echo "=========================================="

# 開発環境
dev-up: ## 開発環境を起動
	@echo "🚀 開発環境を起動中..."
	docker compose --profile dev up -d
	@echo "✅ 開発環境が起動しました"
	@echo "   API: http://localhost:8080"
	@echo "   Swagger: http://localhost:8080/swagger/index.html"

dev-down: ## 開発環境を停止
	@echo "🛑 開発環境を停止中..."
	docker compose --profile dev down
	@echo "✅ 開発環境が停止しました"

dev-logs: ## 開発環境のログを表示
	docker compose --profile dev logs -f

# テスト環境
test-setup: ## テスト環境を起動
	@echo "🔧 テスト環境をセットアップ中..."
	docker compose --profile test up -d
	@echo "⏳ テスト用APIサーバーの起動を待機中（最大60秒）..."
	@for i in 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15; do \
		if curl -s http://localhost:8081/health > /dev/null 2>&1; then \
			echo "✅ テスト環境のセットアップが完了しました ($$i秒)"; \
			echo "   Test API: http://localhost:8081"; \
			exit 0; \
		fi; \
		if [ $$i -eq 1 ]; then echo "   初回起動時は依存関係のダウンロードとビルドに時間がかかります..."; fi; \
		sleep 4; \
	done; \
	echo "❌ APIサーバーの起動に失敗しました（60秒タイムアウト）"; \
	echo "   ログを確認してください: docker compose logs api_test"; \
	exit 1

test-teardown: ## テスト環境を停止・削除
	@echo "🧹 テスト環境をクリーンアップ中..."
	docker compose --profile test down
	@echo "✅ テスト環境をクリーンアップしました"
	@echo "   ※ テスト用DBはtmpfs（メモリ）を使用しているため、コンテナ停止で自動削除されます"

# バックエンドテスト
test-backend: ## バックエンドの単体テストを実行
	@echo "=========================================="
	@echo " バックエンド単体テスト実行"
	@echo "=========================================="
	@$(MAKE) test-setup
	@echo ""
	@echo "🧪 単体テストを実行中..."
	@echo ""
	@docker compose run --rm api_test go test -v -cover -p 1 ./... || (echo "❌ テストが失敗しました" && $(MAKE) test-teardown && exit 1)
	@echo ""
	@echo "=========================================="
	@echo "✅ バックエンドテストが完了しました"
	@echo "=========================================="
	@$(MAKE) test-teardown

# E2Eテスト
test-e2e: ## フロントエンドのE2Eテストを実行
	@echo "=========================================="
	@echo " E2Eテスト実行"
	@echo "=========================================="
	@echo ""
	@echo "⚠️  重要: E2Eテストを実行する前に、別ターミナルで"
	@echo "   フロントエンド開発サーバーを**テストモード**で起動してください:"
	@echo ""
	@echo "   cd frontend && npm run dev:test"
	@echo ""
	@echo "   テストサーバー: http://localhost:3001 (ポート3001)"
	@echo "   開発サーバー:   http://localhost:3000 (ポート3000)"
	@echo "   ※ 両方を同時に起動できます"
	@echo ""
	@read -p "フロントエンドサーバー（テストモード: 3001）が起動していますか？ (y/N): " confirm && [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ] || (echo "キャンセルしました" && exit 1)
	@$(MAKE) test-setup
	@echo ""
	@echo "🎭 Playwrightテストを実行中..."
	@echo ""
	@cd frontend && npm run test:e2e || (echo "❌ E2Eテストが失敗しました" && cd .. && $(MAKE) test-teardown && exit 1)
	@echo ""
	@echo "=========================================="
	@echo "✅ E2Eテストが完了しました"
	@echo "=========================================="
	@$(MAKE) test-teardown

# 全テスト実行
test: ## 全テストを実行（backend + e2e）
	@echo "=========================================="
	@echo " 全テスト実行"
	@echo "=========================================="
	@echo ""
	@$(MAKE) test-backend
	@echo ""
	@$(MAKE) test-e2e
	@echo ""
	@echo "=========================================="
	@echo "🎉 全テストが完了しました！"
	@echo "=========================================="

# デプロイ
deploy-frontend: ## フロントエンドをFirebase Hostingにデプロイ
	@echo "=========================================="
	@echo " Firebase Hosting デプロイ"
	@echo "=========================================="
	@cd frontend && ./scripts/deploy.sh
