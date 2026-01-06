package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/dto"
	"backend/internal/common/utils/JsonUtil"
	"backend/internal/common/utils/PathUtil"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"strconv"
	"strings"
)

func (l *UnsTemplateService) PageList(ctx context.Context, req *types.TemplateQueryVo) (resp *types.TemplatePageResp, err error) {
	db := dao.GetDb(ctx)
	dbQuery := dto.TemplateQueryVo{Key: req.Key, SubscribeEnable: req.SubscribeEnable}
	total := int64(0)
	rs, er := l.unsMapper.PageListTemplates(db, dbQuery, req.PageNo, req.PageSize, &total)
	if er != nil {
		err = er
		return
	}
	resp = &types.TemplatePageResp{Data: base.Map[*dao.UnsNamespace, *types.TemplateSearchResult](rs, func(e *dao.UnsNamespace) *types.TemplateSearchResult {
		return &types.TemplateSearchResult{ID: strconv.FormatInt(e.Id, 10), Name: PathUtil.GetName(e.Path), Description: e.Description, Alias: e.Alias}
	})}
	resp.Code, resp.Msg = 200, "OK"
	resp.PageNo, resp.PageSize, resp.Total = req.PageNo, req.PageSize, total
	return
}
func (l *UnsTemplateService) PageListUnsByTemplate(ctx context.Context, req *types.PageListUnsByTemplateReq) (resp *types.PageListUnsByTemplateResp, err error) {
	db := dao.GetDb(ctx)
	total := int64(0)
	rs, er := l.unsMapper.PageListByTemplateId(db, req.TemplateId, req.PageNo, req.PageSize, &total)
	if er != nil {
		err = er
		return
	}
	resp = &types.PageListUnsByTemplateResp{Data: base.Map[*dao.UnsNamespace, *types.FileVo](rs, func(e *dao.UnsNamespace) *types.FileVo {
		return &types.FileVo{UnsId: strconv.FormatInt(e.Id, 10), PathType: e.PathType, Path: e.Path, Name: PathUtil.GetName(e.Path)}
	})}
	resp.Code, resp.Msg = 200, "OK"
	resp.PageNo, resp.PageSize, resp.Total = req.PageNo, req.PageSize, total
	return
}
func (l *UnsTemplateService) DetailById(ctx context.Context, req *types.WithID) (resp *types.TemplateDetailResp, err error) {
	db := dao.GetDb(ctx)
	var po *dao.UnsNamespace
	po, err = l.unsMapper.SelectById(db, req.ID)
	resp = &types.TemplateDetailResp{}
	if po == nil {
		resp.Code, resp.Msg = 404, I18nUtils.GetMessage("uns.template.not.exists")
		return
	}
	resp.Code, resp.Msg = 200, "OK"
	resp.Data, err = l.uns2TemplateVo(po)
	return
}

func (l *UnsTemplateService) DetailByAlias(ctx context.Context, req *types.WithAlias) (resp *types.TemplateDetailResp, err error) {
	db := dao.GetDb(ctx)
	var po *dao.UnsNamespace
	po, err = l.unsMapper.GetByAlias(db, req.Alias)
	resp = &types.TemplateDetailResp{}
	if po == nil {
		resp.Code, resp.Msg = 404, I18nUtils.GetMessage("uns.template.not.exists")
		return
	}
	resp.Code, resp.Msg = 200, "OK"
	resp.Data, err = l.uns2TemplateVo(po)
	return
}

func (l *UnsTemplateService) uns2TemplateVo(po *dao.UnsNamespace) (*types.TemplateVo, error) {
	name := PathUtil.GetName(po.Path)
	vo := &types.TemplateVo{
		Topic:       "template/" + name,
		ID:          strconv.FormatInt(po.Id, 10),
		Name:        name,
		Alias:       po.Alias,
		Fields:      po.Fields,
		CreateTime:  po.CreateAt.UnixMilli(),
		Description: base.P2v(po.Description),
	}
	if flags := base.P2v(po.WithFlags); flags > 0 {
		enable := constants.WithSubscribeEnable(flags)
		vo.SubscribeEnable = &enable
		if enable {
			protocol := base.P2v(po.Protocol)
			if strings.HasPrefix(protocol, "{") {
				var pMap = make(map[string]string)
				JsonUtil.FromJson(protocol, &pMap)
				if freq, has := pMap["frequency"]; has {
					vo.SubscribeFrequency = freq
				}
			}
		}
	}
	return vo, nil
}
