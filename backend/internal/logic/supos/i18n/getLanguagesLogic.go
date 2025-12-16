package i18n

import (
	"context"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLanguagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 返回支持的语言
func NewGetLanguagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLanguagesLogic {
	return &GetLanguagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLanguagesLogic) GetLanguages(req *types.GetLanguagesReq) (resp *types.GetLanguagesResp, err error) {
	resp = &types.GetLanguagesResp{}
	resp.List = make([]types.I18nLanguageVO, 0, 5)
	var i int64 = 1
	for it := range l.svcCtx.I18n {
		i++
		vo := types.I18nLanguageVO{
			Id:           i,
			LanguageType: 1,
			HasUsed:      true,
		}
		vo.LanguageCode = it.String()
		lang, ok := GetLangInfoByCode(vo.LanguageCode)
		if ok {
			vo.LanguageName = lang.NativeName
		}
		resp.List = append(resp.List, vo)
	}
	return resp, nil
}
