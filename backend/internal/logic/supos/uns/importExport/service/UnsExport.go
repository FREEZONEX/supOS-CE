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
	return &FileData{Name: values[0]}
}

// 流式写入json 返回给客户端
func (l *UnsImportExportService) streamedExportUns(out io.Writer, exportReq *types.ExportReq) {
	fmt.Fprintln(out, "{") //开始 JSON 对象
	//if flusher, ok := w.(http.Flusher); ok {
	//	flusher.Flush()
	//}
	exportCtx := &exportContext{unsWrapTemplate: EXPORT_TYPE_ALL != exportReq.ExportType}
	unsCsv2FileData := func(headers, values []string) *FileData {
		uns := l.unsMapper.Csv2Model(headers, values)
		return uns2DataVo(exportCtx, uns)
	}
	jsonWriter := bufio.NewWriter(out)
	if EXPORT_TYPE_ALL == exportReq.ExportType {
		{
			fmt.Fprintf(out, `"%s":`, Label)
			_, err := jsonstream.Csv2JsonStream(l.labelMapper.ExportCsv, jsonWriter, nodeGetChildren, nodeSetChildren,
				func(f *FileData) int64 {
					return 0
				}, func(f *FileData) int64 {
					return -1
				}, l.labelCsv2FileData, true)
			if err != nil {
				l.log.Error("Label Csv2JsonStream err:", err)
			}
		}
		{
			fmt.Fprintf(out, `,"%s":`, Template)
			_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return l.unsMapper.ExportCsv(context.Background(), []int16{constants.PathTypeTemplate}, writer)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, unsCsv2FileData, true)
			if err != nil {
				l.log.Error("Template Csv2JsonStream err:", err)
			}
		}
		{
			fmt.Fprintf(out, `,"%s":`, UNS)
			_, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return l.unsMapper.ExportCsv(context.Background(), []int16{constants.PathTypeDir, constants.PathTypeFile}, writer)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, unsCsv2FileData, true)
			if err != nil {
				l.log.Error("UNS Csv2JsonStream err:", err)
			}
		}
	} else if dirIds, fileIds := exportReq.Folders, exportReq.Files; len(dirIds) > 0 || len(fileIds) > 0 {
		fmt.Fprintf(out, `"%s":[`, UNS)
		dirLayRecs, ids := l.getLayAndIds(dirIds, fileIds)
		if len(dirLayRecs) > 0 || len(ids) > 0 {
			countUns, err := jsonstream.Csv2JsonStream(func(writer io.Writer) error {
				return l.unsMapper.ExportCsvByLayRecAndIds(context.Background(), dirLayRecs, ids, writer, true)
			}, jsonWriter, nodeGetChildren, nodeSetChildren, nodeGetId, nodeGetParentId, unsCsv2FileData, false)
			l.log.Info("UNS Csv2JsonStream:", err, countUns)
		}
		fmt.Fprintln(out, "]")
	}
	fmt.Fprintln(out, "}")
}
func (l *UnsImportExportService) getLayAndIds(dirIds, fileIds []int64) (layRec []string, ids []int64) {
	layRecs, err := l.unsMapper.ListLayRecByIds(dao.GetDb(context.Background()), append(dirIds, fileIds...))
	if len(layRecs) == 0 {
		l.log.Error("UNS Mapper ListLayRecByIds:", err)
		return
	}
	return getLayAndIdsInner(dirIds, fileIds, layRecs)
}

func getLayAndIdsInner(dirIds, fileIds []int64, layRecs []string) (layRec []string, ids []int64) {

	dirMap := make(map[int64]string, len(dirIds))
	for _, id := range dirIds {
		dirMap[id] = ""
	}
	fMap := make(map[int64]string, len(fileIds))

	dirIdMap := make(map[int64]int, len(dirIds))
	fileIdMap := make(map[int64]int, len(fileIds))
	for _, layerRec := range layRecs {
		parts := strings.Split(layerRec, "/")
		var idMap map[int64]int
		{
			id, _ := strconv.ParseInt(parts[len(parts)-1], 10, 64)
			if _, has := dirMap[id]; has {
				idMap = dirIdMap
				dirMap[id] = layerRec
			} else {
				idMap = fileIdMap
				fMap[id] = layerRec
			}
			idMap[id] += 1
		}
		for i := len(parts) - 2; i >= 0; i-- {
			id, _ := strconv.ParseInt(parts[i], 10, 64)
			idMap[id] += 1
		}
	}
	dirLayRecs := base.Filter(base.MapValues(dirMap), func(e string) bool {
		return len(e) > 0
	})
	sort.Strings(dirLayRecs)

	for id := range fileIdMap {
		if dirIdMap[id] > 0 {
			delete(fileIdMap, id)
		}
	}
	if len(fMap) > 0 && len(dirLayRecs) > 0 {
		for _, fileLayRec := range fMap {
			i := base.BinarySearchArray(dirLayRecs, fileLayRec, func(a, b string) int {
				if strings.HasPrefix(b, a) {
					return 0
				} else {
					return strings.Compare(a, b)
				}
			})
			if i >= 0 {
				parts := strings.Split(fileLayRec, "/")
				for _, part := range parts {
					id, _ := strconv.ParseInt(part, 10, 64)
					countUsed := fileIdMap[id]
					if countUsed > 0 {
						fileIdMap[id] -= 1
						if countUsed == 1 {
							delete(fileIdMap, id)
						}
					}
				}
			}
		}
	}

	return dirLayRecs, base.MapKeys(fileIdMap)
}
