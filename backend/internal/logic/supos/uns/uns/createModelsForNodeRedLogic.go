// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/common"
	sysconfig "backend/internal/common/config"
	"backend/internal/common/constants"
	"backend/internal/common/enums"
	"backend/internal/common/utils/PathUtil"
	"backend/internal/logic/supos/uns/uns/bo"
	"backend/internal/logic/supos/uns/uns/service"
	dao "backend/internal/repo/relationDB"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"encoding/json"
	"strings"
	"time"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateModelsForNodeRedLogic struct {
	logx.Logger
	ctx       context.Context
	svcCtx    *svc.ServiceContext
	unsMapper dao.UnsNamespaceRepo
}

// 批量创建文件夹和文件(node-red导入专用)
func NewCreateModelsForNodeRedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateModelsForNodeRedLogic {
	return &CreateModelsForNodeRedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateModelsForNodeRedLogic) CreateModelsForNodeRed(requestDto []*types.CreateUnsNodeRedDto) (resp *types.ResultVO, err error) {
	// 创建Topic
	sysConfig := spring.GetBean[*sysconfig.SystemConfig]()
	createTopicDtos := modelTransfer(l.ctx, sysConfig, l.unsMapper, requestDto)
	resp = &types.ResultVO{Code: 200, Msg: "success"}
	if len(createTopicDtos) == 0 {
		resp.Code, resp.Msg = 400, "NoFields"
		return
	}
	args := bo.CreateModelInstancesArgs{
		Topics:              createTopicDtos,
		FromImport:          false,
		ThrowModelExistsErr: false,
		SkipWhenExists:      true,
		FlowName:            "Nd" + time.Now().Format("20060102150405"), // yyyyMMddHHmmss
		StatusConsumer: func(runningStatus *common.RunningStatus) {
			bs, _ := json.Marshal(runningStatus)
			logx.Info("[NodeRed] ", string(bs))
		},
	}
	errTipMap := spring.GetBean[*service.UnsAddService]().CreateModelAndInstancesInner(l.ctx, args)
	if len(errTipMap) > 0 {
		errJson, _ := json.Marshal(errTipMap)
		resp.Msg = string(errJson)
		if len(errTipMap) == len(createTopicDtos) {
			resp.Code = 400
			return
		}
	}

	// 填充结果
	results := make([][]string, 0, len(requestDto))
	for _, dto := range requestDto {
		fullpath := dto.Path
		if strings.HasSuffix(dto.Path, "/") {
			fullpath = dto.Path[:len(dto.Path)-1]
		}

		alias := dto.Alias
		if alias == "" && dto.FieldName != "" {
			alias = PathUtil.GenerateFileAlias(fullpath)
		}

		// path, alias, fname, ftype, tag
		row := []string{dto.Path, alias, dto.FieldName, dto.FieldType, dto.Tag}
		results = append(results, row)
	}
	resp.Data = results
	return
}

var __metric_ = int16(enums.METRIC)

// modelTransfer 对应Java方法
func modelTransfer(ctx context.Context, sysConfig *sysconfig.SystemConfig, unsMapper dao.UnsNamespaceRepo, requestDtoList []*types.CreateUnsNodeRedDto) []*types.CreateTopicDto {
	dtos := make([]*types.CreateTopicDto, 0, len(requestDtoList))
	aliasMap := make(map[string]string)
	fileMap := make(map[string][]*types.FieldDefine)

	// 根据路径，对文件属性进行归类
	for _, requestDto := range requestDtoList {
		if fullPath := requestDto.Path; len(fullPath) > 0 {
			if last := len(requestDto.Path) - 1; requestDto.Path[last] == '/' {
				fullPath = requestDto.Path[:last]
			}
			ft, ok := types.GetFieldTypeByNameIgnoreCase(requestDto.FieldType)
			if !ok {
				ft = types.FieldTypeString
			}
			if requestDto.FieldName != "" {
				fd := &types.FieldDefine{
					Type: ft.Name(),
					Name: requestDto.FieldName,
				}
				fs, _ := fileMap[fullPath]
				if len(fs) == 0 {
					fieldDefines := make([]*types.FieldDefine, 0, 4)
					fieldDefines = append(fieldDefines, fd)
					fileMap[fullPath] = fieldDefines
				}
			} else {
				_, exists := fileMap[fullPath]
				if !exists {
					fileMap[fullPath] = []*types.FieldDefine{}
				}
			}

			// 如果传了alias，那么就以参数为准
			if alias := requestDto.Alias; alias != "" {
				aliasMap[fullPath] = alias
			}
		}
	}

	if len(fileMap) == 0 {
		return dtos
	}

	var db = dao.GetDb(ctx)
	directFolder := make(map[string]bool, len(fileMap))
	// 构建 path 和 alias 的关系
	for path := range fileMap {
		first := true
		for tmp := path; len(tmp) > 0; {
			if _, exists := aliasMap[tmp]; !exists {
				alias := unsMapper.GetAliasByPath(db, tmp)
				if alias == "" {
					alias = PathUtil.GenerateFileAlias(tmp)
				}
				aliasMap[tmp] = alias
			}
			x := strings.LastIndex(tmp, "/")
			if x < 1 {
				break
			} else {
				tmp = tmp[:x]
				if first {
					first = false
					if len(tmp) > 0 {
						directFolder[tmp] = true
					}
				}
			}
		}
	}
	EnableAutoCategorization := sysConfig.EnableAutoCategorization

	// 创建CreateTopicDto列表
	for path, alias := range aliasMap {
		dto := &types.CreateTopicDto{
			Alias: alias,
			Path:  path,
		}

		fields, exists := fileMap[path]
		if exists && len(fields) > 0 {
			dto.Fields = fields
			dto.PathType = constants.PathTypeFile
			dto.DataType = base.V2p(constants.TimeSequenceType)
			if EnableAutoCategorization {
				dto.ParentDataType = &__metric_
			}
		} else {
			dto.PathType = constants.PathTypeDir
			if EnableAutoCategorization && directFolder[path] {
				dto.DataType = &__metric_
			}
		}

		// 根据当前path获取父级path
		lastSlashIndex := strings.LastIndex(path, "/")
		if lastSlashIndex > 0 {
			parentPath := path[:lastSlashIndex]
			pa := aliasMap[parentPath]
			if len(pa) > 0 {
				dto.ParentAlias = &pa
			}
		}

		index := lastSlashIndex + 1
		if lastSlashIndex >= 0 {
			dto.Name = path[index:]
		} else {
			dto.Name = path
		}

		dtos = append(dtos, dto)
	}

	return dtos
}
