// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package uns

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	respx "backend/internal/httpx"
	"backend/internal/logic/openapi/v1/uns"
	"backend/internal/svc"
	"backend/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiUnsWriteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WriteReq
		if err := decodeOpenapiUnsWriteReq(r, &req); err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}

		l := uns.NewOpenapiUnsWriteLogic(r.Context(), svcCtx)
		resp, err := l.OpenapiUnsWrite(&req)
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		httpx.OkJson(w, resp)
	}
}

func decodeOpenapiUnsWriteReq(r *http.Request, req *types.WriteReq) error {
	if r.Body == nil {
		return fmt.Errorf("body is required")
	}
	defer r.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(r.Body, 10<<20))
	decoder.UseNumber()
	if err := decoder.Decode(req); err != nil {
		return err
	}
	return nil
}
