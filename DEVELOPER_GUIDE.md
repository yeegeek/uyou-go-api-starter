# UYou Go API Starter 开发者指南

**版本**: v2.0.0  
**最后更新**: 2026-01-20  
**作者**: Manus AI

## 1. 引言

本文档是 UYou Go API Starter 项目的**终极开发者指南**，旨在提供一个全面、统一的参考，涵盖项目的设计理念、架构、核心功能、开发流程和最佳实践。它合并了 `CODE_OVERVIEW.md`, `MICROSERVICE_IMPROVEMENTS.md`, `IMPLEMENTATION_SUMMARY.md` 和 `SCHEDULER_GUIDE.md` 的核心内容，并进行了精简和优化。

## 2. 架构设计：清晰架构 (Clean Architecture)

项目遵循**清晰架构**原则，实现关注点分离，核心是三层架构：

| 层次       | 职责                               | 对应目录/模块                     |
| :--------- | :--------------------------------- | :-------------------------------- |
| **处理器 (Handler)** | 负责处理 HTTP/gRPC 请求和响应，解析输入，调用服务层 | `internal/<domain>/handler.go`        |
| **服务 (Service)**   | 包含核心业务逻辑，不关心数据如何传输或存储   | `internal/<domain>/service.go`        |
| **仓库 (Repository)**| 负责数据持久化，封装与数据库的交互         | `internal/<domain>/repository.go`     |

**依赖规则**：外层依赖内层（Handler → Service → Repository），内层永远不能知道外层的任何信息。

## 3. 目录结构

项目遵循 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)。

```
.
├── cmd/                    # 应用程序入口 (server, scheduler, migrate)
├── internal/              # 内部应用代码 (不对外暴露)
│   ├── auth/             # 认证与授权
│   ├── user/             # 用户管理模块 (示例)
│   ├── middleware/       # Gin 中间件
│   ├── scheduler/        # 定时任务
│   ├── redis/            # Redis 客户端
│   ├── mongodb/          # MongoDB 客户端
│   ├── grpc/             # gRPC 服务
│   ├── messaging/        # RabbitMQ 消息队列
│   └── metrics/          # Prometheus 监控
├── migrations/            # SQL 数据库迁移文件
├── configs/              # 配置文件 (YAML)
├── api/                  # API 定义 (Protobuf, Swagger)
├── scripts/              # 辅助脚本
├── tests/                # 集成测试和端到端测试
├── docker-compose.yml    # 开发环境 Docker Compose 配置
├── Makefile              # 常用命令集合
└── go.mod                # Go 模块依赖
```

## 4. 核心功能与实现

### 4.1. 基础设施

| 功能 | 实现库 | 配置文件 | 使用示例 |
|:---|:---|:---|:---|
| **PostgreSQL** | `gorm.io/gorm` | `database` | `internal/db/db.go` |
| **MongoDB** | `go.mongodb.org/mongo-driver` | `mongodb` | `internal/mongodb/mongodb.go` |
| **Redis** | `github.com/redis/go-redis/v9` | `redis` | `internal/redis/redis.go` |
| **RabbitMQ** | `github.com/rabbitmq/amqp091-go` | `rabbitmq` | `internal/messaging/rabbitmq.go` |
| **gRPC** | `google.golang.org/grpc` | `grpc` | `internal/grpc/server/server.go` |
| **Prometheus** | `github.com/prometheus/client_golang` | `metrics` | `internal/metrics/metrics.go` |
| **定时任务** | `github.com/robfig/cron/v3` | `scheduler` | `internal/scheduler/scheduler.go` |

### 4.2. 认证与授权

- **JWT 认证**：通过 `AuthMiddleware` 实现，支持访问令牌和刷新令牌。
- **RBAC 权限控制**：通过 `RequireAdmin()` 和 `RequireRole()` 中间件实现。
- **Context 辅助函数**：使用 `contextutil.GetUserID(c)` 等函数安全地获取用户信息。

