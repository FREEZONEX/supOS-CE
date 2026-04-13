// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/serviceApi"
	"backend/internal/common/utils/finddatautil"
	logiccommon "backend/internal/logic/supos/uns/common"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"

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
var _initOnce sync.Once

func (l *BatchUpdateFileLogic) BatchUpdateFile(req *types.UpdateFileDTOArray) (resp *types.BatchUpdateFileResult, err error) {
	resp = &types.BatchUpdateFileResult{}
	resp.Code = 200
	resp.Msg = "ok"

	list := req.Items
	if len(list) == 0 {
		return
	} else if len(list) > 100 {
		resp.Code = 400
		resp.Msg = "仅支持最大一次性对100个文件进行数据更新"
		return
	}

	// 初始化服务
	_initOnce.Do(func() {})
	if defService == nil {
		_initOnce.Do(func() {
			defService = spring.GetBean[serviceApi.IUnsDefinitionService]()
		})
	}

	// 获取所有别名
	aliasList := base.Map(list, func(e types.UpdateFileDTO) string {
		return e.Alias
	})

	// 检查不存在的别名
	notExists := base.Filter(aliasList, func(alias string) bool {
		return defService.GetDefinitionByAlias(alias) == nil
	})

	ctx := l.ctx
	errorFields := make(map[string]string)

	// 遍历处理每个文件数据
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

		// 关系型数据需要主键
		if types.GetSrcJdbcTypeByID(def.DataSrcID).TypeCode() == constants.RelationType {
			for k := range def.GetFieldDefines().UniqueKeys {
				if !strings.HasPrefix(k, constants.SystemFieldPrev) && !base.MapContainsKey(body, k) {
					resp.Code = 400
					resp.Msg = ""
					resp.Data = types.UnsDataResponseVo{NotExists: notExists, ErrorFields: errorFields}
					errorFields[alias+"."+k] = I18nUtils.GetMessageWithCtx(ctx, "uns.write.value.relation.pk.is.null")
					return
				}
			}
		}

		// 处理字段数据
		newBody := make(map[string]interface{})
		qosField := def.GetQualityField()
		for fieldName, v := range body {
			fieldDefine := fMap[fieldName]
			if fieldDefine == nil {
				errorFields[alias+"."+fieldName] = I18nUtils.GetMessageWithCtx(ctx, "uns.field.not.found")
				continue
			}

			// 处理QoS字段
			if fieldName == qosField {
				if qosStr, ok := v.(string); ok {
					hex, er := strconv.ParseInt(qosStr, 16, 64)
					if er != nil {
						errorFields[alias+"."+fieldName] = I18nUtils.GetMessageWithCtx(ctx, "uns.field.type.un.match")
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

		// 数据验证
		rs := finddatautil.FindDataList(newBody, 1, def.GetFieldDefines())
		if fieldName := rs.ErrorField; len(fieldName) > 0 {
			errorFields[alias+"."+fieldName] = I18nUtils.GetMessageWithCtx(ctx, "uns.field.type.un.match")
		}
		if fieldName := rs.ToLongField; len(fieldName) > 0 {
			errorFields[alias+"."+fieldName] = I18nUtils.GetMessageWithCtx(ctx, "uns.field.value.out.of.size")
		}

		// 通过 MQTT 发送消息，由现有订阅链路统一回流处理
		jsonBs, _ := json.Marshal(newBody)
		if err = logiccommon.PublishUnsMessage(l.svcCtx, def.GetTopic(), jsonBs); err != nil {
			errorFields[alias] = err.Error()
		}
	}

	// 处理结果
	if len(notExists)+len(errorFields) > 0 {
		resp.Code = 206
		resp.Msg = ""
		resp.Data = types.UnsDataResponseVo{NotExists: notExists, ErrorFields: errorFields}
	}

	return
}
