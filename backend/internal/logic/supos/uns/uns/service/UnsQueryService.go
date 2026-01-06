package service

import (
	"backend/internal/adapters/msg_consumer"
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"errors"
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
	calcService  UnsCalcService
	defService   serviceApi.IUnsDefinitionService
}

func init() {
	spring.RegisterLazy[*UnsQueryService](func() *UnsQueryService {
		return &UnsQueryService{
			log:        logx.WithContext(context.Background()),
			defService: spring.GetBean[*msg_consumer.UnsDefinitionService](),
		}
	})
}

var unsNotFoundError = errors.New("unsNotFoundError")

// GetLastMsg returns the last message for a UNS by ID
func (l *UnsQueryService) GetLastMsg(id int64) ([]byte, error) {
	def := l.defService.GetDefinitionById(id)
	if def == nil {
		return nil, unsNotFoundError
	}
	return l.getLastMsg(def), nil
}

// GetLastMsgByAlias returns the last message for a UNS by alias
func (l *UnsQueryService) GetLastMsgByAlias(alias string) ([]byte, error) {
	def := l.defService.GetDefinitionByAlias(alias)
	if def == nil {
		return nil, unsNotFoundError
	}
	return l.getLastMsg(def), nil
}

// GetLastMsgByPath returns the last message for a UNS by path
func (l *UnsQueryService) GetLastMsgByPath(path string) ([]byte, error) {
	def := l.defService.GetDefinitionByPath(path)
	if def == nil {
		return nil, unsNotFoundError
	}
	return l.getLastMsg(def), nil
}
func (l *UnsQueryService) getLastMsg(def *types.CreateTopicDto) []byte {
	wsMsg := serviceApi.WebsocketMessage{Def: def}
	return processWsMsg(wsMsg)
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
	resp = &types.TopicPaginationSearchResult{Code: 500, Page: &types.PageResultDTO{PageNo: int64(req.PageNo), PageSize: int64(req.PageSize)}}
	switch req.SearchType {
	case 0, 2: // 查询普通文件或目录，可指定dataType
		var list []*dao.SimpleUns
		pageInfo.Orders = []stores.OrderBy{
			{Field: "path", Sort: stores.OrderAsc}, {Field: "create_at", Sort: stores.OrderDesc},
		}
		filter := &dao.UnsPathFilter{Key: keyword, PathType: req.SearchType, DataTypes: req.DataTypes}
		if req.Normal && len(req.DataTypes) == 0 {
			filter.DataTypes = []int16{constants.TimeSequenceType, constants.RelationType, constants.JsonbType}
		}
		list, err = l.unsMapper.ListPaths(db, filter, pageInfo, &resp.Page.Total)
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
		list, err = l.unsMapper.ListNotCalcSeqFiles(db, keyword, req.NumberFieldCount, pageInfo, &resp.Page.Total)
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		resp.Data = l.po2TopicInfo(list)
	case 4: // 查询时序类文件
		var list []*dao.UnsNamespace
		list, err = l.unsMapper.ListTimeSeriesFiles(db, keyword, pageInfo, &resp.Page.Total)
		if err != nil {
			resp.Msg = err.Error()
			return
		}
		resp.Data = l.po2TopicInfo(list)
	case 5: //Alarm
		var list []*dao.UnsNamespace
		list, err = l.unsMapper.ListAlarmRules(db, keyword, pageInfo, &resp.Page.Total)
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
		var fs = base.FilterAndMap[*types.FieldDefine, *types.FieldDefine](e.Fields, func(e *types.FieldDefine) (v *types.FieldDefine, ok bool) {
			if !e.IsSystemField() && (types.FieldType(e.Type).IsNumber() || e.Type == types.FieldTypeBoolean) {
				ok = true
				v = &types.FieldDefine{Name: e.Name, Type: e.Type}
			}
			return
		})
		if len(fs) == 0 {
			return
		}
		ok = true
		v = types.TopicInfo{
			Id: strconv.FormatInt(e.Id, 10),
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
		if v.DataType == int(constants.AlarmRuleType) {
			//TODO Alarm
			v.Topic = e.Path
			v.Description = base.P2v(e.Description)
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
