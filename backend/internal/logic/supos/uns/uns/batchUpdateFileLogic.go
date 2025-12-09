// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/finddatautil"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchUpdateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 批量写文件实时值
func NewBatchUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchUpdateFileLogic {
	return &BatchUpdateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var defService serviceApi.IUnsDefinitionService
var msgConsumer serviceApi.TopicMessageConsumer
var _initOnce sync.Once

func (l *BatchUpdateFileLogic) BatchUpdateFile(list []types.UpdateFileDTO) (resp *types.UnsDataResponse, err error) {
	resp = &types.UnsDataResponse{}
	resp.Code, resp.Msg = 200, "ok"
	if len(list) == 0 {
		return
	} else if len(list) > 100 {
		resp.Code, resp.Msg = 400, "仅支持最大一次性对100个文件进行数据更新"
		return
	}
	initOnce.Do(func() {})
	if defService == nil {
		_initOnce.Do(func() {
			defService = spring.GetBean[serviceApi.IUnsDefinitionService]()
			msgConsumer = spring.GetBean[serviceApi.TopicMessageConsumer]()
		})
	}
	aliasList := base.Map(list, func(e types.UpdateFileDTO) string {
		return e.Alias
	})
	notExists := base.Filter(aliasList, func(alias string) bool {
		return defService.GetDefinitionByAlias(alias) == nil
	})
	errorFields := make(map[string]string, 8)
	for _, dto := range list {
		body := dto.Data
		if len(body) == 0 {
			continue
		}
		alias := dto.Alias
		def := defService.GetDefinitionByAlias(alias)
		if def == nil || def.DataType == nil {
			continue
		}
		fMap := def.GetFieldDefines().FieldsMap
		if types.GetSrcJdbcTypeByID(def.DataSrcID).TypeCode() == constants.RelationType {
			for k := range def.GetFieldDefines().UniqueKeys {
				if !strings.HasPrefix(k, constants.SystemFieldPrev) && !base.MapContainsKey(body, k) {
					resp.Code, resp.Msg = 400, ""
					resp.Data = &types.UnsDataResponseVo{NotExists: notExists, ErrorFields: errorFields}
					errorFields[k] = I18nUtils.GetMessage("uns.write.value.relation.pk.is.null")
					return
				}
			}
		}
		newBody := make(map[string]interface{}, 8)
		qosField := def.GetQualityField()
		for fieldName, v := range body {
			fieldDefine := fMap[fieldName]
			if fieldDefine == nil {
				errorFields[alias+"."+fieldName] = I18nUtils.GetMessage("uns.field.not.found")
				continue
			}
			if fieldName == qosField {
				if qosStr, ok := v.(string); ok {
					hex, er := strconv.ParseInt(qosStr, 16, 64)
					if er != nil {
						errorFields[alias+"."+fieldName] = I18nUtils.GetMessage("uns.field.type.un.match")
						continue
					}
					v = hex
				}
			}
			newBody[fieldName] = v
		}
		if len(newBody) == 0 {
			continue
		}
		rs := finddatautil.FindDataList(newBody, 1, def.GetFieldDefines())
		if fieldName := rs.ErrorField; len(fieldName) > 0 {
			errorFields[alias+"."+fieldName] = I18nUtils.GetMessage("uns.field.type.un.match")
		}
		if fieldName := rs.ToLongField; len(fieldName) > 0 {
			errorFields[alias+"."+fieldName] = I18nUtils.GetMessage("uns.field.value.out.of.size")
		}
		jsonBs, _ := json.Marshal(newBody)
		msgConsumer.OnMessageByAlias(alias, string(jsonBs))
	}
	if len(notExists)+len(errorFields) > 0 {
		resp.Code, resp.Msg = 206, ""
		resp.Data = &types.UnsDataResponseVo{NotExists: notExists, ErrorFields: errorFields}
	}
	return
}
