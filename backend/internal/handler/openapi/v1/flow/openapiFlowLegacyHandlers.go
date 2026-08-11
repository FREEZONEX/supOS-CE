// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2-1

package flow

import (
	"net/http"

	respx "backend/internal/httpx"
	openapiflow "backend/internal/logic/openapi/v1/flow"
	"backend/internal/svc"
	"backend/internal/types"

	gozerohttpx "github.com/zeromicro/go-zero/rest/httpx"
)

func OpenapiFlowLegacyListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyListReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.List(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func OpenapiFlowLegacyGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyGetReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.Get(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func OpenapiFlowLegacyCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyCreateReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.Create(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func OpenapiFlowLegacyUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyUpdateReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.Update(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func OpenapiFlowLegacyDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyDeleteReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.Delete(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func OpenapiFlowLegacyFlowDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyGetReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.FlowData(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func OpenapiFlowLegacyDeployHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyDeployReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.Deploy(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func OpenapiFlowLegacyNodesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OpenapiFlowLegacyNodesReq
		if !parseOpenapiFlowReq(w, r, &req) {
			return
		}
		l := openapiflow.NewOpenapiFlowLegacyLogic(r.Context(), svcCtx)
		resp, err := l.Nodes(&req)
		writeOpenapiFlowResp(w, resp, err)
	}
}

func parseOpenapiFlowReq(w http.ResponseWriter, r *http.Request, req any) bool {
	if err := gozerohttpx.Parse(r, req); err != nil {
		respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
		return false
	}
	return true
}

func writeOpenapiFlowResp(w http.ResponseWriter, resp *types.Envelope, err error) {
	if err != nil {
		respx.WriteError(w, err)
		return
	}
	gozerohttpx.OkJson(w, resp)
}
