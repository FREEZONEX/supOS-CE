package service

import (
	"backend/internal/common/I18nUtils"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
)

func (l *UnsQueryService) DetectIfFieldReferenced(ctx context.Context, req *types.UpdateModeRequestVo) (resp *types.ResultVO, err error) {
	dataMap := map[string]any{"referred": false}
	resp = &types.ResultVO{Code: 200, Msg: "ok", Data: dataMap}

	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.GetByAlias(db, req.Alias)
	if er != nil || uns == nil {
		err = er
		return
	}
	if uns.PathType == 0 {
		return
	}
	files, _ := l.unsMapper.ListFileByTemplateId(db, uns.Id)
	if len(files) == 0 {
		return
	}
	dataMap["referred"] = true
	dataMap["tips"] = I18nUtils.GetMessage("uns.update.field.tips1")
	return
}
