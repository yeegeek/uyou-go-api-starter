.PHONY: help quick-start up down restart logs build test test-coverage lint lint-fix swag migrate-create migrate-up migrate-down migrate-status migrate-goto migrate-force migrate-drop build-binary run-binary clean generate-jwt-secret check-env

# Container name (from docker-compose.yml)
CONTAINER_NAME := go_api_app

# Check if container is running
CONTAINER_RUNNING := $(shell docker ps --format '{{.Names}}' 2>/dev/null | grep -E '^$(CONTAINER_NAME)$$')

# Determine execution command
ifdef CONTAINER_RUNNING
	EXEC_CMD = docker exec $(CONTAINER_NAME)
	EXEC_CMD_INTERACTIVE = docker exec -i $(CONTAINER_NAME)
	ENV_MSG = 🐳 Running in Docker container
else
	EXEC_CMD = 
	EXEC_CMD_INTERACTIVE = 
	ENV_MSG = 💻 Running on host (Docker not available)
endif

## help: 显示帮助信息
help:
	@echo "Go REST API 脚手架 - 可用命令"
	@echo "=============================================="
	@echo ""
	@echo "🚀 快速开始:"
	@echo "  make quick-start    - 完整设置并启动（需要 Docker）"
	@echo ""
	@echo "🐳 Docker 命令:"
	@echo "  make up             - 启动容器"
	@echo "  make down           - 停止容器"
	@echo "  make restart        - 重启容器"
	@echo "  make logs           - 查看容器日志"
	@echo "  make build          - 重新构建容器"
	@echo ""
	@echo "🧪 开发命令:"
	@echo "  make test           - 运行测试"
	@echo "  make test-coverage  - 运行测试并生成覆盖率报告"
	@echo "  make lint           - 运行代码检查"
	@echo "  make lint-fix       - 运行代码检查并自动修复问题"
	@echo "  make swag           - 生成 Swagger 文档"
	@echo ""
	@echo "🔒 安全命令:"
	@echo "  make generate-jwt-secret  - 生成并设置 JWT 密钥到 .env"
	@echo "  make check-env            - 检查必需的环境变量"
	@echo ""
	@echo "👤 管理员管理:"
	@echo "  make create-admin         - 创建新管理员用户（交互式）"
	@echo "  make promote-admin ID=<n> - 将现有用户提升为管理员"
	@echo ""
	@echo "📊️  数据库命令:"
	@echo "  make migrate-create NAME=<name>  - 创建新迁移"
	@echo "  make migrate-up                  - 应用所有待执行的迁移"
	@echo "  make migrate-down                - 回滚最后一次迁移（或使用 STEPS=N 回滚 N 次）"
	@echo "  make migrate-status              - 显示当前迁移版本"
	@echo "  make migrate-goto VERSION=<n>    - 跳转到指定版本"
	@echo "  make migrate-force VERSION=<n>   - 强制设置版本（恢复用）"
	@echo "  make migrate-drop                - 删除所有表"
	@echo ""
	@echo "⏰ 定时任务:"
	@echo "  make scheduler      - 运行定时任务调度器"
	@echo ""
	@echo "⚙️  本地构建（需要宿主机安装 Go）:"
	@echo "  make build-binary   - 直接在宿主机构建 Go 二进制文件（不使用 Docker）"
	@echo "  make run-binary     - 直接在宿主机构建并运行二进制文件（不使用 Docker）"
	@echo ""
	@echo "🧹 工具命令:"
	@echo "  make clean          - 清理构建产物"
	@echo ""
	@echo "💡 大多数命令会自动检测 Docker/宿主机环境"
	@echo "💡 本地构建命令需要您的机器上安装 Go"

## quick-start: Complete setup and start the project
quick-start:
	@chmod +x scripts/quick-start.sh
	@./scripts/quick-start.sh

## up: Start Docker containers
up:
	@echo "🐳 Starting Docker containers..."
	@docker compose up -d --build --wait
	@echo "✅ Containers started and healthy"
	@echo "📍 API: http://localhost:8080"

## down: Stop Docker containers
down:
	@echo "🛑 Stopping Docker containers..."
	@docker compose down
	@echo "✅ Containers stopped"

## restart: Restart Docker containers
restart:
	@echo "🔄 Restarting Docker containers..."
	@docker compose restart
	@echo "✅ Containers restarted"

