package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/utils/PathUtil"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"context"
)

func (l *UnsTemplateService) Create(ctx context.Context, req *types.CreateTemplateReq) (rs *types.CreateTemplateResp, err error) {
	db := dao.GetDb(ctx)
	rs = &types.CreateTemplateResp{}
	rs.Code, rs.Msg = 200, "OK"
	if req.Alias == "" {
		req.Alias = PathUtil.GenerateAlias(req.Name, 1)
	}
	tmpl, _ := l.unsMapper.GetByAlias(db, req.Alias)
	if tmpl != nil {
		rs.Code, rs.Msg = 400, I18nUtils.GetMessage("uns.template.alias.already.exists")
		return
	}
	unsDto := &types.CreateTopicDto{}
	copyTemplate2Uns(req, unsDto)
	idRs := l.unsAddService.CreateModelInstance(ctx, unsDto)
	rs = &types.CreateTemplateResp{}
	rs.Code, rs.Msg = 500, "Err"
	if idRs != nil {
		rs.Code, rs.Msg = idRs.Code, idRs.Msg
		rs.Id = idRs.Data.Id
	}
	return
}

func (l *UnsTemplateService) BatchSaveTemplates(ctx context.Context, list []*types.CreateTemplateReq, fromImport bool) (errorMap map[string]string, err error) {
	unsDtos := make([]*types.CreateTopicDto, 0, len(list))
	for _, t := range list {
		uns := &types.CreateTopicDto{}
		copyTemplate2Uns(t, uns)
		unsDtos = append(unsDtos, uns)
	}
	errorMap = l.unsAddService.CreateModelAndInstance(ctx, unsDtos, fromImport)
	return errorMap, nil
}
func copyTemplate2Uns(req *types.CreateTemplateReq, unsDto *types.CreateTopicDto) {
	unsDto.PathType = constants.PathTypeTemplate
	unsDto.Batch, unsDto.Index, unsDto.FlagNo = req.Batch, req.Index, req.FlagNo
	alias := req.Alias
	if alias == "" {
		alias = PathUtil.GenerateAlias(req.Name, 1)
	}
	unsDto.Alias = alias
	unsDto.Fields = req.Fields
	unsDto.Name = req.Name
	if des := req.Description; len(des) > 0 {
		unsDto.Description = &des
	}
}
