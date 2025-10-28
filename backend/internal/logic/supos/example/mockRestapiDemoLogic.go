package example

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MockRestapiDemoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 模拟电器厂IT数据
func NewMockRestapiDemoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockRestapiDemoLogic {
	return &MockRestapiDemoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MockRestapiDemoLogic) MockRestapiDemo() error {
	// todo: add your logic here and delete this line

	return nil
}
