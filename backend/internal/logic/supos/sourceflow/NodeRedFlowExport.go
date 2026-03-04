package sourceflow

import (
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	dao "backend/internal/repo/relationDB"
	"backend/share/base"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"unsafe"

	"github.com/buger/jsonparser"
	"github.com/zeromicro/go-zero/core/logx"
)

func NodeRedFlowExport(ctx context.Context, groupIds []int64, ids []int64, srcFlow bool, apiHost string) func(io.Writer) {
	return func(jsonWriter io.Writer) {
		mapper := dao.NodeRedFlowExporter{}
		var flowIds map[string]bool
		var hasFlowId bool
		if sz := 4*len(groupIds) + len(ids); sz > 0 {
			flowIds = make(map[string]bool, sz)
			hasFlowId = true
		}
		csv2po := func(headers, values []string) *dao.NoderedFlow {
			node := mapper.Csv2Model(headers, values)
			if flowId := node.FlowID; hasFlowId && node.ExportType == "" && len(flowId) > 0 {
				flowIds[flowId] = true
			}
			return node
		}
		flowGetId := func(node *dao.NoderedFlow) int64 {
			return node.ID
		}
		log := logx.WithContext(ctx)
		//sourceFlow: flows, nodes(flowNodes+globalNodes), unsRefs, tags
		//eventFlow: flows, nodes(flowNodes+globalNodes)
		fmt.Fprintln(jsonWriter, `{ "flows":`)
		_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
			return mapper.ExportByGroupAndIds(ctx, groupIds, ids, srcFlow, writer)
		}, jsonWriter, flowGetChildren, flowSetChildren, flowGetId, flowGetParentId, csv2po, true)
		if err != nil {
			log.Error(apiHost, " flow Csv2JsonStream err:", err)
		}
		var supModelIds *[]string
		if srcFlow {
			mids := make([]string, 0, 4+len(ids))
			supModelIds = &mids
		}
		// nodes(flowNodes+globalNodes)
		err = exportNodes(ctx, jsonWriter, apiHost, flowIds, supModelIds)
		if err != nil {
			log.Error("listNodeFlows err:", err)
		}
		if srcFlow {
			fmt.Fprintln(jsonWriter, `, "unsRefs":`)
			var flowUnsMapper dao.NodeRedFlowUnsExporter
			flowUnsCsv2po := func(headers, values []string) *dao.NoderedFlowNode {
				return flowUnsMapper.Csv2Model(headers, values)
			}
			flowUnsGetChildren := func(node *dao.NoderedFlowNode) []*dao.NoderedFlowNode {
				return nil
			}
			flowUnsSetChildren := func(node *dao.NoderedFlowNode, children []*dao.NoderedFlowNode) {
			}
			flowUnsGetId := func(node *dao.NoderedFlowNode) int64 {
				return node.ParentID
			}
			flowUnsGetParentId := func(node *dao.NoderedFlowNode) int64 {
				return -1
			}
			_, err = jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return flowUnsMapper.ExportByFlowIds(ctx, ids, writer)
			}, jsonWriter, flowUnsGetChildren, flowUnsSetChildren, flowUnsGetId, flowUnsGetParentId, flowUnsCsv2po, true)
			if len(*supModelIds) > 0 && err == nil {
				fmt.Fprintln(jsonWriter, `, "tags":[`)
				countOk := 0
				parts := base.Partition(*supModelIds, 80)
				for i, nodeIds := range parts {
					nodeTags, err := listNodeTags(ctx, apiHost, nodeIds)
					log.Infof("导出 nodeTags[%d/%d]: er:%v, ids: %v~%v", i+1, len(parts), err, nodeIds[0], nodeIds[len(nodeIds)-1])
					if nodeTags != nil && err == nil {
						je := json.NewDecoder(nodeTags)
						_, er := je.Token()
						if er == nil {
							for je.More() {
								var raw json.RawMessage
								er = je.Decode(&raw)
								if er == nil {
									if countOk > 0 {
										fmt.Fprintln(jsonWriter, ",")
									}
									countOk++
									jsonWriter.Write(raw)
								}
							}
							_ = nodeTags.Close()
						}
					}
				}
				fmt.Fprintln(jsonWriter, "]")
			}
		}
		fmt.Fprintln(jsonWriter, "}")
	}
}

