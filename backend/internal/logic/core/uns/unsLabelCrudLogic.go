package uns

import (
	"context"

	domainuns "backend/internal/domain/uns"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsLabelCrudLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnsLabelCrudLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnsLabelCrudLogic {
	return &UnsLabelCrudLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnsLabelCrudLogic) Detail(req *types.IdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.LabelDetail(l.ctx, req.Id)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}

func (l *UnsLabelCrudLogic) Create(req *types.UnsLabelReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.CreateLabel(l.ctx, domainuns.LabelCommand{
		Name: req.Name, Color: req.Color, Description: req.Description, UserID: logicx.UserID(l.ctx),
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}

func (l *UnsLabelCrudLogic) Update(req *types.UnsLabelReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.UpdateLabel(l.ctx, domainuns.LabelCommand{
		ID: req.Id, Name: req.Name, Color: req.Color, Description: req.Description, UserID: logicx.UserID(l.ctx),
	})
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}

func (l *UnsLabelCrudLogic) Delete(req *types.IdReq) (resp *types.Envelope, err error) {
	data, err := l.svcCtx.App.UNS.DeleteLabel(l.ctx, req.Id)
	if err != nil {
		return nil, logicx.Error(err)
	}
	return respx.Envelope(data), nil
}
