package utils

import (
	"context"
	"runtime/debug"

	"github.com/zeromicro/go-zero/core/logx"
)

// Go 启动一个带有 panic recover 的 goroutine，避免单个后台任务崩溃导致整个进程退出。
// 传入的 ctx 仅用于日志和恢复时追踪，不会被用来中断 goroutine 执行。
func Go(ctx context.Context, f func()) {
	go func() {
		defer Recover(ctx)
		f()
	}()
}

// Recover 捕获 panic 并记录错误日志，供 defer 使用。
func Recover(ctx context.Context, infos ...string) {
	if p := recover(); p != nil {
		HandleThrow(ctx, p, infos...)
	}
}

// HandleThrow 统一处理 panic 日志输出。
func HandleThrow(ctx context.Context, p any, msgs ...string) {
	logx.WithContext(ctx).Errorf("panic recovered: msgs=%v error=%v stack=%s", msgs, p, string(debug.Stack()))
}
