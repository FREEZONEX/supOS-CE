package service

import (
	"backend/internal/common/I18nUtils"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
)

func (l *UnsTemplateService) Delete(ctx context.Context, req *types.WithID) (rs *types.BaseResult, err error) {
	templateId := req.ID
	db := dao.GetDb(ctx)
	template, _ := l.unsMapper.SelectById(db, templateId)
	if template == nil {
		return &types.BaseResult{Code: 404, Msg: I18nUtils.GetMessage("uns.template.not.exists")}, nil
	}
	True := true
	options := types.RemoveUnsOptions{
		WithFlow:      &True,
		WithDashboard: &True,
		RemoveRefer:   &True,
	}
	delRs, delErr := l.unsRemoveService.Remove(ctx, options, []*dao.UnsNamespace{template})
	if delErr != nil {
		err = delErr
	} else if delRs != nil {
		rs = &types.BaseResult{Code: delRs.Code, Msg: delRs.Msg}
	}
	return
}