## logs: View container logs
logs:
	@docker compose logs -f app

## build: Rebuild Docker containers
build:
	@echo "🔨 Building Docker containers..."
	@docker compose build
	@echo "✅ Build complete"

## test: Run tests
test:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go test ./... -v
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go test ./... -v; \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## test-coverage: Run tests with coverage
test-coverage:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go test ./... -v -coverprofile=coverage.out
	@$(EXEC_CMD) go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go test ./... -v -coverprofile=coverage.out; \
		go tool cover -html=coverage.out -o coverage.html; \
		echo "✅ Coverage report: coverage.html"; \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## lint: Run linter
lint:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@echo "🔍 Running golangci-lint..."
	@$(EXEC_CMD) golangci-lint run --timeout=5m && echo "✅ No linting issues found!" || exit 1
else
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		echo "🔍 Running golangci-lint..."; \
		golangci-lint run --timeout=5m && echo "✅ No linting issues found!" || exit 1; \
	else \
		echo "❌ Error: Docker container not running and golangci-lint not installed"; \
		echo "Please run: make up"; \
		echo "Or install: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
		exit 1; \
	fi
endif

## lint-fix: Run linter and fix issues
lint-fix:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@echo "🔧 Running golangci-lint with auto-fix..."
	@$(EXEC_CMD) golangci-lint run --fix --timeout=5m && echo "✅ Linting complete! Issues auto-fixed where possible." || exit 1
else
	@if command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		echo "🔧 Running golangci-lint with auto-fix..."; \
		golangci-lint run --fix --timeout=5m && echo "✅ Linting complete! Issues auto-fixed where possible." || exit 1; \
	else \
		echo "❌ Error: Docker container not running and golangci-lint not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## swag: Generate Swagger documentation
swag:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) swag init -g ./cmd/server/main.go -o ./api/docs
	@echo "✅ Swagger docs generated"
else
	@if command -v swag >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		swag init -g ./cmd/server/main.go -o ./api/docs; \
		echo "✅ Swagger docs generated"; \
	else \
		echo "❌ Error: Docker container not running and swag not installed"; \
		echo "Please run: make up"; \
		echo "Or install: go install github.com/swaggo/swag/cmd/swag@latest"; \
		exit 1; \
	fi
endif

## migrate-create: Create a new migration
migrate-create:
ifndef NAME
	@echo "❌ Error: NAME is required"
	@echo "Usage: make migrate-create NAME=add_user_avatar"
	@exit 1
endif
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go run cmd/migrate/main.go create $(NAME)
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/migrate/main.go create $(NAME); \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## migrate-up: Apply all pending migrations
migrate-up:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go run cmd/migrate/main.go up
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/migrate/main.go up; \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## migrate-down: Rollback last migration (or N migrations with STEPS=N)
migrate-down:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
ifdef STEPS
	@$(EXEC_CMD_INTERACTIVE) go run cmd/migrate/main.go down $(STEPS)
else
	@$(EXEC_CMD_INTERACTIVE) go run cmd/migrate/main.go down
endif
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
ifdef STEPS
		go run cmd/migrate/main.go down $(STEPS); \
else
		go run cmd/migrate/main.go down; \
endif
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## migrate-status: Show current migration version
migrate-status:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go run cmd/migrate/main.go version
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/migrate/main.go version; \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## migrate-goto: Go to specific version
migrate-goto:
ifndef VERSION
	@echo "❌ Error: VERSION is required"
	@echo "Usage: make migrate-goto VERSION=5"
	@exit 1
endif
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go run cmd/migrate/main.go goto $(VERSION)
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/migrate/main.go goto $(VERSION); \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## migrate-force: Force set version (recovery)
migrate-force:
ifndef VERSION
	@echo "❌ Error: VERSION is required"
	@echo "Usage: make migrate-force VERSION=1"
	@exit 1
endif
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD_INTERACTIVE) go run cmd/migrate/main.go force $(VERSION)
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/migrate/main.go force $(VERSION); \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## migrate-drop: Drop all tables
migrate-drop:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD_INTERACTIVE) go run cmd/migrate/main.go drop --force
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/migrate/main.go drop --force; \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## create-admin: Create new admin user (interactive)
create-admin:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@docker exec -it $(CONTAINER_NAME) go run cmd/createadmin/main.go
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/createadmin/main.go; \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## promote-admin: Promote existing user to admin by ID
promote-admin:
ifndef ID
	@echo "❌ Error: User ID is required"
	@echo "Usage: make promote-admin ID=123"
	@exit 1
