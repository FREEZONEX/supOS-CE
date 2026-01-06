package service

import (
	"backend/internal/common/I18nUtils"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsRemoveService struct {
	log       logx.Logger
	unsMapper dao.UnsNamespaceRepo
}

func init() {
	spring.RegisterBean(&UnsRemoveService{
		log: logx.WithContext(context.Background()),
	})
}

var TRUE = true
var FALSE = false

func defTrue(b *bool) *bool {
	if b == nil || *b {
		return &TRUE
	}
	return &FALSE
}
func defFalse(b *bool) *bool {
	if b == nil || !*b {
		return &FALSE
	}
	return &TRUE
}
func (r *UnsRemoveService) RemoveModelOrInstance(ctx context.Context, req *types.RemoveReq) (resp *types.RemoveResult, err error) {

	return r.removeModelOrInstance(ctx, req.Id, types.BatchRemoveUnsDto{
		RemoveUnsOptions: types.RemoveUnsOptions{
			WithFlow:        defTrue(req.WithFlow),
			WithDashboard:   defTrue(req.WithDashboard),
			RemoveRefer:     req.RemoveRefer,
			CheckMount:      &TRUE,
			OnlyRemoveChild: &FALSE,
		},
	})
}
func (r *UnsRemoveService) DetectIfRemove(ctx context.Context, req *types.DetectRemoveReq) (resp *types.RemoveResult, err error) {
	resp = &types.RemoveResult{}
	resp.Code = 200
	resp.Msg = "ok"
	db := dao.GetDb(ctx)
	po, er := r.unsMapper.SelectById(db, req.Id)
	if er != nil {
		err = er
		resp.Code = 500
		resp.Msg = "DBErr!"
	} else if po == nil {
		resp.Code = 400
		resp.Msg = I18nUtils.GetMessage("uns.folder.or.file.not.found")
	}
	return
}
func (r *UnsRemoveService) BatchRemoveResultByAliasList(ctx context.Context, req *types.BatchRemoveUnsDto) (resp *types.RemoveResult, err error) {
	req.CheckMount = defTrue(req.CheckMount)
	req.OnlyRemoveChild = defFalse(req.OnlyRemoveChild)
	return r.removeModelOrInstance(ctx, 0, *req)
}
