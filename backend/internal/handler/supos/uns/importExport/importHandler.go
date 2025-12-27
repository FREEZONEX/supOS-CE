package importExport

import (
	"backend/internal/logic/supos/uns/importExport"
	"backend/internal/types"
	"net/http"

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
