package importExport

import (
	"backend/internal/logic/supos/uns/importExport"
	"backend/internal/types"
	"io"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ImportHandler(w http.ResponseWriter, r *http.Request) {
	// 删除Keep-Alive头
	w.Header().Del("Keep-Alive")

	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("transfer-encoding", "chunked")

	req, err := types.FormFile(r, "file")
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	importExport.ImportUnsByReader(req.FileName, req.Size, w, req.Reader)
}

// ImportHandler UNS 导入
func ImportHandler2(w http.ResponseWriter, r *http.Request) {
	// 删除Keep-Alive头
	w.Header().Del("Keep-Alive")

	fileSizeStr := r.Header.Get("Content-Length")
	fileSize, _ := strconv.ParseInt(fileSizeStr, 10, 64)

	mr, err := r.MultipartReader()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	for {
		part, err2 := mr.NextPart()
		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err2.Error()))
			return
		}
		if fileName := part.FileName(); fileName != "" {
			if fileSize < 1 {
				fileSizeStr = part.Header.Get("Content-Length")
				fileSize, err = strconv.ParseInt(fileSizeStr, 10, 64)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					_, _ = w.Write([]byte("Missing Header: Content-Length"))
					return
				}
			}
			w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("transfer-encoding", "chunked")
			// 文件部分
			writer, wait := importExport.Import(fileName, fileSize, w)
			n, er := io.CopyN(writer, part, fileSize)
			if er != nil {
				logx.Error("文件上传失败：", er, n, fileName)
				_ = part.Close()
				return
			}
			logx.Debug("文件上传...")
			wait()
			_ = part.Close()
			return
		} else {
			// 字段部分
			_, _ = io.Copy(io.Discard, part)
		}
		_ = part.Close()
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("FileNotFound!"))
}
