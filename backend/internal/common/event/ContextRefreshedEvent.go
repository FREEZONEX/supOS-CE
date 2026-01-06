package event

import "backend/internal/svc"

// ContextRefreshedEvent 类似java spring ContextRefreshedEvent 事件，表示系统初始化完成
type ContextRefreshedEvent struct {
	SvcContext *svc.ServiceContext
}
