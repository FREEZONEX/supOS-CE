package sourceflow

import (
	"backend/internal/common"
	dao "backend/internal/repo/relationDB"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"github.com/h2non/gock"
)

func TestSrcFlowImport(t *testing.T) {
	defer removeMock()
	initMock(t)
	go func() {
		// 每10秒打印一次所有Goroutine的堆栈
		for {
			time.Sleep(10 * time.Second)
			buf := make([]byte, 1024*1024) // 1MB buffer
			n := runtime.Stack(buf, true)  // true 表示打印所有Goroutine
			t.Logf("=== Goroutine Dump ===\n%s\n", string(buf[:n]))
		}
	}()
	//saveRefTags(t.Context(), "http://nodered:1880", bytes.NewBufferString(`{"tags":121}`))
	const API_HOST = "http://nodered:1880"
	//body, err := listNodeFlows(t.Context(), API_HOST)
	//if err != nil {
	//	t.Error(err)
	//} else {
	//	bs, _ := io.ReadAll(body)
	//	t.Log(string(bs))
	//}
	//groupIds := []int64{3, 4}
	//ids := []int64{2004190137830871040, 2004190506053013504, 1963148065022947328, 1910955436975697920}
	//buf := bytes.NewBuffer(make([]byte, 0, 1024))
	//NodeRedFlowExport(t.Context(), groupIds, ids, true, API_HOST)(buf)
	////t.Log(buf.String())
	buf, err := os.Open("F:\\data\\global-flows.json")
	if err != nil {
		t.Error(err)
	}

	saveFlow := func(ctx context.Context, data []*dao.NoderedSourceFlow) error {
		bs, _ := json.Marshal(data)
		t.Logf("数据库 saveFlows: %v", string(bs))
		return nil
	}
	saveGroup := func(ctx context.Context, groups []*dao.GroupModel) error {
		bs, _ := json.Marshal(groups)
		t.Logf("数据库 saveGroups: %v", string(bs))
		return nil
	}
	saveFlowUns := func(ctx context.Context, list []*dao.NoderedFlowNode) error {
		bs, _ := json.Marshal(list)
		t.Logf("数据库 saveFlowUns: %v", string(bs))
		return nil
	}
	removeMock()
	initMock(t)
	// 创建自定义 Transport
	recorder := &recordingTransport{}
	http.DefaultClient.Transport = recorder

	nodeRedImport(t.Context(), API_HOST, true, buf, func(status *common.RunningStatus) {
		bs, _ := json.Marshal(status)
		t.Log("进度：", string(bs))
	}, saveFlow, saveGroup, saveFlowUns)
	time.Sleep(3 * time.Second)
}
func removeMock() {
	gock.Off()
}
func initMock(t *testing.T) {
	dao.InitDbConfig(conf.Database{DSN: "postgres://postgres:postgres@100.100.100.20:31014/postgres?search_path=supos"})

	gock.New("http://nodered:1880").
		Get("/flows").
		Reply(200).
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		BodyString(`[
    {
        "type": "tab",
        "label": "性能测试5",
        "id": "96f6240ba4faf850",
        "info": "",
        "disabled": false
    },
    {
        "type": "tab",
        "label": "性能测试4",
        "id": "62954259cbe020bb",
        "info": "",
        "disabled": false
    },
    {
        "type": "supmodel",
        "label": "超级模型",
        "id": "3c1c52aacbb93a2b",
        "info": "",
        "disabled": false
    },   {
        "endpoint": "opc.tcp://192.168.235.134:53530/OPCUA/SimulationServer",
        "id": "74a9679cba1deda8",
        "login": false,
        "none": true,
        "secmode": "None",
        "secpol": "None",
        "type": "OpcUa-Endpoint",
        "usercert": false,
        "usercertificate": "",
        "userprivatekey": ""
    }
]`)

	gock.New("http://nodered:1880").
		Post("/flows").Reply(200).JSON(map[string]string{"code": "200", "msg": "flow保存成功"})

	gock.New("http://nodered:1880").
		Post("/nodered-api/batchSave/tags").Reply(200).JSON(map[string]string{"code": "200"})

	gock.Observe(func(request *http.Request, mock gock.Mock) {
		if request.Method == "POST" {
			switch request.URL.Path {
			case "/flows":
				t.Log("--保存 flows之前---")
				bs, err := io.ReadAll(request.Body)
				t.Log("保存flows: ", err, string(bs))
				request.Body = io.NopCloser(bytes.NewBuffer(bs))
			case "/nodered-api/batchSave/tags":
				t.Log("--保存tags之前---")
				bs, _ := io.ReadAll(request.Body)
				t.Log("保存tags: ", string(bs))
				request.Body = io.NopCloser(bytes.NewBuffer(bs))
			}

		}
	})

	mockTags := make([]map[string]any, 0, 16)
	mockTags = append(mockTags, map[string]any{
		"12321": []string{
			"tag", "12321",
		},
	})
	gock.New("http://nodered:1880").Get("/nodered-api/tags").Reply(200).JSON(mockTags)
}

// 可记录请求的自定义 Transport
type recordingTransport struct {
	capturedRequests []requestData
	base             http.RoundTripper
}

type requestData struct {
	url          string
	method       string
	body         []byte
	capturedTime time.Time
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// 捕获请求体（如果是 POST/PUT 等）
	var bodyBytes []byte
	if req.Body != nil && (req.Method == "POST" || req.Method == "PUT") {
		// 读取并替换 Body
		bodyBytes, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 记录
		r.capturedRequests = append(r.capturedRequests, requestData{
			url:          req.URL.String(),
			method:       req.Method,
			body:         bodyBytes,
			capturedTime: time.Now(),
		})
	}

	// 使用 gock 的拦截
	if mk, _ := gock.MatchMock(req); mk != nil {
		return gock.DefaultTransport.RoundTrip(req)
	}

	// 如果没有被 gock 拦截，使用基础 Transport
	if r.base == nil {
		r.base = http.DefaultTransport
	}
	return r.base.RoundTrip(req)
}

func TestWithGock(t *testing.T) {
	defer gock.Off() // 测试结束后关闭 mock

	// 设置 mock 响应
	gock.New("http://api.example.com").
		Get("/users/1").
		Reply(200).
		JSON(map[string]string{"name": "John"})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// 发起请求
		resp, err := http.Get("http://api.example.com/users/1")
		if err != nil {
			t.Fatal(err)
		}

		// 验证请求是否被拦截
		if !gock.IsDone() {
			t.Error("Not all requests were intercepted")
		}

		// 处理响应...
		bs := bytes.Buffer{}
		io.Copy(&bs, resp.Body)
		resp.Body.Close()
		t.Log(bs.String())
	}()
	wg.Wait()
}

func Test_writer2reader(t *testing.T) {
	writer2reader(func(writer io.Writer) {
		writer.Write([]byte(`{}`))
	}, func(reader io.Reader) {
		bs, _ := io.ReadAll(reader)
		t.Log("结果：", string(bs))
	})
}
