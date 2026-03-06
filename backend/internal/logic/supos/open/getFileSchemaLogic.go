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

type GetFileSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询文件schema 元数据结构
func NewGetFileSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileSchemaLogic {
	return &GetFileSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileSchemaLogic) GetFileSchema() (resp map[string]interface{}, err error) {
	// 根据语言环境变量选择对应的schema文件
	lang := os.Getenv("SYS_OS_LANG")
	var filename string

	if lang == "en-US" {
		filename = "file-schema_en.json"
	} else {
		filename = "file-schema.json"
	}

	// 构建文件路径（使用嵌入的文件系统）
	filePath := "templates/" + filename

	// 读取文件内容
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
