package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"cmp"
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const EXPORT_TYPE_ALL = "ALL"

func (l *UnsImportExportService) Export(ctx context.Context, w http.ResponseWriter, req *types.ExportReq) (resp *types.JsonResult, err error) {
	resp = &types.JsonResult{Code: 200, Msg: "ok"}
	if EXPORT_TYPE_ALL != req.ExportType && len(req.Files)+len(req.Folders) == 0 {
		resp.Code, resp.Msg = 400, "NoArgs"
		return
	}
	if base.P2v(req.CheckSmallFile) {
		countRows := int64(len(req.Files))
		limitSmallFileRows := l.exportConfig.LimitSmallFileRows
		if EXPORT_TYPE_ALL == req.ExportType {
			count, er := l.unsMapper.CountAll(dao.GetDb(ctx))
			if er != nil {
				resp.Code, resp.Msg = 500, er.Error()
				err = er
				return
			}
			countRows = count
		} else if len(req.Folders) > 0 && countRows < limitSmallFileRows {
			count, er := l.unsMapper.CountChildrenTree(dao.GetDb(ctx), req.Folders)
			if er != nil {
				resp.Code, resp.Msg = 500, er.Error()
				err = er
				return
			}
			countRows += count
		}
		resp.Msg = fmt.Sprintf("%d VS %d", countRows, limitSmallFileRows)
		if countRows == 0 {
			resp.Code, resp.Msg = 200, "NoData"
		} else {
			resp.Data = countRows < limitSmallFileRows
		}
		return
	} else {
		l.doExport(w, datetimeutils.DateSimple()+".json", req)
		return nil, nil
	}
}

func (l *UnsImportExportService) doExport(w http.ResponseWriter, attachmentName string, req *types.ExportReq) {
	// 设置附件下载头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename=UNS_`+attachmentName)
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	fmt.Fprintln(w, "{") //开始 JSON 对象
	//if flusher, ok := w.(http.Flusher); ok {
	//	flusher.Flush()
	//}
	l.streamedExportUns(w, req)
}

func labelGetId(lb *dao.UnsLabel) int64 {
	return lb.ID
}
func labelGetParentId(lb *dao.UnsLabel) int64 {
	return -1
}
func label2FileData(lb *dao.UnsLabel) *FileData {
	return &FileData{Name: lb.LabelName}
}

// 流式写入json 返回给客户端
func (l *UnsImportExportService) streamedExportUns(out io.Writer, exportReq *types.ExportReq) {
	if EXPORT_TYPE_ALL == exportReq.ExportType {
		hasData := false
		{
			encoder := jsonstream.NewStreamJsonEncoder(out, nodeGetChildren, nodeSetChildren, labelGetId, labelGetParentId, label2FileData)
			hasData = queryDbAndSendJson(out, encoder, Label, func(page, pageSize int) ([]*dao.UnsLabel, error) {
				return l.labelMapper.ListAll(dao.GetDb(context.Background()), page, pageSize)
			})
		}
		{
			if hasData {
				fmt.Fprint(out, ",")
			}
			encoder := jsonstream.NewStreamJsonEncoder(out, nodeGetChildren, nodeSetChildren, poGetId, poGetParentId, po2DataVo)
			hasTemplate := queryDbAndSendJson(out, encoder, Template, func(page, pageSize int) ([]*dao.UnsNamespace, error) {
				return l.unsMapper.ListAll(dao.GetDb(context.Background()), []int16{constants.PathTypeTemplate}, page, pageSize)
			})
			if !hasTemplate {
				fmt.Fprintf(out, `"%s": []`, Template)
			} else {
				hasData = true
			}
		}
		{
			if hasData {
				fmt.Fprint(out, ",")
			}
			encoder := jsonstream.NewStreamJsonEncoder(out, nodeGetChildren, nodeSetChildren, poGetId, poGetParentId, po2DataVo)
			hasUns := queryDbAndSendJson(out, encoder, UNS, func(page, pageSize int) ([]*dao.UnsNamespace, error) {
				return l.unsMapper.ListAll(dao.GetDb(context.Background()), []int16{constants.PathTypeDir, constants.PathTypeFile}, page, pageSize)
			})
			if !hasUns {
				fmt.Fprintf(out, `"%s": []`, UNS)
			} else {
				hasData = true
			}
		}
	} else if len(exportReq.Folders)+len(exportReq.Files) > 0 {
		fmt.Fprintf(out, `"%s":`, UNS)
		encoder := jsonstream.NewStreamJsonEncoder(out, nodeGetChildren, nodeSetChildren, poGetId, poGetParentId, po2DataVo)
		hashDir, hasFile := false, false
		if len(exportReq.Folders) > 0 {
			layRecPrev := base.Map[int64, string](exportReq.Folders, func(e int64) string {
				return strconv.FormatInt(e, 10)
			})
			hashDir = queryDbAndSendJson(out, encoder, "", func(page, pageSize int) ([]*dao.UnsNamespace, error) {
				return l.unsMapper.PageListByLayRecs(dao.GetDb(context.Background()), layRecPrev, page, pageSize)
			})
		}
		if len(exportReq.Files) > 0 {
			layRecs, _ := l.unsMapper.ListLayRecByIds(dao.GetDb(context.Background()), exportReq.Files)
			ids := make(map[int64]bool, len(layRecs))
			for _, layerRec := range layRecs {
				parts := strings.Split(layerRec, "/")
				for _, part := range parts {
					num, err := strconv.ParseInt(part, 10, 64)
					if err == nil {
						ids[num] = true
					}
				}
			}
			idValues := base.MapKeys(ids)
			sort.Sort(base.LongSlice(idValues))
			hasFile = queryDbAndSendJson(out, encoder, "", func(page, pageSize int) ([]*dao.UnsNamespace, error) {
				return l.unsMapper.PageListByIds(dao.GetDb(context.Background()), idValues, page, pageSize)
			})
		}
		if hashDir || hasFile {
			_ = jsonstream.WriteBatch(encoder, nil, true)
		}
	}
	fmt.Fprintln(out, "}")
}
func queryDbAndSendJson[Node any, ID cmp.Ordered, TreeNode any](
	out io.Writer,
	encoder *jsonstream.StreamJsonEncoder[Node, ID, TreeNode],
	propName string,
	pageQuery func(page, pageSize int) ([]*TreeNode, error)) bool {
	page, pageSize := 1, 1000
	for {
		list, er := pageQuery(page, pageSize)
		if er != nil || len(list) == 0 {
			break
		}
		if page == 1 && propName != "" {
			fmt.Fprintf(out, `"%s":`, propName)
		}
		page++
		_ = jsonstream.WriteBatch(encoder, list, false)
	}
	if page > 1 && propName != "" {
		_ = jsonstream.WriteBatch(encoder, nil, true)
	}
	return page > 1
}
