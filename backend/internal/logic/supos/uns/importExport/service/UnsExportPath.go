package service

import (
	"backend/internal/common/constants"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/common/utils/fileutil"
	"backend/internal/logic/supos/auth"
	dao "backend/internal/repo/relationDB"
	"backend/internal/types"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash/crc32"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (l *UnsImportExportService) ExportPath(ctx context.Context, req *types.ExportReq) (resp *types.ExportResp, err error) {
	resp = &types.ExportResp{}
	resp.Code, resp.Msg = 200, "ok"
	if EXPORT_TYPE_ALL != req.ExportType && len(req.Files)+len(req.Folders) == 0 {
		resp.Code, resp.Msg = 400, "NoArgs"
		return
	}

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
		return
	}
	resp.Data = &types.ExportPathResult{SmallFile: countRows < limitSmallFileRows}

	if EXPORT_TYPE_ALL == req.ExportType {
		relativePath := filepath.Join(constants.ExportRoot, fmt.Sprintf("export%s_%s.json", datetimeutils.DateSimple(), req.Language))
		resp.Data.FilePath = relativePath
	} else {
		jsonBs, _ := json.Marshal(req)
		if len(jsonBs) <= 500 {
			var buf bytes.Buffer
			gz := gzip.NewWriter(&buf)
			_, err = gz.Write(jsonBs)
			_ = gz.Close()
			if err == nil {
				encoded := base64.URLEncoding.EncodeToString(buf.Bytes())
				resp.Data.FilePath = filepath.Join(constants.ExportRoot, encoded+".json")
				return
			}
		}
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
	}
	return
}

func (l *UnsImportExportService) tryExportByParamFile(ctx context.Context, paramFilePath string, w http.ResponseWriter) bool {
	if filepath.Base(filepath.Dir(paramFilePath)) != filepath.Base(constants.ExportRoot) {
		return false
	}
	req, err := l.decodeExportParams(ctx, paramFilePath)
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`{"code":400, "msg":"%s"}`, err.Error())))
		return true
	}
	attachmentName := filepath.Base(paramFilePath)
	if len(attachmentName) > 100 {
		attachmentName = datetimeutils.DateSimple() + ".json"
	}
	l.doExport(w, attachmentName, req)
	return true
}
func (l *UnsImportExportService) decodeExportParams(ctx context.Context, paramFilePath string) (*types.ExportReq, error) {
	baseFile := filepath.Base(paramFilePath)
	var exportReq types.ExportReq
	if !strings.HasPrefix(baseFile, "export") && strings.HasSuffix(baseFile, ".json") {
		l.log.Info("导出 ByUrl:", paramFilePath)

		b64 := baseFile[:len(baseFile)-5]
		bs, er := base64.URLEncoding.DecodeString(b64)
		if er == nil {
			gr, err := gzip.NewReader(bytes.NewReader(bs))
			if err == nil {
				bs, er = io.ReadAll(gr)
			}
		}
		if er != nil {
			return nil, er
		}
		er = json.NewDecoder(bytes.NewReader(bs)).Decode(&exportReq)
		if er != nil {
			l.log.Error("json decode failed:", er, ",", paramFilePath)
			return nil, er
		}
	} else {
		x := strings.Index(baseFile, "_")
		end := strings.LastIndex(baseFile, ".")
		var lang = baseFile[x+1 : end]
		_, hashEr := strconv.ParseInt(lang, 0, 64)
		if hashEr != nil {
			exportReq.UserId = auth.ResolveUserID(ctx)
			exportReq.Language = lang
			exportReq.ExportType = EXPORT_TYPE_ALL
			l.log.Infof("导出全部: %s, language=%s\n", paramFilePath, lang)
		} else {
			l.log.Infof("导出: %s, file=%s\n", paramFilePath, baseFile)
			targetPath := filepath.Join(fileutil.GetFileRootPath(), paramFilePath)
			paramFile, err := os.Open(targetPath)
			if err != nil {
				l.log.Error("open file failed", err, ",", paramFilePath)
				return nil, err
			}

			err = json.NewDecoder(paramFile).Decode(&exportReq)
			if err != nil {
				l.log.Error("json decode failed", err, ",", paramFilePath)
				return nil, err
			}
		}
	}
	return &exportReq, nil
}
