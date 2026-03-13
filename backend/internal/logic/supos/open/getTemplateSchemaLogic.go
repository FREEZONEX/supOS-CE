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

type GetTemplateSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询模板schema 元数据结构
func NewGetTemplateSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTemplateSchemaLogic {
	return &GetTemplateSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTemplateSchemaLogic) GetTemplateSchema() (resp map[string]interface{}, err error) {
	// 根据语言环境变量选择对应的schema文件
	lang := os.Getenv("SYS_OS_LANG")
	var filename string

	if lang == "en-US" {
		filename = "template-schema_en.json"
	} else {
		filename = "template-schema.json"
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