func exportNodes(ctx context.Context, jsonWriter io.Writer, apiHost string, flowIds map[string]bool, supModelIds *[]string) error {
	allFlowConfigs, err := listNodeFlows(ctx, apiHost)
	if allFlowConfigs != nil {
		defer allFlowConfigs.Close()
	}
	if err != nil {
		return err
	}

	fmt.Fprintln(jsonWriter, `, "nodes":[`)
	if len(flowIds) == 0 && supModelIds == nil { //没有过滤条件，返回全部
		_, err = io.Copy(jsonWriter, allFlowConfigs)
		return err
	}
	decoder := json.NewDecoder(allFlowConfigs)

	// 读取起始数组标记
	if token, err := decoder.Token(); err != nil {
		fmt.Fprintln(jsonWriter, `]`)
		return err
	} else if delim, ok := token.(json.Delim); !ok || delim != '[' {
		fmt.Fprintln(jsonWriter, `]`)
		return fmt.Errorf("json parse error: token=%s", token)
	}
	// 流式读取数组元素
	for i, n := 0, 0; decoder.More(); i++ {
		var raw json.RawMessage
		if err := decoder.Decode(&raw); err != nil {
			if i > 0 {
				fmt.Fprintln(jsonWriter, `]`)
			} else {
				fmt.Fprintln(jsonWriter, `[]`)
			}
			return err
		}
		zBs, _, _, _ := jsonparser.Get(raw, "z")
		idBs, _, _, _ := jsonparser.Get(raw, "id")
		tpBs, _, _, _ := jsonparser.Get(raw, "type")
		z, id, Type := b2s(zBs), b2s(idBs), b2s(tpBs)
		if len(flowIds) == 0 || ((z == "" && Type != "tab") || flowIds[z] || flowIds[id]) {
			if n > 0 {
				jsonWriter.Write([]byte(","))
			}
			n++
			_, err = jsonWriter.Write(raw)
			if supModelIds != nil && "supmodel" == Type {
				*supModelIds = append(*supModelIds, id)
			}
		}
	}
	// 读取结束数组标记
	_, err = decoder.Token()
	fmt.Fprintln(jsonWriter, "\n]")
	return err
}
func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func flowGetChildren(node *dao.NoderedFlow) []*dao.NoderedFlow {
	return node.Children
}
func flowSetChildren(node *dao.NoderedFlow, children []*dao.NoderedFlow) {
	node.Children = children
}

func flowGetParentId(node *dao.NoderedFlow) int64 {
	return base.P2v(node.GroupId)
}

// 查询所有流程
func listNodeFlows(ctx context.Context, apiHost string) (io.ReadCloser, error) {
	url := apiHost + "/flows"
	resp, err := http.Get(url)
	err = logNodeRedErr(ctx, url, resp, err, "查询所有流程失败")
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

// 查询节点引用的位号
func listNodeTags(ctx context.Context, apiHost string, nodeIds []string) (io.ReadCloser, error) {
	urlBuilder := base.StringBuilder{}
	urlBuilder.Grow(128 + 20*len(nodeIds))
	urlBuilder.Append(apiHost).Append("/nodered-api/tags")
	if len(nodeIds) > 0 {
		urlBuilder.Append(`?`)
		for i, nodeId := range nodeIds {
			if i > 0 {
				urlBuilder.Append(`&`)
			}
			urlBuilder.Append(`nodeId=`).Append(nodeId)
		}
	}
	URL := urlBuilder.String()
	resp, err := http.Get(URL)
	err = logNodeRedErr(ctx, apiHost+"/nodered-api/tags", resp, err, "查询节点引用的位号失败")
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func logNodeRedErr(ctx context.Context, url string, resp *http.Response, err error, tag string) error {
	if err != nil {
	} else if resp == nil {
		err = fmt.Errorf("resp is nil, url=%s", url)
	} else if resp.StatusCode/200 != 1 {
		err = fmt.Errorf("code=%d", resp.StatusCode)
	}
	if err != nil {
		logx.WithContext(ctx).Error(tag, " ", err.Error())
	}
	return err
}
