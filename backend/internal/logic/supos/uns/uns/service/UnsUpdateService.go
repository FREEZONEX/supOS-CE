package service

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/logic/supos/uns/uns/UnsConverter"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/spring"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsUpdateService struct {
	log           logx.Logger
	unsMapper     dao.UnsNamespaceRepo
	unsAddService *UnsAddService
}

func init() {
	spring.RegisterLazy[*UnsUpdateService](func() *UnsUpdateService {
		return &UnsUpdateService{
			log:           logx.WithContext(context.Background()),
			unsAddService: spring.GetBean[*UnsAddService](),
		}
	})
}
func (s *UnsUpdateService) UpdateDetail(ctx context.Context, unsDto *types.UpdateUnsDto) (resp *types.StringResult, err error) {
	var unsPo *dao.UnsNamespace
	db := dao.GetDb(ctx)
	ctx = dao.SetDb(ctx, db)
	if unsDto.Id != 0 {
		unsPo, err = s.unsMapper.SelectById(db, unsDto.Id)
	} else if unsDto.Alias != "" {
		unsPo, err = s.unsMapper.GetByAlias(db, unsDto.Alias)
	}

	if err != nil || unsPo == nil {
		resp = &types.StringResult{BaseResult: types.BaseResult{Code: 400, Msg: I18nUtils.GetMessage("uns.folder.or.file.not.found")}}
		return
	}

	createTopicDto := UnsConverter.ConvertApiUpdateDto(unsDto)

	if createTopicDto.Name == "" {
		createTopicDto.Name = unsPo.Name
	}
	createTopicDto.Id = unsPo.Id
	createTopicDto.PathType = unsPo.PathType

	var flags *int32
	dash := createTopicDto.AddDashBoard
	flow := createTopicDto.AddFlow
	save2db := createTopicDto.Save2Db
	accessLevel := unsDto.AccessLevel

	if createTopicDto.WithFlags == nil && (dash != nil || flow != nil || save2db != nil || accessLevel != "") {
		// 使用原有的flags或默认值
		if unsPo.WithFlags != nil {
			flags = unsPo.WithFlags
		} else {
			zero := int32(0)
			flags = &zero
		}

		fl := *flags
		if dash != nil {
			if *dash {
				fl = fl | constants.UnsFlagWithDashboard
			} else {
				fl = fl &^ constants.UnsFlagWithDashboard
			}
		}
		if flow != nil {
			if *flow {
				fl = fl | constants.UnsFlagWithFlow
			} else {
				fl = fl &^ constants.UnsFlagWithFlow
			}
		}
		if save2db != nil {
			if *save2db {
				fl = fl | constants.UnsFlagWithSave2DB
			} else {
				fl = fl &^ constants.UnsFlagWithSave2DB
			}
		}
		if accessLevel != "" {
			if accessLevel == enums.FileModeReadOnly.String() {
				fl = (fl &^ constants.UnsFlagAccessLevelReadWrite) | constants.UnsFlagAccessLevelReadOnly
			} else {
				fl = (fl &^ constants.UnsFlagAccessLevelReadOnly) | constants.UnsFlagAccessLevelReadWrite
			}
		}
		flags = &fl
	} else {
		flags = unsPo.WithFlags
	}

	createTopicDto.WithFlags = flags
	createTopicDto.ParentAlias = unsPo.ParentAlias
	createTopicDto.DataType = unsPo.DataType

	createUnsResp := s.unsAddService.CreateModelInstance(ctx, createTopicDto)
	resp = &types.StringResult{BaseResult: createUnsResp.BaseResult, Data: createUnsResp.Data.Id}
	return
}

func (s *UnsUpdateService) UpdateName(ctx context.Context, req *types.UpdateNameVo) (resp *types.StringResult, err error) {
	db := dao.GetDb(ctx)
	ctx = dao.SetDb(ctx, db)
	unsPo, err := s.unsMapper.SelectById(db, req.Id)
	if err != nil || unsPo == nil {
		resp = &types.StringResult{BaseResult: types.BaseResult{Code: 400, Msg: I18nUtils.GetMessage("uns.folder.or.file.not.found")}}
		return
	}

	createTopicDto := &types.CreateTopicDto{
		Name:        req.Name,
		WithFlags:   unsPo.WithFlags,
		PathType:    unsPo.PathType,
		Alias:       unsPo.Alias,
		ParentAlias: unsPo.ParentAlias,
		DataType:    unsPo.DataType,
	}

	createUnsResp := s.unsAddService.CreateModelInstance(ctx, createTopicDto)
	resp = &types.StringResult{BaseResult: createUnsResp.BaseResult, Data: createUnsResp.Data.Id}
	return
}

func (s *UnsUpdateService) SubscribeModel(ctx context.Context, req *types.SubscribeModelReq) (resp *types.ResultVO, err error) {
	return
}
