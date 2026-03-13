// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"encoding/json"
	"os"

	"backend/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLabelSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询标签schema 元数据结构
func NewGetLabelSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLabelSchemaLogic {
	return &GetLabelSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLabelSchemaLogic) GetLabelSchema() (resp map[string]interface{}, err error) {
	// 根据语言环境变量选择对应的schema文件
	lang := os.Getenv("SYS_OS_LANG")
	var filename string

	if lang == "en-US" {
		filename = "label-schema_en.json"
	} else {
		filename = "label-schema.json"
	}

	// 构建文件路径
	filePath := "templates/uns/" + filename

	// 读取文件内容（使用embed包读取嵌入的文件）
	content, err := templates.ReadFile(filePath)
	if err != nil {
		l.Errorf("Failed to read file: %v", err)
		return nil, err
	}

	// 解析JSON内容
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	if err != nil {
		l.Errorf("Failed to parse JSON: %v", err)
		return nil, err
	}

	// 直接返回JSON数据
	return jsonData, nil
}