endif
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go run cmd/createadmin/main.go --promote=$(ID)
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/createadmin/main.go --promote=$(ID); \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## build-binary: Build Go binary directly on host (requires Go)
build-binary:
	@if ! command -v go >/dev/null 2>&1; then \
		echo "❌ Error: Go is not installed on your machine"; \
		echo ""; \
		echo "Please install Go first:"; \
		echo "  https://golang.org/doc/install"; \
		echo ""; \
		echo "Or use Docker instead:"; \
		echo "  make up"; \
		exit 1; \
	fi
	@echo "🔨 Building Go binary..."
	@mkdir -p bin
	@go build -o bin/server ./cmd/server
	@echo "✅ Binary built successfully: bin/server"
	@echo ""
	@echo "To run the binary:"
	@echo "  make run-binary"
	@echo "  OR"
	@echo "  ./bin/server"

## run-binary: Build and run Go binary directly on host (requires Go)
run-binary: build-binary
	@echo ""
	@echo "🚀 Starting server..."
	@echo ""
	@echo "⚠️  Note: Ensure PostgreSQL is running on localhost:5432"
	@echo "⚠️  Note: Set environment variables or use .env file"
	@echo ""
	@./bin/server

## clean: Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -f coverage.out coverage.html
	@rm -f bin/*
	@docker compose down -v 2>/dev/null || true
	@echo "✅ Clean complete"

## generate-jwt-secret: Generate and set JWT_SECRET in .env if not exists
generate-jwt-secret:
	@if [ ! -f .env ]; then \
		echo "� Creating .env file from .env.example..."; \
		cp .env.example .env 2>/dev/null || touch .env; \
	fi
	@if grep -q "^JWT_SECRET=.\+" .env 2>/dev/null; then \
		echo "✅ JWT_SECRET already exists in .env"; \
		echo "💡 Current value is set (not displayed for security)"; \
		echo ""; \
		echo "To regenerate, remove the current JWT_SECRET line from .env first"; \
	else \
		echo "🔐 Generating JWT secret..."; \
		SECRET=$$(openssl rand -base64 48 | tr -d '\n'); \
		if grep -q "^JWT_SECRET=" .env 2>/dev/null; then \
			sed -i.bak "s|^JWT_SECRET=.*|JWT_SECRET=$$SECRET|" .env && rm -f .env.bak; \
		else \
			echo "JWT_SECRET=$$SECRET" >> .env; \
		fi; \
		echo "✅ JWT_SECRET generated and saved to .env"; \
		echo ""; \
		echo "⚠️  NEVER commit .env to git!"; \
		fi

## setup-db: Configure database options (DEPRECATED - all databases are now included by default)
# setup-db:
# 	@chmod +x scripts/setup-database.sh
# 	@./scripts/setup-database.sh

## scheduler: Run scheduler service
scheduler:
ifdef CONTAINER_RUNNING
	@echo "$(ENV_MSG)"
	@$(EXEC_CMD) go run cmd/scheduler/main.go
else
	@if command -v go >/dev/null 2>&1; then \
		echo "$(ENV_MSG)"; \
		go run cmd/scheduler/main.go; \
	else \
		echo "❌ Error: Docker container not running and Go not installed"; \
		echo "Please run: make up"; \
		exit 1; \
	fi
endif

## check-env: Check if required environment variables are set
check-env:
	@echo "🔍 Checking required environment variables..."
	@if [ -f .env ]; then \
		echo "✅ .env file exists"; \
		if grep -q "^JWT_SECRET=.\+" .env 2>/dev/null; then \
			echo "✅ JWT_SECRET is set in .env"; \
		else \
			echo "❌ JWT_SECRET is missing or empty in .env"; \
			echo "   Run: make generate-jwt-secret"; \
			exit 1; \
		fi \
	else \
		echo "❌ .env file not found"; \
		echo "   Copy .env.example to .env and set JWT_SECRET"; \
		exit 1; \
	fi
