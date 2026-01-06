package handler

import (
	"fmt"
	"io"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type NodeRedHandler struct {
	nodeRedHost string
	nodeRedPort string
	client      *resty.Client
}

func NewNodeRedHandler(nodeRedHost, nodeRedPort string) *NodeRedHandler {
	return &NodeRedHandler{
		nodeRedHost: nodeRedHost,
		nodeRedPort: nodeRedPort,
		client:      resty.New(),
	}
}

// ProxyNodeRedFlowsHandler 代理 NodeRed /flows 接口
// 对应 GET /test/nodered 和 /flows/test/nodered
// 通过此代理方法，实现请求 nodeRed /flows 接口，只返回当前流程ID的数据，ID从cookie中获取；
// 如果cookie中不包含ID，则返回空数组
func (h *NodeRedHandler) ProxyNodeRedFlowsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookies := r.Cookies()

		// 构建要转发的 cookie
		var cookieStrs []string
		for _, cookie := range cookies {
			cookieStrs = append(cookieStrs, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
			logx.Infof("cookie key = %s, value = %s", cookie.Name, cookie.Value)
		}

		// 请求 NodeRed 的 /flows 接口
		nodeRedURL := fmt.Sprintf("http://%s:%s/flows", h.nodeRedHost, h.nodeRedPort)
		resp, err := h.client.R().
			SetHeader("Cookie", fmt.Sprintf("%v", cookieStrs)).
			Get(nodeRedURL)

		if err != nil {
			logx.Errorf("request node-red failed: %v", err)
			httpx.Error(w, err)
			return
		}

		logx.Infof("<=== Get node response: %s", resp.Body())

		// 直接返回 NodeRed 的响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode())
		_, _ = io.WriteString(w, resp.String())
	}
}
