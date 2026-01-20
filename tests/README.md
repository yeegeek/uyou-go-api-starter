# 测试指南

本目录包含 API 的集成测试和端到端测试。

## 测试分类

**集成测试**: 验证完整的请求/响应周期，包括处理器、服务和仓库层的测试。

**单元测试**: 应放置在代码旁边：
- `internal/user/service_test.go` - 用户服务测试
- `internal/auth/middleware_test.go` - 认证中间件测试

## 运行测试

```bash
# 运行所有测试
make test

# 或直接使用 go test
go test ./...

# 详细输出模式
go test -v ./...

# 带覆盖率
go test -cover ./...

# 仅运行本目录的测试
go test ./tests/...
```

## 编写新测试

### 1. 创建测试文件
```bash
# 文件名必须以 *_test.go 结尾
touch tests/my_feature_test.go
```

### 2. 基本结构
```go
package tests

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "github.com/gin-gonic/gin"
)

func TestMyFeature(t *testing.T) {
    // 设置
    gin.SetMode(gin.TestMode)
    
    // 创建测试数据库（内存 SQLite）
    db := setupTestDB(t)
    
    // 创建路由
    router := server.SetupRouter(db)
    
    // 发起请求
    req := httptest.NewRequest("GET", "/api/v1/endpoint", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    
    // 断言
    if w.Code != http.StatusOK {
        t.Errorf("期望状态码 200，实际得到 %d", w.Code)
    }
}
```

### 3. 使用现有辅助函数
参考 `handler_test.go` 中的辅助函数：
- `setupTestDB(t)` - 创建内存 SQLite 数据库
- `createTestUser(t, db)` - 创建测试用户
- `getAuthToken(t, db)` - 获取测试用的 JWT 令牌

## 测试数据库

测试使用 **SQLite 内存数据库**，而非 PostgreSQL。这使得测试：
- ✅ 快速
- ✅ 隔离
- ✅ 无外部依赖

## 最佳实践

1. **清理资源**: 每个测试应该独立运行
2. **使用子测试**: 用 `t.Run()` 组织相关测试
3. **测试错误情况**: 不要只测试正常流程
4. **模拟外部服务**: 不要发起真实的 API 调用
5. **使用表驱动测试**: 用于测试多种场景

## 示例：表驱动测试

```go
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"有效邮箱", "user@example.com", false},
        {"无效邮箱", "not-an-email", true},
        {"空邮箱", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("得到错误 %v，期望错误 %v", err, tt.wantErr)
            }
        })
    }
}
```

## 持续集成

测试会在以下情况自动运行：
- 每次推送到 `main` 或 `develop` 分支
- 每个 Pull Request

---

**需要帮助？** 查看 `handler_test.go` 中的现有测试示例。
