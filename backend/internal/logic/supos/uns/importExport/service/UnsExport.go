package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

const EXPORT_TYPE_ALL = "ALL"

func (l *UnsImportExportService) Export(ctx context.Context, w http.ResponseWriter, req *types.ExportReq) (resp *types.BaseResult, err error) {
	resp = &types.BaseResult{Code: 200, Msg: "ok"}
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
		if countRows == 0 {
			resp.Code, resp.Msg = 204, "NoData"
			return
		} else if countRows < limitSmallFileRows {
			l.doExport(w, datetimeutils.DateSimple()+".json", req, fmt.Sprintf("%d VS %d", countRows, limitSmallFileRows))
			return nil, nil
		} else {
			resp.Msg = fmt.Sprintf("%d VS %d", countRows, limitSmallFileRows)
		}
		return
	} else {
		l.doExport(w, datetimeutils.DateSimple()+".json", req, "")
		return nil, nil
	}
}

func (l *UnsImportExportService) doExport(w http.ResponseWriter, attachmentName string, req *types.ExportReq, msg string) {
	// 设置附件下载头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename=UNS_`+attachmentName)
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	if len(msg) > 0 {
		w.Header().Set("X-Msg", msg)
	}
	l.streamedExportUns(w, req)
}

func label2FileData(lb *dao.UnsLabel) *FileData {
	return &FileData{Name: lb.LabelName}
}
func (l *UnsImportExportService) labelCsv2FileData(headers, values []string) *FileData {
	label := l.labelMapper.Csv2Model(headers, values)
	return label2FileData(label)
}
func (l *UnsImportExportService) unsCsv2FileData(headers, values []string) *FileData {
	uns := l.unsMapper.Csv2Model(headers, values)
	return po2DataVo(uns)
}

// 流式写入json 返回给客户端
func (l *UnsImportExportService) streamedExportUns(out io.Writer, exportReq *types.ExportReq) {
	fmt.Fprintln(out, "{") //开始 JSON 对象
	//if flusher, ok := w.(http.Flusher); ok {
	//	flusher.Flush()
	//}
	jsonWriter := bufio.NewWriter(out)
	if EXPORT_TYPE_ALL == exportReq.ExportType {
		{
			fmt.Fprintf(out, `"%s":`, Label)
			_, err := jsonstream.Csv2JsonStream(l.labelMapper.ExportCsv, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, l.labelCsv2FileData, true)
			if err != nil {
				l.log.Error("Label Csv2JsonStream err:", err)
			}
		}
		{
			fmt.Fprintf(out, `,"%s":`, Template)
			_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return l.unsMapper.ExportCsv([]int16{constants.PathTypeTemplate}, writer)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, l.unsCsv2FileData, true)
			if err != nil {
				l.log.Error("Template Csv2JsonStream err:", err)
			}
		}
		{
			fmt.Fprintf(out, `,"%s":`, UNS)
			_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return l.unsMapper.ExportCsv([]int16{constants.PathTypeDir, constants.PathTypeFile}, writer)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, l.unsCsv2FileData, true)
			if err != nil {
				l.log.Error("UNS Csv2JsonStream err:", err)
			}
		}
	} else if len(exportReq.Folders)+len(exportReq.Files) > 0 {
		fmt.Fprintf(out, `"%s":`, UNS)
		var countUns = 0
		var err error
		if folderIds := exportReq.Folders; len(folderIds) > 0 {
			countUns, err = jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return l.unsMapper.ExportCsvByFolderIds(folderIds, writer)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, l.unsCsv2FileData, false)
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
			if countUns > 0 {
				err = jsonWriter.WriteByte(',')
			}
			countUns, err = jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return l.unsMapper.ExportCsvByIds(idValues, writer)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, l.unsCsv2FileData, false)
		}
		if err != nil {
			l.log.Error("UNS Csv2JsonStream err:", err)
		}
	}
	fmt.Fprintln(out, "}")
}
