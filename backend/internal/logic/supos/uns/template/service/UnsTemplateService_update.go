package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
)

func (l *UnsTemplateService) UpdateBaseInfo(ctx context.Context, req *types.UpdateTemplateBaseInfoReq) (rs *types.BaseResult, err error) {
	rs = &types.BaseResult{Code: 400, Msg: ""}
	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.SelectById(db, req.ID)
	if er != nil || uns == nil {
		rs.Msg = I18nUtils.GetMessage("uns.template.not.exists")
		return
	}
	unsDto := &types.CreateTopicDto{
		PathType: constants.PathTypeTemplate,
		Id:       req.ID,
		Name:     req.Name,
	}
	if req.Description != "NULL" {
		unsDto.Description = base.V2p(req.Description)
	}
	idRs := l.unsAddService.CreateModelInstance(ctx, unsDto)
	if idRs != nil {
		rs.Code, rs.Msg = idRs.Code, idRs.Msg
	}
	return
}
func (l *UnsTemplateService) UpdateFieldsAndDesc(ctx context.Context, req *types.UpdateTemplateFieldsAndDescReq) (rs *types.BaseResult, err error) {
	rs = &types.BaseResult{Code: 400, Msg: ""}
	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.GetByAlias(db, req.Alias)
	if er != nil || uns == nil {
		rs.Msg = I18nUtils.GetMessage("uns.template.not.exists")
		return
	}
	unsDto := &types.CreateTopicDto{
		PathType:    uns.PathType,
		Id:          uns.Id,
		Alias:       req.Alias,
		Fields:      req.Fields,
		JsonFields:  req.JsonFields,
		Description: req.Description,
	}
	idRs := l.unsAddService.CreateModelInstance(ctx, unsDto)
	if idRs != nil {
		rs.Code, rs.Msg = idRs.Code, idRs.Msg
	}
	return
}
func (l *UnsTemplateService) UpdateSubscribe(ctx context.Context, req *types.UpdateTemplateSubscribeReq) (rs *types.BaseResult, err error) {
	rs = &types.BaseResult{Code: 400, Msg: ""}
	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.SelectById(db, req.ID)
	if er != nil || uns == nil {
		rs.Msg = I18nUtils.GetMessage("uns.template.not.exists")
		return
	}
	unsDto := &types.CreateTopicDto{
		PathType:  constants.PathTypeTemplate,
		Id:        req.ID,
		Frequency: req.SubscribeFrequency,
	}
	if req.SubscribeEnable != "" {
		unsDto.SubscribeEnable = base.V2p(req.SubscribeEnable == "true")
	}
	idRs := l.unsAddService.CreateModelInstance(ctx, unsDto)
	if idRs != nil {
		rs.Code, rs.Msg = idRs.Code, idRs.Msg
	}
	return
}
