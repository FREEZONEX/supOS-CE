package example

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MockRestapiOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 模拟订单数据
func NewMockRestapiOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockRestapiOrderLogic {
	return &MockRestapiOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MockRestapiOrderLogic) MockRestapiOrder() error {
	// todo: add your logic here and delete this line

	return nil
}
