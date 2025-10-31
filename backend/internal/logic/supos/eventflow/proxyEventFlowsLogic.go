// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package eventflow

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

type ProxyEventFlowsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// Proxy Node-RED event /flows endpoint using cookie scoped id
func NewProxyEventFlowsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProxyEventFlowsLogic {
	return &ProxyEventFlowsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProxyEventFlowsLogic) ProxyEventFlows(flowID string) (string, error) {
	client := l.svcCtx.EventNodeRed
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
		l.Errorf("proxy event flows failed: code=%d err=%v body=%s", code, errs, string(body))
		return "", errors.System.WithMsg(i18ns.LocalizeMsg("error.sys.systemError")).AddDetail(string(body))
	}
	if out == nil {
		return string(body), nil
	}
	data, err := json.Marshal(out)
	if err != nil {
		l.Errorf("marshal event node-red response failed: %v", err)
		return string(body), nil
	}
	return string(data), nil
}
