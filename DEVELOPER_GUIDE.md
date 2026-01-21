# UYou Go API Starter 开发者指南

**版本**: v2.0.0  
**最后更新**: 2026-01-20


## 目录

- [1. 架构与目录结构](#1-架构与目录结构)
- [2. 5 分钟快速开始](#2-5-分钟快速开始)
- [3. 如何确认各服务是否启用/可用](#3-如何确认各服务是否启用可用)
- [4. 核心功能最小上手（按功能）](#4-核心功能最小上手按功能)
- [5. 新模块开发（Handler → Service → Repository）](#5-新模块开发handler--service--repository)
- [6. 测试与质量流程](#6-测试与质量流程)
- [7. 微服务落地建议（拆分、通信、部署）](#7-微服务落地建议拆分通信部署)
- [8. 常用命令速查](#8-常用命令速查)
- [9. 常见问题](#9-常见问题)

---

## 1. 架构与目录结构

### 1.1 清晰架构（Clean Architecture）

项目遵循三层架构，关注点分离：

| 层次 | 职责 | 对应文件 |
|---|---|---|
| **Handler** | 处理 HTTP/gRPC 请求、校验输入、返回响应 | `internal/<domain>/handler.go` |
| **Service** | 业务逻辑（不依赖 HTTP） | `internal/<domain>/service.go` |
| **Repository** | 数据访问与持久化（只关心 DB） | `internal/<domain>/repository.go` |

依赖方向：Handler → Service → Repository（内层不依赖外层）。

### 1.2 目录结构

```
.
├── cmd/                 # 程序入口（server / scheduler / migrate / createadmin）
├── internal/            # 业务与基础设施（auth/user/db/redis/mongodb/grpc/messaging/...）
├── configs/             # 配置文件（YAML）+ 环境变量覆盖
├── migrations/          # SQL 迁移（golang-migrate）
├── api/                 # Swagger / Protobuf
├── tests/               # 集成/端到端测试
├── docker-compose.yml   # 开发环境依赖（db/mongodb/redis/rabbitmq/app）
└── Makefile             # 常用命令
```

---

## 2. 5 分钟快速开始

### 2.1 启动（Docker 推荐）

```bash
# 1) 生成 JWT Secret（必需）
make generate-jwt-secret

# 2) 一键启动（构建并启动容器）
make quick-start

# 3) 看日志
make logs

# 4) 健康检查
curl http://localhost:8080/health
```

### 2.2 常用访问地址

- **API**: `http://localhost:8080/api/v1`
- **Swagger**: `http://localhost:8080/swagger/index.html`
- **健康检查**: `http://localhost:8080/health`
- **RabbitMQ 管理台**: `http://localhost:15672`（默认 `guest/guest`）

> 监控 `/metrics` 的暴露方式见下方第 4.6 节（Prometheus）。

---

## 3. 如何确认各服务是否启用/可用

### 3.1 通过健康检查（推荐）

```bash
curl http://localhost:8080/health | jq
```

查看响应中的 `components` 字段（例如 `database/redis/mongodb`）来判断是否连通与健康。

### 3.2 通过配置开关（启用/禁用）

配置文件：`configs/config.yaml`  
任意项都可被环境变量覆盖（优先级：环境变量 > 环境配置文件 > `config.yaml`）。

常用开关：

```bash
export MONGODB_ENABLED=false
export REDIS_ENABLED=false
export RABBITMQ_ENABLED=false
export GRPC_ENABLED=false
export METRICS_ENABLED=false
make restart
```

### 3.4 配置加载优先级与“哪些 env 会生效”

#### 配置加载优先级（从高到低）

1. **环境变量**（在 Docker 下通常来自 `env_file: .env` 注入）
2. **环境配置文件**：`configs/config.<APP_ENVIRONMENT>.yaml`
3. **基础配置文件**：`configs/config.yaml`

也就是说：同一个配置项如果在 env 里设置了值，会覆盖 `config.<env>.yaml` 和 `config.yaml`。

#### 重要：只有绑定过的 env 才“稳定覆盖”

项目通过 `internal/config/config.go` 的 `bindEnvVariables()` 显式绑定一组环境变量（`v.BindEnv(key, ENV_NAME)`）。
**只有这些绑定过的 env，才能保证一定覆盖 YAML。**


### 3.3 手动连通性验证（排障用）

```bash
# PostgreSQL
make exec-db
psql -U postgres -d uyou_api -c "SELECT version();"

# Redis
make exec
redis-cli -h redis ping

# MongoDB
make exec
mongosh mongodb://mongodb:27017 --eval "db.adminCommand('ping')"
```

---

## 4. 核心功能最小上手（按功能）

本节只保留“最小可用步骤”，避免重复与长篇代码，完整实现参考现有模块：`internal/user/`。

### 4.1 PostgreSQL（迁移 + CRUD）

```bash
make migrate-create NAME=create_posts_table
make migrate-up
make migrate-status
```

### 4.2 Redis（缓存 + 锁）

代码示例参考：
- `internal/redis/cache.go`
- `internal/user/cache_service.go`

### 4.3 MongoDB（基础仓库）

代码示例参考：
- `internal/mongodb/mongodb.go`
- `internal/mongodb/repository.go`

### 4.4 RabbitMQ（事件发布订阅）

消息队列已抽象为接口，便于后续切换云厂商实现：
- 接口：`internal/messaging/queue.go`
- RabbitMQ 实现：`internal/messaging/rabbitmq.go`
- EventBus：`internal/messaging/events.go`
- 工厂（可选）：`internal/messaging/factory.go`

### 4.5 gRPC（服务间通信）

参考：
- 服务端：`internal/grpc/server/`
- 客户端：`internal/grpc/client/`

### 4.6 Prometheus（指标）

当前仓库已集成指标“采集与记录”（`internal/metrics/` + `internal/middleware/metrics.go`），但**默认未提供暴露 `/metrics` 的 HTTP 端点**。

启用方式二选一：

- **方式 A：在主 HTTP 服务器上暴露 `/metrics`**
  - 在 `internal/server/router.go` 注册 `router.GET("/metrics", promhttp.Handler())`（按你们的路径配置 `cfg.Metrics.Path`）。
- **方式 B：独立 metrics 服务器**
  - 在 `cmd/server/main.go` 启动第二个 `http.Server`，监听 `cfg.Metrics.Port`。

---

## 5. 新模块开发（Handler → Service → Repository）

以新增 `post` 模块为例（步骤清单）：

1. 创建目录：`mkdir -p internal/post`
2. 定义模型：`internal/post/model.go`
3. 定义 DTO：`internal/post/dto.go`
4. Repository：`internal/post/repository.go`
5. Service：`internal/post/service.go`
6. Handler：`internal/post/handler.go`（包含 Swagger 注解）
7. 创建迁移：`make migrate-create NAME=create_posts_table`（如需 SQL）
8. 注册路由：`internal/server/router.go`
9. 写测试：`internal/post/*_test.go` + 必要时 `tests/*`
10. 本地验证：`make migrate-up && make test && make lint && make swag`

模块模板请直接参考 `internal/user/`（最完整示例）。

---

## 6. 测试与质量流程

### 6.1 测试

```bash
make test
make test-coverage
```

建议：
- 单元测试：贴近模块（`internal/<domain>/*_test.go`）
- 集成测试：放在 `tests/`

### 6.2 提交前流程（推荐）

```bash
make lint-fix
make lint
make test
make swag
```

---

## 7. 微服务落地建议（拆分、通信、部署）

### 7.1 拆分建议

按业务边界拆服务（示例）：
- `centralService`（网关/聚合/管理）
- `chatService`
- `feedService`
- `socketService`

### 7.2 通信策略

- **同步**：HTTP/gRPC（低延迟、强一致读写）
- **异步**：消息队列（削峰填谷、解耦、最终一致）

### 7.3 依赖服务选择

- 开发：`docker-compose.yml` 一套依赖即可
- 生产：优先用托管服务（RDS/ElastiCache/DocumentDB 或 Atlas/云消息队列），降低运维成本

---

## 8. 常用命令速查

```bash
# 启停/日志
make up
make down
make restart
make logs

# 测试/质量
make test
make test-coverage
make lint
make lint-fix

# Swagger
make swag

# 迁移
make migrate-create NAME=create_xxx_table
make migrate-up
make migrate-down
make migrate-status

# 管理员
make create-admin
make promote-admin ID=1

# 定时任务
make scheduler
```

---

## 9. 常见问题

### 9.1 如何添加新的环境变量？

1. 在 `internal/config/config.go` 增加字段
2. 在 `configs/config.yaml` 增加默认值
3. 在 `.env.example` 增加示例

### 9.2 如何添加新的 gRPC 服务？

1. 在 `api/proto/` 新增 `.proto`
2. 生成代码（如仓库有对应命令则用 `make proto`，否则按现有 proto 生成方式执行）
3. 在 `internal/grpc/server/` 实现服务并注册
