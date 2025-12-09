// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package system

import (
	"backend/internal/common/serviceApi"
	"backend/share/spring"
	"context"
	"errors"
	"strconv"
	"unicode"

	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUnsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnsLogic {
	return &GetUnsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

var defService serviceApi.IUnsDefinitionService

func (l *GetUnsLogic) GetUns(req *types.GetDefRequest) (resp *types.CreateTopicDto, err error) {
	if defService == nil {
		defService = spring.GetBean[serviceApi.IUnsDefinitionService]()
	}
	if unicode.IsDigit(rune(req.Uns[0])) {
		id, _ := strconv.ParseInt(req.Uns, 10, 64)
		if id > 0 {
			resp = defService.GetDefinitionById(id)
		}
	} else {
		resp = defService.GetDefinitionByAlias(req.Uns)
	}
	if resp == nil {
		err = errors.New("UnsDefinitionNotFound")
	}
	return
}
