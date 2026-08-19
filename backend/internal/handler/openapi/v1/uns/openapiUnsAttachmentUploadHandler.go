// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"backend/internal/logic/openapi/v1/uns"
	"backend/internal/svc"
	"backend/internal/types"
	"net/http"
	"strconv"
	"strings"

	respx "backend/internal/httpx"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiUnsAttachmentUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UnsAttachmentUploadReq

		contentType := r.Header.Get("Content-Type")
		if !strings.HasPrefix(strings.ToLower(contentType), "multipart/") {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "multipart/form-data is required"))
			return
		}
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid multipart form: "+err.Error()))
			return
		}
		req.Topic = strings.TrimSpace(r.FormValue("topic"))
		req.FileName = strings.TrimSpace(r.FormValue("fileName"))
		req.Sha256 = strings.TrimSpace(r.FormValue("sha256"))
		if rawUnsID := strings.TrimSpace(r.FormValue("unsId")); rawUnsID != "" {
			unsID, err := strconv.ParseInt(rawUnsID, 10, 64)
			if err != nil {
				respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid unsId"))
				return
			}
			req.UnsId = unsID
		}
		if req.Topic == "" && req.UnsId <= 0 {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "topic or unsId is required"))
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "file is required"))
			return
		}
		defer file.Close()

		l := uns.NewOpenapiUnsAttachmentUploadLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiUnsAttachmentUpload(&req, file, header)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		gozerohttpx.OkJson(w, resp)
	}
}
