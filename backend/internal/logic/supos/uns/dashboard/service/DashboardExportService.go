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
		mapper := dao.DashboardMapper{}
		csv2po := func(headers, values []string) *dao.DashboardModel {
			return mapper.Csv2Model(headers, values)
		}
		fmt.Fprintln(jsonWriter, `{ "data":`)
		_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
			return mapper.ExportByGroupAndIds(ctx, req.GroupIds, req.DashIds, writer)
		}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, csv2po, true)
		if err == nil {
			fmt.Fprintln(jsonWriter, `}`)
		} else {
			log.Error("Dashboard Csv2JsonStream err:", err)
		}

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

func (*DashboardExportService) ImportStream(ctx context.Context, fileName string, size int64, reader io.Reader, statusConsumer func(status *common.RunningStatus)) {
	var dashMapper dao.DashboardMapper
	var groupMapper dao.GroupMapper

	importDashboards(ctx, size, statusConsumer, reader, groupMapper.Save, dashMapper.Save)
}

func importDashboards(ctx context.Context, size int64, statusConsumer func(status *common.RunningStatus), reader io.Reader,
	groupSave func(ctx context.Context, list []*dao.GroupModel) error, dashSave func(ctx context.Context, list []*dao.DashboardModel) error) {

	tree2flat := func(propName string, node, parent *dao.DashboardModel) *dao.DashboardModel {
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

	progress := common.Float3(0.0)
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
			progress = 0.2 * common.Float3(readSize) / common.Float3(size)
		} else if progress < 90 {
			progress += 5
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
					}
				}

			}
		}
		if statusConsumer != nil {
			statusConsumer(&common.RunningStatus{Progress: &progress, Code: code, Msg: msg, Task: task})
		}
	}
	errConsumer := func(node *dao.DashboardModel) {
	}
	err := jsonstream.DecodeJsonTreeToFlat(reader, 1000, tree2flat, consumer, errConsumer)
	if statusConsumer != nil {
		code, msg := 200, ""
		if err != nil {
			code = 500
			msg = err.Error()
		}
		progress = 100
		statusConsumer(&common.RunningStatus{
			Code: code, Msg: msg,
			Task:     I18nUtils.GetMessage("uns.create.task.name.final"),
			Progress: &progress, Finished: base.OptionalTrue})
	}
}