### 4.3. 错误处理

- **统一错误响应**：`internal/errors` 模块定义了标准错误格式。
- **错误处理中间件**：`ErrorHandler` 捕获错误并返回统一的 JSON 响应。

### 4.4. 配置管理

- **Viper**：支持 YAML 文件和环境变量覆盖。
- **数据库选择**：`make setup-db` 交互式配置数据库。

## 5. 开发流程

### 5.1. 添加新业务模块

1. **创建目录**：`mkdir -p internal/<domain>`
2. **定义模型**：在 `model.go` 中定义 GORM 或 MongoDB 模型。
3. **创建 DTO**：在 `dto.go` 中定义请求和响应的数据结构。
4. **实现 Repository**：在 `repository.go` 中实现数据访问逻辑。
5. **实现 Service**：在 `service.go` 中实现业务逻辑。
6. **创建 Handler**：在 `handler.go` 中实现 API 处理器，并添加 Swagger 注解。
7. **生成迁移**（如果使用 SQL）：`make migrate-create NAME=create_<table>_table`
8. **注册路由**：在 `internal/server/router.go` 中注册新的 API 路由。
9. **编写测试**：为所有层编写单元测试和集成测试。
10. **运行**：`make migrate-up && make test && make lint && make swag`

### 5.2. 添加定时任务

1. **创建任务文件**：在 `internal/scheduler/tasks/` 目录下创建文件，实现 `scheduler.Task` 接口。
2. **注册任务**：在 `cmd/scheduler/main.go` 中注册任务和 cron 表达式。
3. **运行测试**：`make scheduler`

**Cron 表达式格式**（6 字段，支持秒级）：
```
秒 分 时 日 月 周

示例： "0 */1 * * * *"  # 每分钟执行
```

### 5.3. 提交前工作流

```bash
make lint-fix    # 自动修复问题
make lint        # 检查剩余问题
make test        # 运行测试
make swag        # 更新 Swagger（如果 API 有变化）
```

## 6. 微服务架构演进

当前框架已具备微服务基础设施，但要构建完整的微服务体系，还需以下改进：

### 6.1. 需要补充的核心功能

| 功能 | 推荐技术 | 目的 |
|:---|:---|:---|
| **服务注册与发现** | Consul / etcd | 解耦服务物理地址，实现动态扩缩容 |
| **分布式配置中心** | Consul KV / Nacos | 集中管理配置，支持动态更新 |
| **分布式追踪** | OpenTelemetry + Jaeger | 跨服务追踪请求链路，快速定位问题 |
| **API 网关** | KrakenD / Tyk | 统一入口，负责路由、认证、限流等 |

### 6.2. 需要改进的现有功能

- **弹性与容错**：引入断路器（`gobreaker`）和重试机制。
- **分布式事务**：引入 Saga 模式或 TCC 模式。
- **容器化与部署**：提供 Helm Chart 用于 Kubernetes 部署。

### 6.3. 改进路线图

1.  **第一阶段：基础增强**：引入 OpenTelemetry 和 Prometheus。
2.  **第二阶段：服务化改造**：引入 gRPC 和 Consul。
3.  **第三阶段：生产级加固**：引入断路器、重试机制和 Helm Chart。

## 7. 快速参考

| 任务 | 命令 |
|:---|:---|
| 启动开发环境 | `make up` |
| 运行测试 | `make test` |
| 代码检查 | `make lint` |
| 创建迁移 | `make migrate-create NAME=<name>` |
| 应用迁移 | `make migrate-up` |
| 更新 Swagger | `make swag` |
| 查看日志 | `make logs` |
| 进入应用容器 | `make exec` |
| 运行定时任务 | `make scheduler` |
| 配置数据库 | `make setup-db` |

---


**仓库地址**：https://github.com/yeegeek/uyou-go-api-starter
