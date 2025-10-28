package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/dto"
	"backend/internal/common/enums"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"strconv"
	"strings"

	"gitee.com/unitedrhino/share/stores"
	"github.com/zeromicro/go-zero/core/logx"
)

type UnsQueryService struct {
	log          logx.Logger
	unsMapper    dao.UnsNamespaceRepo
	labelMapper  dao.UnsLabelRepo
	mountService UnsMountService
}

func init() {
	spring.RegisterBean(&UnsQueryService{
		log: logx.WithContext(context.Background()),
	})
}

func (l *UnsQueryService) SearchPaged(ctx context.Context, req *types.SearchPagedReq) (resp *types.TopicPaginationSearchResult, err error) {
	db := dao.GetDb(ctx)
	keyword := req.Key
	if len(keyword) > 0 {
		keyword = strings.Replace(keyword, "_", "\\_", -1)
		keyword = strings.Replace(keyword, "%", "\\%", -1)
		keyword = "%" + keyword + "%"
	}
	pageInfo := &stores.PageInfo{Page: int64(req.PageNo), Size: int64(req.PageSize), Orders: []stores.OrderBy{
		{Field: "create_at", Sort: stores.OrderDesc},
	}}
	resp = &types.TopicPaginationSearchResult{Code: 500, PageNo: int64(req.PageNo), PageSize: int64(req.PageSize)}
	switch req.SearchType {
	case 0, 2: // 查询普通文件或目录，可指定dataType
		var list []*dao.SimpleUns
		pageInfo.Orders = []stores.OrderBy{
			{Field: "path", Sort: stores.OrderAsc}, {Field: "create_at", Sort: stores.OrderDesc},
		}
		list, err = l.unsMapper.ListPaths(db, &dao.UnsPathFilter{Key: keyword, PathType: req.SearchType, DataTypes: req.DataTypes}, pageInfo, &resp.Total)
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		resp.Data = base.Map[*dao.SimpleUns, types.TopicInfo](list, func(e *dao.SimpleUns) types.TopicInfo {
			return types.TopicInfo{
				Id:       e.ID,
				DataType: e.DataType,
				Alias:    e.Alias,
				Path:     e.Path,
				Topic:    base.SanYuan(constants.UseAliasAsTopic, e.Alias, e.Path),
			}
		})
	case 3: // 为计算类型的下拉框查询其他类型文件
		var list []*dao.UnsNamespace
		list, err = l.unsMapper.ListNotCalcSeqFiles(db, keyword, req.NumberFieldCount, pageInfo, &resp.Total)
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		resp.Data = l.po2TopicInfo(list)
	case 4: // 查询时序类文件
		var list []*dao.UnsNamespace
		list, err = l.unsMapper.ListTimeSeriesFiles(db, keyword, pageInfo, &resp.Total)
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		resp.Data = l.po2TopicInfo(list)
	case 5: //Alarm
		var list []*dao.UnsNamespace
		list, err = l.unsMapper.ListAlarmRules(db, keyword, pageInfo, &resp.Total)
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		resp.Data = l.po2TopicInfo(list)
	}
	resp.Code = 200
	return
}

func (l *UnsQueryService) po2TopicInfo(list []*dao.UnsNamespace) []types.TopicInfo {
	return base.FilterAndMap[*dao.UnsNamespace, types.TopicInfo](list, func(e *dao.UnsNamespace) (v types.TopicInfo, ok bool) {
		var fs = base.FilterAndMap[*dto.FieldDefine, *types.FieldDefine](e.Fields, func(e *dto.FieldDefine) (v *types.FieldDefine, ok bool) {
			if !e.IsSystemField() && (e.Type.IsNumber() || e.Type == enums.FieldTypeBoolean) {
				ok = true
				v = &types.FieldDefine{Name: e.Name, Type: e.Type.Name()}
			}
			return
		})
		if len(fs) == 0 {
			return
		}
		ok = true
		v = types.TopicInfo{
			Id: strconv.FormatInt(e.ID, 10),
			DataType: base.SanA(e.DataType == nil, 0, func() int {
				return int(*e.DataType)
			}),
			ParentDataType: e.ParentDataType,
			Alias:          e.Alias,
			Path:           e.Path,
			Name:           e.Name,
			Fields:         fs,
			Topic:          base.SanYuan(constants.UseAliasAsTopic, e.Alias, e.Path),
		}
		return
	})
}

func (l *UnsQueryService) SearchEmptyFolder(ctx context.Context, req *types.EmptyFolderReq) (resp *types.EmptyFolderResp, err error) {
	db := dao.GetDb(ctx)
	var list []*dao.UnsNamespace
	list, err = l.unsMapper.ListAllEmptyFolder(db)
	resp = &types.EmptyFolderResp{}
	resp.Code = 500
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Code = 200
	resp.Data = UnsConverter.Po2ApiDtos(list)
	return
}
