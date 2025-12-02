package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/common/utils/fileutil"
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"backend/share/base"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const EXPORT_TYPE_ALL = "ALL"

func (l *UnsImportExportService) ExportPath(ctx context.Context, req *types.ExportReq) (resp *types.ExportResp, err error) {
	resp = &types.ExportResp{}
	resp.Code, resp.Msg = 200, "ok"
	if EXPORT_TYPE_ALL != req.ExportType && len(req.Files)+len(req.Folders) == 0 {
		resp.Msg = "NoArgs"
		return
	}

	countRows := int64(len(req.Files))
	limitSmallFileRows := l.exportConfig.LimitSmallFileRows
	if len(req.Folders) > 0 && countRows < limitSmallFileRows {
		count, er := l.unsMapper.CountChildrenTree(dao.GetDb(ctx), req.Folders)
		if er != nil {
			resp.Code, resp.Msg = 500, "Internal server error"
			err = er
			return
		}
		countRows += count
	}

	resp.Data = &types.ExportPathResult{SmallFile: countRows < limitSmallFileRows}

	jsonBs, _ := json.Marshal(req)
	hash := crc32.ChecksumIEEE(jsonBs)
	relativePath := filepath.Join(constants.ExportRoot, fmt.Sprintf("export%s_%d.json", datetimeutils.DateSimple(), hash))
	targetPath := filepath.Join(fileutil.GetFileRootPath(), relativePath)
	_ = os.MkdirAll(filepath.Dir(targetPath), os.ModeDir)
	paramFile, err := os.Create(targetPath)
	if err != nil {
		resp.Code, resp.Msg = 500, "Internal server error"
		return resp, err
	}
	defer func() {
		err = paramFile.Close()
	}()
	_, err = paramFile.Write(jsonBs)

	resp.Data.FilePath = relativePath
	return
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
func (l *UnsImportExportService) tryExportByParamFile(paramFilePath string, w http.ResponseWriter) bool {
	if filepath.Base(filepath.Dir(paramFilePath)) != filepath.Base(constants.ExportRoot) {
		return false
	}
	l.log.Infof("导出: %s, file=%s\n", paramFilePath, filepath.Base(paramFilePath))
	targetPath := filepath.Join(fileutil.GetFileRootPath(), paramFilePath)
	paramFile, err := os.Open(targetPath)
	if err != nil {
		l.log.Error("open file failed", "err", err, "paramFilePath", paramFilePath)
		return true
	}
	var exportReq types.ExportReq
	err = json.NewDecoder(paramFile).Decode(&exportReq)
	if err != nil {
		l.log.Error("json decode failed", "err", err, "paramFilePath", paramFilePath)
		return true
	}
	// 设置附件下载头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename=UNS_`+filepath.Base(paramFilePath)+``)
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	fmt.Fprintln(w, "{") //开始 JSON 对象
	//if flusher, ok := w.(http.Flusher); ok {
	//	flusher.Flush()
	//}
	l.streamedExportUns(w, &exportReq)
	return true
}

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
