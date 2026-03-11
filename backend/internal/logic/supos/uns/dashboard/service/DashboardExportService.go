package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/utils/grafanautil"
	"backend/internal/logic/supos/auth"
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type DashboardExportService struct {
}

func init() {
	spring.RegisterBean(&DashboardExportService{})
}

func (*DashboardExportService) FileName() string {
	return "dashboard.json"
}
func (*DashboardExportService) Order() int {
	return 1000
}

func (*DashboardExportService) ExportStream(ctx context.Context, arg *types.GlobalExportParam) (exporter func(writer io.Writer)) {
	req := arg.DashboardExportParam
	if req == nil {
		return
	}
	log := logx.WithContext(ctx)
	return func(jsonWriter io.Writer) {
		{
			var mapper dao.DashboardMapper
			csv2po := func(headers, values []string) *dao.DashboardModel {
				po := mapper.Csv2Model(headers, values)
				jsonContent, err := grafanautil.Get(ctx, po.ID)
				if err != nil {
					logx.WithContext(ctx).Errorf("get grafana json content err: %v", err)
				} else if len(jsonContent) > 0 {
					po.JsonContent = jsonContent
				}
				return po
			}
			fmt.Fprintln(jsonWriter, `{ "data":`)
			_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return mapper.ExportByGroupAndIds(ctx, req.GroupIds, req.DashIds, writer)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, csv2po, true)
			if err != nil {
				log.Error("Dashboard Csv2JsonStream err:", err)
			}
		}
		{
			var mapper dao.DashboardRefMapper
			csv2po := func(headers, values []string) *dao.DashboardRefModel {
				return mapper.Csv2Model(headers, values)
			}
			fmt.Fprintln(jsonWriter, `, "unsRefs":`)
			_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return mapper.ExportByGroupAndIds(ctx, req.GroupIds, req.DashIds, writer)
			}, jsonWriter, refGetChildren, refSetChildren, refGetId, refGetParentId, csv2po, true)
			if err != nil {
				log.Error("DashboardRef Csv2JsonStream err:", err)
			}

		}
		fmt.Fprintln(jsonWriter, `}`)
	}
}
func nodeGetChildren(node *dao.DashboardModel) []*dao.DashboardModel {
	if node.Type == -1 && node.Children == nil { //group
		return []*dao.DashboardModel{}
	}
	return node.Children
}
func nodeSetChildren(node *dao.DashboardModel, children []*dao.DashboardModel) {
	node.Children = children
}
func nodeGetId(node *dao.DashboardModel) string {
	return node.ID
}
func nodeGetParentId(node *dao.DashboardModel) string {
	gid := base.P2v(node.GroupId)
	if gid < 1 {
		return ""
	}
	return strconv.FormatInt(gid, 10)
}
func refGetChildren(node *dao.DashboardRefModel) []*dao.DashboardRefModel {
	return nil
}
func refSetChildren(node *dao.DashboardRefModel, children []*dao.DashboardRefModel) {
}
func refGetId(node *dao.DashboardRefModel) string {
	return node.DashboardID
}
func refGetParentId(node *dao.DashboardRefModel) string {
	return ""
}
func (*DashboardExportService) ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {
	var dashMapper dao.DashboardMapper
	var groupMapper dao.GroupMapper
	var unsLinkMapper dao.DashboardRefMapper

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
	progress := common.Float3(0)
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
		switch propName {
		case "data":
			importDashboards(ctx, size, &progress, statusConsumer, decoder, groupMapper.Save, dashMapper.Save)
			if progress < 50 && statusConsumer != nil {
				progress = 50
				status := &common.RunningStatus{Task: propName, Progress: &progress}
				statusConsumer(status)
			}
		case "unsRefs":
			importDashUnsLink(ctx, &progress, statusConsumer, decoder, unsLinkMapper.SaveBatch)
		}
	}
	if statusConsumer != nil {
		progress = 100
		status := &common.RunningStatus{
			Code:     200,
			Task:     I18nUtils.GetMessageWithCtx(ctx, "uns.create.task.name.final"),
			Progress: &progress, Finished: base.OptionalTrue}
		statusConsumer(status)
	}
}
func jsonErr(ctx context.Context, err error) error {
	if err != nil {
		if je, is := err.(*json.SyntaxError); is {
			return fmt.Errorf("%s: %d: %v", I18nUtils.GetMessageWithCtx(ctx, "uns.import.json.error"), je.Offset, je.Error())
		}
	}
	return err
}
func importDashboards(ctx context.Context, size int64, progress *common.Float3, statusConsumer func(status *common.RunningStatus), dec *json.Decoder,
	groupSave func(ctx context.Context, list []*dao.GroupModel) error, dashSave func(ctx context.Context, list []*dao.DashboardModel) error) {

	tree2flat := func(c context.Context, propName string, node, parent *dao.DashboardModel) *dao.DashboardModel {
		if parent != nil {
			pid, err := strconv.ParseInt(parent.ID, 10, 64)
			if err == nil {
				node.GroupId = &pid
			}
		}
		return node
	}
	flowType := int16(3)
	v2g := func(v *dao.DashboardModel) *dao.GroupModel {
		id, er := strconv.ParseInt(v.ID, 10, 64)
		if er != nil {
			return nil
		}
		return &dao.GroupModel{
			ID:          id,
			Name:        v.Name,
			Description: v.Description,
			Type:        &flowType,
			Sort:        v.Sort,
			Creator:     v.Creator,
		}
	}

	creatorUser := auth.ResolveUserName(ctx)
	consumer := func(readSize int64, propName string, nodes []*dao.DashboardModel) {
		dashes := make([]*dao.DashboardModel, 0, 512)
		groups := make([]*dao.GroupModel, 0, 16)
		for _, node := range nodes {
			node.Creator = creatorUser
			node.CreateTime = time.Now()
			node.UpdateTime = time.Time{}
			if node.ExportType == "group" {
				group := v2g(node)
				if group != nil {
					groups = append(groups, group)
				}
			} else {
				dashes = append(dashes, node)
			}
		}
		if readSize < size {
			*progress = 0.2 * common.Float3(readSize) / common.Float3(size)
		} else if *progress < 90 {
			*progress += 5
		}
		code, msg := 200, ""
		task := "save"
		if len(groups) > 0 {
			task = "save groups"
			err := groupSave(ctx, groups)
			if err != nil {
				code = 500
				msg = err.Error()
			}
		}
		if len(dashes) > 0 {
			err := dashSave(ctx, dashes)
			if err != nil {
				logx.WithContext(ctx).Error("dash save error :%v", err)
				code = 500
				msg = err.Error()
			} else {
				for _, dash := range dashes {
					jsonContent := dash.JsonContent
					if len(jsonContent) > 0 && jsonContent[0] == '{' {
						_, er := grafanautil.Create(ctx, jsonContent)
						if er != nil {
							err = er
							code = 500
							msg = err.Error()
						}
					} else {
						logx.WithContext(ctx).Errorf("面板json格式错误 : %v, id=%v", jsonContent, dash.ID)
					}
				}

			}
		}
		if statusConsumer != nil {
			statusConsumer(&common.RunningStatus{Progress: progress, Code: code, Msg: msg, Task: task})
		}
	}
	errConsumer := func(node *dao.DashboardModel) {
	}
	_ = jsonstream.DecodeJsonTreeToFlatByJsonDecoder(ctx, dec, 1000, tree2flat, consumer, errConsumer)
}
func importDashUnsLink(ctx context.Context, progress *common.Float3, statusConsumer func(status *common.RunningStatus),
	dec *json.Decoder, saveDashUns func(ctx context.Context, refers []*dao.DashboardRefModel) error) {

	tree2flat := func(c context.Context, propName string, node, parent *dao.DashboardRefModel) *dao.DashboardRefModel {
		node.CreateAt = time.Now()
		return node
	}
	consumer := func(readSize int64, propName string, nodes []*dao.DashboardRefModel) {
		err := saveDashUns(ctx, nodes)
		task := "save dashboard unsLinks"
		code, msg := 200, "ok"
		if err != nil {
			code, msg = 500, err.Error()
		}
		if statusConsumer != nil {
			if *progress < 90 {
				*progress += 1
			}
			statusConsumer(&common.RunningStatus{Progress: progress, Code: code, Msg: msg, Task: task})
		}
	}
	errConsumer := func(node *dao.DashboardRefModel) {
	}
	_ = jsonstream.DecodeJsonTreeToFlatByJsonDecoder(ctx, dec, 1000, tree2flat, consumer, errConsumer)
}
