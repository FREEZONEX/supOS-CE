package nodered

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/zeromicro/go-zero/core/logx"
)

type NodeRedConfig struct {
	Source SourceFlowConfig `json:",optional" `
	Event  EventFlowConfig  `json:",optional" `
}

type EventFlowConfig struct {
	Host string `json:",optional,env=NODE_RED_EVENT_HOST,default=eventflow" `
	Port string `json:",optional,env=NODE_RED_EVENT_PORT,default=1889" `
}

type SourceFlowConfig struct {
	Host string `json:",optional,env=NODE_RED_HOST,default=nodered" `
	Port string `json:",optional,env=NODE_RED_PORT,default=1880" `
}

type Client struct {
	BaseURL string        // 例如: http://nodered:1880 或 http://gateway/nodered/home
	Timeout time.Duration // 默认 10m
	Retry   int           // 默认 1
}

func NewClient(host, port, prefix string) *Client {
	base := host
	if port != "" && !strings.Contains(host, "://") {
		base = fmt.Sprintf("http://%s:%s", host, port)
	}
	if prefix != "" {
		base = strings.TrimRight(base, "/") + "/" + strings.TrimLeft(prefix, "/")
	}
	return &Client{BaseURL: base, Timeout: 10 * time.Minute, Retry: 1}
}

func (c *Client) url(path string) string {
	return strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

func (c Client) DoJSON(ctx context.Context, method, path string, reqBody any, out any) (int, []byte, []error) {
	greq := gorequest.New().Retry(c.Retry, time.Second*2).Timeout(c.Timeout)
	logx.WithContext(ctx).Debugf("nodered %s %s", method, c.url(path))
	var resp gorequest.Response
	var body []byte
	var errs []error
	switch strings.ToUpper(method) {
	case "GET":
		resp, body, errs = greq.Get(c.url(path)).EndStruct(out)
	case "POST":
		resp, body, errs = greq.Post(c.url(path)).Send(reqBody).EndStruct(out)
	case "PUT":
		resp, body, errs = greq.Put(c.url(path)).Send(reqBody).EndStruct(out)
	case "DELETE":
		// Delete often returns an empty body; using EndStruct can emit unexpected EOF errors.
		resp, body, errs = greq.Delete(c.url(path)).EndBytes()
	default:
		errs = append(errs, fmt.Errorf("unsupported method: %s", method))
	}
	code := 0
	if resp != nil {
		code = resp.StatusCode
	} else {
		errs = append(errs, fmt.Errorf("nodered DELETE no response: %s", c.url(path)))
	}
	return code, []byte(body), errs
}

// GetFlowNodesV1 fetches nodes array of a tab by id via /flow/{id}
func (c *Client) GetFlowNodesV1(ctx context.Context, flowId string, out any) (int, []byte, []error) {
	return c.DoJSON(ctx, "GET", "/flow/"+flowId, nil, out)
}

// GetVersionRevV2 calls /flows with header node-red-api-version:v2 and decodes to out (map)
func (c *Client) GetVersionRevV2(ctx context.Context, out any) (int, []byte, []error) {
	greq := gorequest.New().Retry(c.Retry, time.Second*2).Timeout(c.Timeout)
	url := c.url("/flows")
	logx.WithContext(ctx).Debugf("nodered GET %s (v2)", url)
	resp, body, errs := greq.Get(url).Set("node-red-api-version", "v2").EndStruct(out)
	return resp.StatusCode, []byte(body), errs
}
