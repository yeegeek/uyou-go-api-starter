// Package health 定义健康检查器接口
package health

import "context"

type Checker interface {
	Name() string
	Check(ctx context.Context) CheckResult
}
