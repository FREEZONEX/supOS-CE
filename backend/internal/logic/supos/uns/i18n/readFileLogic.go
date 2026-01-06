package i18n

import (
	"context"
	"sync"

	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/text/language"
)

type ReadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

var once sync.Once

var i18nMap = map[string]string{}

// 获取i18n
func NewReadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReadFileLogic {
	once.Do(func() {

	})
	return &ReadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReadFileLogic) ReadFile(req *types.GetUnsI18nMessagesReq) (resp *types.GetUnsI18nMessagesResp, err error) {
	t, err := language.Parse(req.Lang)
	if err != nil {
		return nil, errors.Parameter.AddMsg("尚未支持").AddDetail(err)
	}
	ms, ok := l.svcCtx.I18n[t]
	if !ok {
		return nil, errors.Parameter.AddMsg("尚未支持")
	}
	var messages = map[string]string{}
	for _, m := range ms.Messages {
		messages[m.ID] = m.Other
	}
	return &types.GetUnsI18nMessagesResp{Messages: messages}, nil
}
