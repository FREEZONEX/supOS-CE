// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"os"
	"path/filepath"

	"backend/internal/svc"
	"backend/internal/types"

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

func (l *GetFileSchemaLogic) GetFileSchema() (resp *types.ResultVO, err error) {
	// 根据语言环境变量选择对应的schema文件
	lang := os.Getenv("SYS_OS_LANG")
	var filename string

	if lang == "en-US" {
		filename = "file-schema_en.json"
	} else {
		filename = "file-schema.json"
	}

	// 构建文件路径
	filePath := filepath.Join("internal", "logic", "supos", "open", "templates", filename)

	// 读取文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		l.Errorf("Failed to read file: %v", err)
		return &types.ResultVO{
			Code: 500,
			Msg:  "Failed to read file schema",
			Data: nil,
		}, err
	}

	// 返回结果
	return &types.ResultVO{
		Code: 200,
		Msg:  "Success",
		Data: string(content),
	}, nil
}
