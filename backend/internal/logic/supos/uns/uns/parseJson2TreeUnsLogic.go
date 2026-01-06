// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package uns

import (
	"context"
	"encoding/json"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParseJson2TreeUnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 外部JSON定义转树结构uns字段定义
func NewParseJson2TreeUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseJson2TreeUnsLogic {
	return &ParseJson2TreeUnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ParseJson2TreeUnsLogic) ParseJson2TreeUns(jsonBs []byte) (resp *types.ParseJson2TreeUnsResp, err error) {
	resp = &types.ParseJson2TreeUnsResp{}
	resp.Code, resp.Msg = 200, "ok"
	resp.Data, err = parserJSON2Tree(jsonBs)
	if err != nil {
		resp.Code, resp.Msg = 400, err.Error()
	}
	return
}

func parserJSON2Tree(jsonBs []byte) ([]*types.TreeOuterStructureVo, error) {
	// 检查JSON是否有效
	if !json.Valid(jsonBs) {
		return nil, nil
	}

	// 创建default作为根节点
	defaultNode := &types.TreeOuterStructureVo{
		Name:     "default",
		DataPath: "default",
		Fields:   []types.FieldDefine{},
		Children: []*types.TreeOuterStructureVo{},
	}

	// 解析JSON
	var data interface{}
	if err := json.Unmarshal(jsonBs, &data); err != nil {
		return nil, err
	}

	// 根据类型进行解析
	switch v := data.(type) {
	case map[string]interface{}:
		// 直接解析JSON对象
		parseJSONObject(v, "default", defaultNode)
	case []interface{}:
		// JSON数组
		if len(v) > 0 {
			// 取数组第一个元素进行解析
			if firstItem, ok := v[0].(map[string]interface{}); ok {
				parseJSONObject(firstItem, "default", defaultNode)
			}
		}
	}

	// 结果只返回default这个根节点
	result := []*types.TreeOuterStructureVo{defaultNode}
	return result, nil
}

// parseJSONObject 递归解析JSON对象
func parseJSONObject(jsonObj map[string]interface{}, parentPath string, parentNode *types.TreeOuterStructureVo) {
	var children []*types.TreeOuterStructureVo
	var fields []types.FieldDefine

	for key, value := range jsonObj {
		currentPath := parentPath + "." + key

		// 处理数组：取数组第一个元素（如果是JSON对象）
		if arr, isArray := value.([]interface{}); isArray && len(arr) > 0 {
			// 如果数组的第一个元素是对象，则使用该对象进行后续解析
			if firstItem, ok := arr[0].(map[string]interface{}); ok {
				value = firstItem
			}
		}

		// 根据类型处理
		switch v := value.(type) {
		case map[string]interface{}:
			// 处理子对象
			childNode := &types.TreeOuterStructureVo{
				Name:     key,
				DataPath: currentPath,
				Fields:   []types.FieldDefine{},
				Children: []*types.TreeOuterStructureVo{},
			}

			parseJSONObject(v, currentPath, childNode)
			children = append(children, childNode)

		default:
			// 处理普通字段
			field := types.FieldDefine{
				Name: key,
				Type: string(guessType(v)), // 使用已有的guessType函数
			}
			fields = append(fields, field)
		}
	}

	parentNode.Children = children
	parentNode.Fields = fields
}

// 辅助函数：检查JSON是否有效
func isValidJSON(jsonStr string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}
