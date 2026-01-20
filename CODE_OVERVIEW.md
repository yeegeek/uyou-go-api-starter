# UYou Go API Starter 代码结构与设计文档

**作者**: Manus AI
**日期**: 2026年1月20日

## 1. 引言

本文档旨在提供对 UYou Go API Starter 项目的全面理解，包括其架构设计、目录结构、核心模块功能以及请求处理流程。通过本文档，开发者可以快速掌握代码库的核心思想，并在此基础上进行二次开发。

## 2. 架构设计：清晰架构 (Clean Architecture)

本项目遵循**清晰架构**（Clean Architecture）原则，旨在实现关注点分离、降低耦合度、提高可测试性和可维护性。核心思想是将软件分为多个层次，并强制执行依赖关系规则：**内层永远不能知道外层的任何信息**。

本项目主要分为三层：

| 层次       | 职责                               | 对应目录/模块                     |
| :--------- | :--------------------------------- | :-------------------------------- |
| **处理器 (Handler)** | 负责处理 HTTP 请求和响应，解析输入，调用服务层 | `internal/user/handler.go`        |
| **服务 (Service)**   | 包含核心业务逻辑，不关心数据如何传输或存储   | `internal/user/service.go`        |
| **仓库 (Repository)**| 负责数据持久化，封装与数据库的交互         | `internal/user/repository.go`     |

这种分层确保了业务逻辑（Service）与具体实现（如 Web 框架 Gin 或数据库 ORM GORM）的解耦。

## 3. 目录结构详解

项目遵循 Go 社区推荐的 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)。

```
.
├── cmd/                    # 应用程序入口
│   ├── server/            # API 服务器主程序
│   ├── migrate/           # 数据库迁移工具
│   └── createadmin/       # 管理员创建工具
├── internal/              # 内部应用代码（不对外暴露）
│   ├── auth/             # 认证与授权模块
│   ├── user/             # 用户管理模块
│   ├── health/           # 健康检查模块
│   ├── middleware/       # Gin 中间件
│   ├── errors/           # 统一错误处理
│   ├── config/           # 配置加载与管理
│   ├── db/               # 数据库连接
│   ├── server/           # HTTP 服务器与路由
│   └── contextutil/      # 上下文工具函数
├── migrations/            # SQL 数据库迁移文件
├── configs/              # 配置文件 (YAML)
├── api/                  # API 文档 (Swagger, Postman)
├── scripts/              # 辅助脚本 (如：启动脚本)
├── tests/                # 集成测试和端到端测试
├── docker-compose.yml    # 开发环境 Docker Compose 配置
├── docker-compose.prod.yml # 生产环境 Docker Compose 配置
├── Dockerfile            # 多阶段 Docker 镜像构建文件
├── Makefile              # 常用命令集合
├── go.mod                # Go 模块依赖
└── README.md             # 项目说明文档
```

## 4. 核心模块分析

所有核心业务代码都位于 `internal/` 目录下，确保了项目内部逻辑的封装性。

### 4.1. `internal/config` - 配置管理

