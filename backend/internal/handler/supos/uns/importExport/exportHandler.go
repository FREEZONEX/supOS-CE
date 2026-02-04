// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package importExport

import (
	"net/http"

	"backend/internal/logic/supos/uns/importExport"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ExportHandler UNS 导出
func ExportHandler(w http.ResponseWriter, r *http.Request) {
	var req types.ExportReq
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	var l importExport.ExportLogic
	resp, err := l.Export(r.Context(), w, &req)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
	} else if resp != nil {
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}

// ExportGlobalHandler 全局导出
func ExportGlobalHandler(w http.ResponseWriter, r *http.Request) {
	var req types.GlobalExportParam
	if err := httpx.Parse(r, &req); err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}
	var l importExport.ExportLogic
	l.ExportGlobal(r.Context(), w, &req)
}
