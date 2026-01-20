# UYou Go API Starter

一个生产就绪的 Go REST API 脚手架项目，采用清晰架构设计，并集成了 gRPC、消息队列和可观测性工具。

## 项目简介

UYou Go API Starter 是一个高质量的 Go 语言 REST API 起始模板，专为快速构建生产级微服务而设计。本项目遵循业界最佳实践，提供完整的认证授权、多数据库支持、gRPC 通信、消息队列、测试框架和可观测性等核心功能。

## 核心特性

### 架构设计
- **清晰架构（Clean Architecture）** - 采用 Handler → Service → Repository 三层架构
- **依赖注入** - 便于测试和维护
- **模块化设计** - 高内聚低耦合

### 安全特性
- **JWT 认证** - 支持访问令牌和刷新令牌机制
- **令牌轮换** - 实现 OAuth 2.0 BCP 最佳实践
- **令牌重用检测** - 自动检测并撤销可疑令牌
- **RBAC 权限控制** - 基于角色的访问控制，支持多对多角色关系
- **密码加密** - 使用 bcrypt 算法进行密码哈希
- **速率限制** - 基于令牌桶算法的 API 限流

### 数据库
- **多数据库支持** - 可选 PostgreSQL、MongoDB 或同时使用
- **PostgreSQL** - 生产级关系型数据库，使用 GORM
- **MongoDB** - 灵活的 NoSQL 数据库，使用官方驱动
- **Redis** - 高性能内存缓存和分布式锁
- **数据库迁移** - 使用 golang-migrate 进行版本化管理
- **连接池** - 优化的数据库连接配置

### 微服务通信
- **gRPC 支持** - 高性能服务间通信
- **消息队列** - 使用 RabbitMQ 实现异步任务和事件驱动架构
- **事件总线** - 统一的事件发布/订阅模型

### 开发体验
- **Docker 支持** - 完整的 Docker 和 Docker Compose 配置
- **热重载** - 使用 Air 实现 2 秒快速重载
- **Swagger 文档** - 自动生成交互式 API 文档
- **Postman 集合** - 预配置的 API 测试集合
- **Make 命令** - 简化的命令行操作

### 可观测性
- **结构化日志** - JSON 格式日志，包含请求 ID 和追踪信息
- **Prometheus 监控** - 预置丰富的监控指标（HTTP、gRPC、数据库、缓存）
- **健康检查** - 支持 Kubernetes 存活探针和就绪探针
- **统一错误处理** - 标准化的错误响应格式
- **优雅关闭** - 零停机部署支持

### 测试
- **单元测试** - 完整的单元测试覆盖
- **集成测试** - 端到端测试支持
- **测试工具** - 使用 testify 断言库

## 快速开始

### ⚠️ 重要：安全配置

**在启动应用之前，必须配置 JWT Secret**：

```bash
# 设置环境变量
export JWT_SECRET="your-32-character-or-longer-secret-key-here"

# 或在 .env 文件中设置
echo 'JWT_SECRET="your-32-character-or-longer-secret-key-here"' >> .env
```

**生成安全的随机密钥**：
```bash
openssl rand -base64 32
```

### 前置要求

