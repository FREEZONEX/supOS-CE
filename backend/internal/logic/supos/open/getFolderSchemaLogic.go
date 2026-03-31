// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/utils/langutil"
	"backend/internal/svc"
	"context"
	"embed"
	"encoding/json"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

//go:embed templates/*
var templates embed.FS

type GetFolderSchemaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 查询文件夹schema 元数据结构
func NewGetFolderSchemaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFolderSchemaLogic {
	return &GetFolderSchemaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFolderSchemaLogic) GetFolderSchema() (resp map[string]interface{}, err error) {
	// 获取用户语言，优先使用用户配置的语言，为空则使用系统语言
	lang := langutil.SystemLocale()
	if user := apiutil.GetUserFromContext(l.ctx); user != nil && user.MainLanguage != "" {
		lang = user.MainLanguage
	}
	var filename string

	if strings.Contains(lang, "en") {
		filename = "folder-schema_en.json"
	} else {
		filename = "folder-schema.json"
	}

	// 构建文件路径
	filePath := "templates/" + filename

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
