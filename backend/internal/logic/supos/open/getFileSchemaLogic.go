// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"backend/internal/common/utils/apiutil"
	"backend/internal/common/utils/langutil"
	"backend/internal/svc"
	"context"
	"encoding/json"
	"strings"

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
	// 获取用户语言，优先使用用户配置的语言，为空则使用系统语言
	lang := langutil.SystemLocale()
	if user := apiutil.GetUserFromContext(l.ctx); user != nil && user.MainLanguage != "" {
		lang = user.MainLanguage
	}
	var filename string

	if strings.Contains(lang, "en") {
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
