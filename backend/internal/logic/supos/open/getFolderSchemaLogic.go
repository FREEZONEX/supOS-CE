// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"context"
	"embed"
	"os"

	"backend/internal/svc"
	"backend/internal/types"

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

func (l *GetFolderSchemaLogic) GetFolderSchema() (resp *types.ResultVO, err error) {
	// 根据语言环境变量选择对应的schema文件
	lang := os.Getenv("SYS_OS_LANG")
	var filename string

	if lang == "en-US" {
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
		return &types.ResultVO{
			Code: 500,
			Msg:  "Failed to read folder schema",
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
