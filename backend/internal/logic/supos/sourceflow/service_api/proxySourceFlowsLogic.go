// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package service_api

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"backend/internal/svc"

	"gitee.com/unitedrhino/share/errors"
	"gitee.com/unitedrhino/share/i18ns"
	"github.com/zeromicro/go-zero/core/logx"
)

type ProxySourceFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Proxy Node-RED /flows endpoint using cookie scoped id
func NewProxySourceFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxySourceFlowsLogic {
	return &ProxySourceFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxySourceFlowsLogic) ProxySourceFlows(flowID string) (string, error) {
	client := l.svcCtx.SourceNodeRed
	if client == nil {
		return "[]", nil
	}
	path := "/flows"
	if strings.TrimSpace(flowID) != "" {
		path = fmt.Sprintf("/flow/%s", strings.TrimSpace(flowID))
	}
	var out any
	code, body, errs := client.DoJSON(l.ctx, "GET", path, nil, &out)
	if len(errs) > 0 || (code != 200 && code != 204) {
		l.Errorf("proxy source flows failed: code=%d err=%v body=%s", code, errs, string(body))
		return "", errors.System.WithMsg(i18ns.LocalizeMsg("error.sys.systemError")).AddDetail(string(body))
	}
	if out == nil {
		return string(body), nil
	}
	data, err := json.Marshal(out)
	if err != nil {
		l.Errorf("marshal node-red response failed: %v", err)
		return string(body), nil
	}
	return string(data), nil
}
