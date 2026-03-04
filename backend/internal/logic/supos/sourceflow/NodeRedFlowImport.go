package sourceflow

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/logic/supos/auth"
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	dao "backend/internal/repo/relationDB"
	"backend/share/base"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/zeromicro/go-zero/core/logx"
)

func NodeRedImport(ctx context.Context, apiHost string, srcFlow bool, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {
	var flowMapper = dao.NewNoderedSourceFlowRepo(ctx)
	var groupMapper dao.GroupMapper
	var flowUnsMapper dao.NodeFlowModelMapper
	nodeRedImport(ctx, apiHost, srcFlow, reader, statusConsumer, flowMapper.MultiInsert, groupMapper.Save, flowUnsMapper.SaveBatch)
}
func nodeRedImport(ctx context.Context, apiHost string, srcFlow bool, reader io.Reader, statusConsumer func(status *common.RunningStatus),
	saveFlow func(ctx context.Context, data []*dao.NoderedSourceFlow) error,
	saveGroup func(ctx context.Context, groups []*dao.GroupModel) error,
	saveFlowUns func(ctx context.Context, list []*dao.NoderedFlowNode) error,
) {
	decoder := json.NewDecoder(reader)
	_, err := decoder.Token()
	if err != nil {
		if statusConsumer != nil {
			prog := common.Float3(0)
			statusConsumer(&common.RunningStatus{
				Code: 500, Msg: err.Error(),
				Progress: &prog, Finished: base.OptionalTrue})
		}
		return
	}
	progress := common.Float3(0.0)
	// 处理对象中的字段
	var someError error
	for decoder.More() {
		// 读取字段名
		fieldName, er := decoder.Token()
		if er != nil {
			logx.WithContext(ctx).Errorf("错误Token :%v", fieldName)
			er = jsonErr(ctx, er)
			if statusConsumer != nil {
				statusConsumer(&common.RunningStatus{
					Code: 500, Msg: err.Error(),
					Progress: &progress, Finished: base.OptionalTrue})
			}
			return
		}

		propName, isString := fieldName.(string)
		if !isString {
			logx.WithContext(ctx).Errorf("未知Token :%v", fieldName)
			// 跳过未知字段的值
			continue
		}
		var okName = true
		switch propName {
		case "flows":
			err = importFlows(ctx, srcFlow, decoder, saveFlow, saveGroup)
		case "nodes":
			err = importNodes(ctx, apiHost, decoder)
		case "unsRefs":
			err = importFlowUnsLink(ctx, decoder, saveFlowUns)
		case "tags":
			err = importTags(ctx, apiHost, decoder)
		default:
			okName = false
		}
		if statusConsumer != nil && okName {
			if progress <= 75 {
				progress += 25.0
			}
			status := &common.RunningStatus{Task: propName, Progress: &progress}
			if err != nil {
				status.Code = 500
				status.Msg = err.Error()
			} else {
				status.Code = 200
			}
			statusConsumer(status)
		} else if !okName {
			logx.WithContext(ctx).Errorf("未知字段 :%v", propName)
		}
		if err != nil {
			someError = err
		}
	}
	if statusConsumer != nil {
		progress = 100
		status := &common.RunningStatus{
			Code:     200,
			Task:     I18nUtils.GetMessageWithCtx(ctx, "uns.create.task.name.final"),
			Progress: &progress, Finished: base.OptionalTrue}
		if someError != nil {
			status.Code = 500
			status.Msg = someError.Error()
		}
		statusConsumer(status)
	}
}

func importTags(ctx context.Context, apiHost string, decoder *json.Decoder) (err error) {
	writer2reader(func(w io.Writer) {
		w.Write([]byte(`[`))
		token, er := decoder.Token()
		if delim, ok := token.(json.Delim); !ok || delim != '[' {
			logx.WithContext(ctx).Errorf("tags token err:%v, token=%v", er, token)
			err = fmt.Errorf("tags token err:%v, token=%v", er, token)
			w.Write([]byte(`]`))
			return
		}
		for next := false; decoder.More(); {
			var raw json.RawMessage
			decoder.Decode(&raw)
			if next {
				w.Write([]byte(`,`))
			} else {
				next = true
			}
			w.Write(raw)
		}
		token, er = decoder.Token()
		if delim, ok := token.(json.Delim); !ok || delim != ']' {
			logx.WithContext(ctx).Errorf("end tags token err:%v, token=%v", er, token)
			err = fmt.Errorf("end tags token err:%v, token=%v", er, token)
			w.Write([]byte(`]`))
			return
		}
		w.Write([]byte(`]`))
	}, func(r io.Reader) {
		er := saveRefTags(ctx, apiHost, r)
		if er != nil {
			err = er
		}
	})
	return err
}

func importFlowUnsLink(ctx context.Context, dec *json.Decoder, SaveBatch func(ctx context.Context, list []*dao.NoderedFlowNode) error) (err error) {
	node2vo := func(c context.Context, prop string, i, parent *dao.NoderedFlowNode) *dao.NoderedFlowNode {
		return i
	}
	consumer := func(readSize int64, propName string, nodes []*dao.NoderedFlowNode) {
		SaveBatch(ctx, nodes)
	}
	errConsumer := func(node *dao.NoderedFlowNode) {
	}
	err = jsonstream.DecodeJsonTreeToFlatByJsonDecoder(ctx, dec, 1000, node2vo, consumer, errConsumer)
	return
}

func importNodes(ctx context.Context, apiHost string, decoder *json.Decoder) (err error) {
	writer2reader(func(w io.Writer) {
		w.Write([]byte(`[`))
		ids := make(map[string]bool, 1024)
		next := false
		token, er := decoder.Token()
		if delim, ok := token.(json.Delim); !ok || delim != '[' {
			logx.WithContext(ctx).Errorf("Nodes token err:%v, token=%v", er, token)
			w.Write([]byte(`]`))
			err = jsonErr(ctx, fmt.Errorf("nodes token err:%v, token=%v", er, token))
			return
		}
		for decoder.More() {
			var raw json.RawMessage
			err := decoder.Decode(&raw)
			if err != nil {
				logx.WithContext(ctx).Errorf("nodes Decode err:%v", err)
				break
			}
			zBs, _, _, _ := jsonparser.Get(raw, "z")
			idBs, _, _, _ := jsonparser.Get(raw, "id")
			z, id := b2s(zBs), b2s(idBs)
			if len(z) > 0 {
				ids[z] = true
			}
			if len(id) > 0 {
				ids[id] = true
			}
			if next {
				w.Write([]byte(`,`))
			} else {
				next = true
			}
			w.Write(raw)
		}
		token, er = decoder.Token()
		if delim, ok := token.(json.Delim); !ok || delim != ']' {
			logx.WithContext(ctx).Errorf("END Nodes token err:%v, token=%v", er, token)
			err = jsonErr(ctx, fmt.Errorf("end nodes token err:%v, token=%v", er, token))
		}
		existsFlows, er := listNodeFlows(ctx, apiHost)
		if er != nil {
			err = er
		}
		if er == nil {
			je := json.NewDecoder(existsFlows)
			_, er := je.Token()
			if er == nil {
				for je.More() {
					var raw json.RawMessage
					err := je.Decode(&raw)
					if err != nil {
						logx.WithContext(ctx).Errorf("existsFlows Decode err:%v", err)
						break
					}
					zBs, _, _, _ := jsonparser.Get(raw, "z")
					idBs, _, _, _ := jsonparser.Get(raw, "id")
					z, id := b2s(zBs), b2s(idBs)
					if (len(z) > 0 && ids[z]) || (len(id) > 0 && ids[id]) {
						continue
					}
					if next {
						w.Write([]byte(`,`))
					} else {
						next = true
					}
					w.Write(raw)
				}
			}
		}
		w.Write([]byte(`]`))
	}, func(r io.Reader) {
		deployAllToNodeRed(ctx, apiHost, r)
	})
	return err
}

func importFlows(ctx context.Context, srcFlow bool, dec *json.Decoder,
	saveFlow func(ctx context.Context, data []*dao.NoderedSourceFlow) error,
	saveGroup func(ctx context.Context, groups []*dao.GroupModel) error,
) error {
	node2vo := func(c context.Context, prop string, i, parent *dao.NoderedFlow) *dao.NoderedFlow {
		return i
	}

	creatorUser := auth.ResolveUserName(ctx)
	Template := base.SanYuan(srcFlow, dao.TemplateTypeSrcFlow, dao.TemplateTypeEventFlow)
	consumer := func(readSize int64, propName string, nodes []*dao.NoderedFlow) {
		flows := make([]*dao.NoderedFlow, 0, len(nodes))
		groups := make([]*dao.NoderedFlow, 0, 64)
		for _, n := range nodes {
			n.Creator = creatorUser
			n.CreateTime = time.Now()
			n.UpdateTime = time.Time{}
			if n.ExportType == "group" {
				groups = append(groups, n)
			} else {
				n.Template = Template
				flows = append(flows, n)
			}
		}
		if len(flows) > 0 {
			er := saveFlow(ctx, flows)
			if er != nil {
				logx.WithContext(ctx).Error("save flow err:", er)
			}
		}
		if len(groups) > 0 {
			flowType := base.V2p(base.SanYuan(srcFlow, int16(1), int16(2)))
			v2g := func(v *dao.NoderedFlow) *dao.GroupModel {
				return &dao.GroupModel{
					ID:          v.ID,
					Name:        v.FlowName,
					Description: v.Description,
					Type:        flowType,
					Sort:        v.Sort,
					Creator:     v.Creator,
				}
			}
			er := saveGroup(ctx, base.Map(groups, v2g))
			if er != nil {
				logx.WithContext(ctx).Error("save group err:", er)
			}
		}
	}
	errConsumer := func(node *dao.NoderedFlow) {
	}
	err := jsonstream.DecodeJsonTreeToFlatByJsonDecoder(ctx, dec, 1000, node2vo, consumer, errConsumer)
	return err
}
func jsonErr(ctx context.Context, err error) error {
	if err != nil {
		if je, is := err.(*json.SyntaxError); is {
			return fmt.Errorf("%s: %d: %v", I18nUtils.GetMessageWithCtx(ctx, "uns.import.json.error"), je.Offset, je.Error())
		}
	}
	return err
}

// 批量保存引用位号
func saveRefTags(ctx context.Context, apiHost string, body io.Reader) error {
	url := apiHost + "/nodered-api/batchSave/tags"
	resp, err := http.Post(url, "application/json; charset=UTF-8", body)
	err = logNodeRedErr(ctx, url, resp, err, "node-red保存引用位号失败")
	return err
}

// 全局部署
func deployAllToNodeRed(ctx context.Context, apiHost string, nodes io.Reader) error {
	url := apiHost + "/flows"
	resp, err := http.Post(url, "application/json; charset=UTF-8", nodes)
	err = logNodeRedErr(ctx, url, resp, err, "node-red 全局部署 失败")
	return err
}
func writer2reader(w func(writer io.Writer), r func(reader io.Reader)) {
	reader, writer := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer func() {
			_ = writer.Close()
			wg.Done()
		}()
		w(writer)
	}()
	go func() {
		defer func() {
			_ = reader.Close()
			wg.Done()
		}()
		r(reader)
	}()
	wg.Wait()
}
