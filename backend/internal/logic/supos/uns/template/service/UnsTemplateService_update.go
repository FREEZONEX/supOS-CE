package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"strings"
)

func (l *UnsTemplateService) UpdateBaseInfo(ctx context.Context, req *types.UpdateTemplateBaseInfoReq) (rs *types.BaseResult, err error) {
	rs = &types.BaseResult{Code: 400, Msg: ""}
	db := dao.GetDb(ctx)
	uns, er := l.unsMapper.SelectById(db, req.ID)
	if er != nil || uns == nil {
		rs.Msg = I18nUtils.GetMessageWithCtx(ctx, "uns.template.not.exists")
		return
	}
	unsDto := buildTemplateBaseInfoUpdateDto(uns, req)
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
		rs.Msg = I18nUtils.GetMessageWithCtx(ctx, "uns.template.not.exists")
		return
	}
	unsDto := buildTemplateFieldsAndDescUpdateDto(uns, req)
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
		rs.Msg = I18nUtils.GetMessageWithCtx(ctx, "uns.template.not.exists")
		return
	}
	unsDto := buildTemplateSubscribeUpdateDto(uns, req)
	idRs := l.unsAddService.CreateModelInstance(ctx, unsDto)
	if idRs != nil {
		rs.Code, rs.Msg = idRs.Code, idRs.Msg
	}
	return
}

func buildTemplateBaseInfoUpdateDto(uns *dao.UnsNamespace, req *types.UpdateTemplateBaseInfoReq) *types.CreateTopicDto {
	name := ""
	if req != nil {
		name = strings.TrimSpace(req.Name)
	}
	if name == "" && uns != nil {
		name = uns.Name
	}
	dto := &types.CreateTopicDto{
		PathType: constants.PathTypeTemplate,
		Id:       req.ID,
		Alias:    uns.Alias,
		Name:     name,
	}
	if req != nil && req.Description != "NULL" {
		dto.Description = base.V2p(req.Description)
	}
	return dto
}

func buildTemplateFieldsAndDescUpdateDto(uns *dao.UnsNamespace, req *types.UpdateTemplateFieldsAndDescReq) *types.CreateTopicDto {
	alias := ""
	if req != nil {
		alias = strings.TrimSpace(req.Alias)
	}
	if alias == "" && uns != nil {
		alias = uns.Alias
	}
	name := ""
	if uns != nil {
		name = uns.Name
	}
	return &types.CreateTopicDto{
		PathType:    uns.PathType,
		Id:          uns.Id,
		Alias:       alias,
		Name:        name,
		Fields:      req.Fields,
		JsonFields:  req.JsonFields,
		Description: req.ResolveDescription(),
	}
}

func buildTemplateSubscribeUpdateDto(uns *dao.UnsNamespace, req *types.UpdateTemplateSubscribeReq) *types.CreateTopicDto {
	dto := &types.CreateTopicDto{
		PathType:  constants.PathTypeTemplate,
		Id:        req.ID,
		Alias:     uns.Alias,
		Name:      uns.Name,
		Frequency: req.SubscribeFrequency,
	}
	if req.SubscribeEnable != "" {
		dto.SubscribeEnable = base.V2p(req.SubscribeEnable == "true")
	}
	return dto
}