- **职责**: 使用 [Viper](https://github.com/spf13/viper) 库加载和管理配置。
- **特性**:
    - 支持多文件配置（`config.yaml`, `config.production.yaml`）。
    - 支持环境变量覆盖，实现云原生友好。
    - 配置优先级：**环境变量 > 环境特定配置 > 基础配置**。
    - 提供强类型的配置结构体，如 `Config`, `DatabaseConfig` 等。

### 4.2. `internal/db` - 数据库连接

- **职责**: 初始化并管理数据库连接。
- **实现**: 使用 [GORM](https://gorm.io/) 作为 ORM 框架，支持 PostgreSQL 和 SQLite（用于测试）。

### 4.3. `internal/server` - HTTP 服务器

- **职责**: 设置 Gin 引擎，注册路由和中间件。
- **`router.go`**: 核心文件，定义了所有 API 端点、应用的中间件栈（日志、恢复、CORS、错误处理等）以及路由分组（如 `/api/v1`）。

### 4.4. `internal/middleware` - 中间件

提供可重用的 Gin 中间件：

- **`logger.go`**: 结构化日志中间件，记录请求详情。
- **`rbac.go`**: 基于角色的访问控制（RBAC）中间件，如 `RequireAdmin()`。
- **`rate_limit.go`**: API 速率限制中间件。
- **`pagination.go`**: 从请求中解析分页参数。

### 4.5. `internal/user` - 用户模块

这是实现一个完整业务功能的典型示例：

- **`model.go`**: 定义 `User` 和 `Role` 的 GORM 模型。
- **`dto.go`**: 定义数据传输对象（DTO），如 `RegisterRequest`, `UserResponse`，用于请求绑定和响应序列化。
- **`repository.go`**: `Repository` 接口和实现，封装了对用户数据的增删改查操作。
- **`service.go`**: `Service` 接口和实现，处理用户注册、认证、更新等核心业务逻辑。
- **`handler.go`**: `Handler` 结构体，包含 Gin 的处理器函数，负责处理 HTTP 请求，调用 `Service` 并返回响应。

### 4.6. `internal/auth` - 认证模块

- **职责**: 处理所有与认证、授权和 JWT 相关的逻辑。
- **`service.go`**: 核心服务，负责生成/验证访问令牌、生成/刷新/撤销刷新令牌。
- **`middleware.go`**: `AuthMiddleware` 中间件，用于保护需要认证的路由。
- **`refresh_token.go`**: 定义了刷新令牌的模型和仓库，实现了令牌轮换和重用检测的安全策略。

### 4.7. `internal/errors` - 统一错误处理

- **职责**: 提供标准化的 API 错误响应格式。
- **`errors.go`**: 定义了 `APIError` 结构体和常用的错误构造函数（如 `NotFound`, `BadRequest`）。
- **`middleware.go`**: `ErrorHandler` 中间件，捕获在上下文中记录的错误，并将其转换为统一的 JSON 错误响应，避免向客户端泄露内部错误细节。

## 5. 请求生命周期

一个典型的 API 请求（例如 `POST /api/v1/auth/register`）会经过以下流程：

1.  **入口**: 请求到达 `cmd/server/main.go`，启动 Gin 服务器。
2.  **路由匹配**: Gin 引擎在 `internal/server/router.go` 中找到匹配的路由，执行相应的中间件栈。
3.  **中间件**: 请求依次通过日志、CORS、速率限制等中间件。
4.  **处理器 (Handler)**: 请求到达 `internal/user/handler.go` 中的 `Register` 方法。
    -   处理器使用 `c.ShouldBindJSON()` 将请求体绑定到 `RegisterRequest` DTO 并进行验证。
5.  **服务 (Service)**: 处理器调用 `internal/user/service.go` 中的 `RegisterUser` 方法，并传入 DTO。
    -   服务层执行业务逻辑：检查邮箱是否已存在、哈希密码等。
6.  **仓库 (Repository)**: 服务层调用 `internal/user/repository.go` 中的 `Create` 方法，将用户数据持久化到数据库。
7.  **返回响应**: 
    -   仓库层将创建的用户模型返回给服务层。
    -   服务层将用户模型返回给处理器。
    -   处理器将用户模型转换为 `UserResponse` DTO（隐藏敏感信息），并以 JSON 格式返回给客户端。

## 6. 数据库迁移

项目使用 `golang-migrate` 工具管理数据库 schema 的演进。

- **迁移文件**: 位于 `migrations/` 目录，每个迁移都包含一个 `.up.sql` 和一个 `.down.sql` 文件。
- **执行迁移**: 通过 `Makefile` 中的命令执行，如 `make migrate-up`。
- **工具入口**: `cmd/migrate/main.go` 是迁移工具的命令行入口，负责解析命令并调用 `internal/migrate` 中的逻辑。

## 7. 总结

UYou Go API Starter 提供了一个结构清晰、功能完备、安全可靠的起点。其核心优势在于良好的分层架构和对最佳实践的遵循。开发者应重点理解 `internal` 目录下的模块划分和它们之间的依赖关系，这将有助于高效地进行后续开发和维护。