- [Docker](https://docs.docker.com/get-docker/) 和 [Docker Compose](https://docs.docker.com/compose/install/)
- [Git](https://git-scm.com/downloads)
- Go 1.24+ (如需本地开发)

### 一键启动

```bash
# 克隆仓库
git clone https://github.com/uyou/uyou-go-api-starter.git
cd uyou-go-api-starter

# 配置数据库并启动服务
make quick-start
```

在执行 `make quick-start` 时，系统会提示您选择要使用的数据库组合。

### 访问服务

启动成功后，可以通过以下地址访问：

- **API 基础地址**: http://localhost:8080/api/v1
- **Swagger 文档**: http://localhost:8080/swagger/index.html
- **健康检查**: http://localhost:8080/health
- **Prometheus 指标**: http://localhost:9091/metrics

### 创建管理员用户

```bash
# 交互式创建管理员
make create-admin

# 将现有用户提升为管理员
make promote-admin ID=1
```

## 项目结构

```
.
├── cmd/                    # 应用程序入口
│   ├── server/            # API 服务器
│   ├── migrate/           # 数据库迁移工具
│   └── createadmin/       # 管理员创建工具
├── internal/              # 内部应用代码
│   ├── auth/             # 认证服务
│   ├── user/             # 用户模块
│   ├── health/           # 健康检查
│   ├── middleware/       # 中间件
│   ├── errors/           # 错误处理
│   ├── config/           # 配置管理
│   ├── db/               # PostgreSQL 数据库连接
│   ├── mongodb/          # MongoDB 数据库连接
│   ├── redis/            # Redis 缓存连接
│   ├── grpc/             # gRPC 服务实现
│   ├── messaging/        # 消息队列和事件总线
│   ├── metrics/          # Prometheus 指标
│   ├── server/           # 服务器路由
│   └── contextutil/      # 上下文工具
├── migrations/            # 数据库迁移文件
├── configs/              # 配置文件
├── api/                  # API 文档、Postman 集合和 Protobuf 定义
│   └── proto/            # Protobuf 文件
├── scripts/              # 脚本文件
├── tests/                # 测试文件
├── docker-compose.yml    # Docker Compose 配置
├── Dockerfile            # Docker 镜像构建
├── Makefile              # Make 命令
└── go.mod                # Go 模块依赖
```

## 常用命令

### 开发命令

```bash
# 启动开发环境（带热重载）
make dev

# 运行测试
make test

# 查看测试覆盖率
make coverage

# 代码格式化
make fmt

# 代码检查
make lint
```

### 数据库迁移

```bash
# 创建新迁移
make migrate-create NAME=add_posts_table

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down

# 查看迁移状态
make migrate-status
```

### Docker 命令

```bash
# 构建镜像
make docker-build

# 启动服务
make docker-up

# 停止服务
make docker-down

# 查看日志
make docker-logs
```

## API 端点

### 认证相关

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/refresh` - 刷新令牌
- `POST /api/v1/auth/logout` - 用户登出

### 用户相关

- `GET /api/v1/auth/me` - 获取当前用户信息
- `PUT /api/v1/auth/me` - 更新当前用户信息
- `GET /api/v1/users/:id` - 获取指定用户信息（需认证）
- `GET /api/v1/users` - 获取用户列表（仅管理员）
- `DELETE /api/v1/users/:id` - 删除用户（仅管理员）

### 健康检查

- `GET /health` - 综合健康检查
- `GET /health/live` - 存活探针
- `GET /health/ready` - 就绪探针

## 配置说明

### 环境变量

项目支持通过环境变量覆盖配置文件中的设置。详见 `.env.example` 文件。

### 配置文件

配置文件位于 `configs/` 目录：

- `config.yaml` - 基础配置
- `config.development.yaml` - 开发环境配置
- `config.staging.yaml` - 预发布环境配置
- `config.production.yaml` - 生产环境配置

配置优先级：环境变量 > 环境特定配置 > 基础配置

## 开发指南

### 添加新模块

1. 在 `internal/` 目录下创建新模块目录
2. 实现以下文件：
   - `model.go` - 数据模型
   - `dto.go` - 数据传输对象
   - `repository.go` - 数据访问层
   - `service.go` - 业务逻辑层
   - `handler.go` - HTTP 处理层
   - `*_test.go` - 测试文件

3. 在 `internal/server/router.go` 中注册路由

### 数据库迁移

创建迁移文件：

```bash
make migrate-create NAME=create_posts_table
```

这将在 `migrations/` 目录下创建两个文件：
- `{timestamp}_create_posts_table.up.sql` - 升级脚本
- `{timestamp}_create_posts_table.down.sql` - 回滚脚本

编写 SQL 后执行迁移：

```bash
make migrate-up
```

### 编写测试

使用 testify 库编写测试：

```go
func TestUserService_RegisterUser(t *testing.T) {
    // 准备测试数据
    req := RegisterRequest{
        Name:     "测试用户",
        Email:    "test@example.com",
        Password: "password123",
    }
    
    // 执行测试
    user, err := service.RegisterUser(ctx, req)
    
    // 断言结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, req.Email, user.Email)
}
```

## 部署

### Docker 部署

```bash
# 构建生产镜像
docker build -t uyou-api:latest .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -e DATABASE_HOST=your-db-host \
  -e DATABASE_PASSWORD=your-password \
  -e JWT_SECRET=your-secret \
  uyou-api:latest
```

### Docker Compose 部署

```bash
# 使用生产配置启动
docker-compose -f docker-compose.prod.yml up -d
```

## 技术栈

- **语言**: Go 1.24+
- **Web 框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL, MongoDB, Redis
- **认证**: JWT (golang-jwt/jwt)
- **配置管理**: Viper
- **数据库迁移**: golang-migrate
- **gRPC**: google.golang.org/grpc
- **消息队列**: RabbitMQ (rabbitmq/amqp091-go)
- **监控**: Prometheus (prometheus/client_golang)
- **测试**: testify
- **文档**: Swagger (swaggo)
- **日志**: slog (Go 标准库)
- **容器**: Docker & Docker Compose

## 许可证

本项目采用 MIT 许可证。详见 [LICENSE](LICENSE) 文件。

## 版权信息

Copyright (c) 2026 UYou Team


## 定时任务

项目集成了 `robfig/cron` 定时任务调度器，可以方便地管理和执行周期性任务。

### 运行调度器

```bash
# 独立运行定时任务调度器
make scheduler
```

### 添加新任务

1. 在 `internal/scheduler/tasks/` 目录下创建新的任务文件，实现 `scheduler.Task` 接口。
2. 在 `cmd/scheduler/main.go` 中注册新任务和对应的 cron 表达式。

### 示例任务

- **Hello World**: 每分钟执行一次，输出日志。
- **清理任务**: 每小时执行一次，用于清理过期数据。
- **统计任务**: 每天凌晨 2 点执行，用于生成统计报表。

## 文档

- [开发者指南](DEVELOPER_GUIDE.md) - 如何开发和贡献代码
- [测试文档](tests/README.md) - 如何运行和编写测试
- [安全审查](docs/security/) - 安全审查报告
- [历史文档](docs/archives/) - 代码审查和优化记录
