package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"context"
	"strconv"

	"gitee.com/unitedrhino/share/errors"
)

func (l *UnsLabelService) AllLabel(ctx context.Context, req *types.UnsLabelListReq) (resp *types.UnsLabelListResp, err error) {
	db := dao.GetDb(ctx)
	list, _ := l.labelMapper.LikeName(db, req.Key)
	resp = &types.UnsLabelListResp{
		Data: base.Map[*dao.UnsLabel, *types.LabelVo](list, func(e *dao.UnsLabel) *types.LabelVo {
			return UnsConverter.LabelPo2Vo(e)
		}),
	}
	resp.Code = 200
	return
}
func (l *UnsLabelService) Detail(ctx context.Context, req *types.WithID) (resp *types.LabelDetailResult, err error) {
	if req.ID <= 0 {
		return nil, errors.Parameter.WithMsg("id无效")
	}
	db := dao.GetDb(ctx)
	resp = &types.LabelDetailResult{}
	item, _ := l.labelMapper.FindOne(db, req.ID)
	if item == nil {
		resp.Code, resp.Msg = 400, I18nUtils.GetMessage("uns.label.not.exists")
		return nil, err
	}
	resp.Code = 200
	resp.Data = UnsConverter.LabelPo2Vo(item)
	return
}

func (l *UnsLabelService) PageListUnsByLabel(ctx context.Context, req *types.LabelPageReq) (rs *types.UnsByLabelPageResp, err error) {
	db := dao.GetDb(ctx)
	labelId, pageNo, pageSize := req.LabelId, req.PageNo, req.PageSize
	var countTotal = int64(0)
	var unsList []*dao.UnsNamespace
	unsList, err = l.unsMapper.PageListByLabel(db, labelId, pageNo, pageSize, &countTotal)
	if err != nil {
		return
	}
	rs = &types.UnsByLabelPageResp{PageNo: pageNo, PageSize: pageSize, Total: countTotal}
	rs.Code, rs.Msg = 200, "OK"
	rs.Data = base.Map[*dao.UnsNamespace, types.FileVo](unsList, func(uns *dao.UnsNamespace) types.FileVo {
		return types.FileVo{
			UnsId:    strconv.FormatInt(uns.Id, 10),
			PathType: uns.PathType,
			Path:     uns.Path,
			Name:     PathUtil.GetName(uns.Path),
		}
	})
	return
}
