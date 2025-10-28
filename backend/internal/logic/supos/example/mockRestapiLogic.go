package example

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MockRestapiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 模拟restApi数据
func NewMockRestapiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockRestapiLogic {
	return &MockRestapiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MockRestapiLogic) MockRestapi() error {
	// todo: add your logic here and delete this line

	return nil
}
