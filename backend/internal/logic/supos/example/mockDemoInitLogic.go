package example

import (
	"context"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type MockDemoInitLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 初始化发电机数据
func NewMockDemoInitLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MockDemoInitLogic {
	return &MockDemoInitLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MockDemoInitLogic) MockDemoInit() error {
	// todo: add your logic here and delete this line

	return nil
}
