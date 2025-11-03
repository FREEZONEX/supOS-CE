package handler

import (
	"backend/internal/common/dto"
	"backend/internal/common/errors"
	"backend/internal/common/utils/apiutil"
	"backend/internal/logic/supos/uns/dashboard/model"
	"backend/internal/logic/supos/uns/dashboard/service"
	"net/http"
	"strconv"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// DashboardHandler Dashboard HTTP 处理器
type DashboardHandler struct {
	dashboardLogic *service.DashboardService
}

// NewDashboardHandler 创建 DashboardHandler 实例
func NewDashboardHandler(dashboardLogic *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardLogic: dashboardLogic,
	}
}

// PageListHandler 分页查询 Dashboard
// GET /inter-api/supos/uns/dashboard
func (h *DashboardHandler) PageListHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析查询参数
		keyword := r.URL.Query().Get("k")
		orderCode := r.URL.Query().Get("orderCode")
		isAsc := r.URL.Query().Get("isAsc")

		var typ *int
		if typeStr := r.URL.Query().Get("type"); typeStr != "" {
			typeVal, err := strconv.Atoi(typeStr)
			if err == nil {
				typ = &typeVal
			}
		}

		// 使用已迁移的 PaginationDTO
		var pagination dto.PaginationDTO
		if err := httpx.Parse(r, &pagination); err != nil {
			// 如果解析失败，使用默认值
			pagination.PageNo = 1
			pagination.PageSize = 10
		}

		// 获取用户 ID（从上下文或 session）
		userID := getUserID(r)

		// 调用 Logic 层
		result, err := h.dashboardLogic.PageList(keyword, typ, orderCode, isAsc, pagination.GetPageNo(), pagination.GetPageSize(), userID)
		if err != nil {
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		httpx.OkJson(w, result)
	}
}

// GetDetailHandler 获取 Dashboard 详情
// GET /inter-api/supos/uns/dashboard/detail
func (h *DashboardHandler) GetDetailHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			errors.Fail(w, http.StatusBadRequest, 400, "id is required")
			return
		}

		dashboard, err := h.dashboardLogic.GetById(id)
		if err != nil {
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		errors.SuccessWithData(w, dashboard)
	}
}

// GetByUuidHandler 根据 UID 获取 Dashboard
// GET /inter-api/supos/uns/dashboard/{uid}
func (h *DashboardHandler) GetByUuidHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数获取 uid
		uid := extractPathParam(r.URL.Path, "/inter-api/supos/uns/dashboard/")
		if uid == "" {
			errors.Fail(w, http.StatusBadRequest, 400, "uid is required")
			return
		}

		dbJson, err := h.dashboardLogic.GetByUuid(uid)
		if err != nil {
			if buzErr, ok := err.(*errors.BuzError); ok {
				errors.Fail(w, http.StatusBadRequest, buzErr.Code, buzErr.Msg)
				return
			}
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		errors.SuccessWithData(w, dbJson)
	}
}

// CreateHandler 创建 Dashboard
// POST /inter-api/supos/uns/dashboard
func (h *DashboardHandler) CreateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dashboard model.DashboardModel
		if err := httpx.Parse(r, &dashboard); err != nil {
			errors.Fail(w, http.StatusBadRequest, 400, err.Error())
			return
		}

		// 获取创建者（从上下文或 session）
		creator := getUsername(r) // TODO: 实现获取用户名的逻辑

		created, err := h.dashboardLogic.Create(&dashboard, creator)
		if err != nil {
			if buzErr, ok := err.(*errors.BuzError); ok {
				errors.Fail(w, http.StatusBadRequest, buzErr.Code, buzErr.Msg)
				return
			}
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		httpx.OkJson(w, dto.NewJsonResult(0, "success", created))
	}
}

// EditHandler 编辑 Dashboard
// PUT /inter-api/supos/uns/dashboard
func (h *DashboardHandler) EditHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dashboard model.DashboardModel
		if err := httpx.Parse(r, &dashboard); err != nil {
			errors.Fail(w, http.StatusBadRequest, 400, err.Error())
			return
		}

		err := h.dashboardLogic.Edit(&dashboard)
		if err != nil {
			if buzErr, ok := err.(*errors.BuzError); ok {
				errors.Fail(w, http.StatusBadRequest, buzErr.Code, buzErr.Msg)
				return
			}
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		httpx.OkJson(w, dto.NewJsonResult[any](0, "success", nil))
	}
}

// DeleteHandler 删除 Dashboard
// DELETE /inter-api/supos/uns/dashboard/{uid}
func (h *DashboardHandler) DeleteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数获取 uid
		uid := extractPathParam(r.URL.Path, "/inter-api/supos/uns/dashboard/")
		if uid == "" {
			errors.Fail(w, http.StatusBadRequest, 400, "uid is required")
			return
		}

		err := h.dashboardLogic.Delete(uid)
		if err != nil {
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		httpx.OkJson(w, dto.NewJsonResult[any](0, "success", nil))
	}
}

// MarkTopHandler 置顶 Dashboard
// POST /inter-api/supos/uns/dashboard/mark
func (h *DashboardHandler) MarkTopHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID string `json:"id"`
		}
		if err := httpx.Parse(r, &req); err != nil {
			errors.Fail(w, http.StatusBadRequest, 400, err.Error())
			return
		}

		userID := getUserID(r) // TODO: 实现获取用户 ID 的逻辑

		err := h.dashboardLogic.MarkTop(req.ID, userID)
		if err != nil {
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		httpx.OkJson(w, dto.NewJsonResult[any](0, "ok", nil))
	}
}

// UnmarkTopHandler 取消置顶 Dashboard
// DELETE /inter-api/supos/uns/dashboard/unmark
func (h *DashboardHandler) UnmarkTopHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			errors.Fail(w, http.StatusBadRequest, 400, "id is required")
			return
		}

		userID := getUserID(r) // TODO: 实现获取用户 ID 的逻辑

		err := h.dashboardLogic.RemoveMarkedTop(id, userID)
		if err != nil {
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		httpx.OkJson(w, dto.NewJsonResult[any](0, "ok", nil))
	}
}

// BindUnsHandler 绑定 UNS
// POST /inter-api/supos/uns/dashboard/bindUns
func (h *DashboardHandler) BindUnsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.DashboardRefModel
		if err := httpx.Parse(r, &req); err != nil {
			errors.Fail(w, http.StatusBadRequest, 400, err.Error())
			return
		}

		err := h.dashboardLogic.BindUns(req.DashboardID, req.UnsAlias)
		if err != nil {
			if buzErr, ok := err.(*errors.BuzError); ok {
				errors.Fail(w, http.StatusBadRequest, buzErr.Code, buzErr.Msg)
				return
			}
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		errors.SuccessWithData(w, "ok")
	}
}

// GetByUnsHandler 根据 UNS 获取 Dashboard
// GET /inter-api/supos/uns/dashboard/getByUns
func (h *DashboardHandler) GetByUnsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unsAlias := r.URL.Query().Get("unsAlias")
		if unsAlias == "" {
			errors.Fail(w, http.StatusBadRequest, 400, "unsAlias is required")
			return
		}

		dashboard, err := h.dashboardLogic.GetByUns(unsAlias)
		if err != nil {
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		errors.SuccessWithData(w, dashboard)
	}
}

// CreateGrafanaByUnsHandler 基于 UNS 创建 Grafana Dashboard
// POST /inter-api/supos/uns/dashboard/createGrafanaByUns/{alias}
func (h *DashboardHandler) CreateGrafanaByUnsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从路径参数获取 alias
		alias := extractPathParam(r.URL.Path, "/inter-api/supos/uns/dashboard/createGrafanaByUns/")
		if alias == "" {
			errors.Fail(w, http.StatusBadRequest, 400, "alias is required")
			return
		}

		uuid, err := h.dashboardLogic.CreateGrafanaByUns(alias)
		if err != nil {
			if buzErr, ok := err.(*errors.BuzError); ok {
				errors.Fail(w, http.StatusBadRequest, buzErr.Code, buzErr.Msg)
				return
			}
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		errors.SuccessWithData(w, uuid)
	}
}

// IsExistHandler 检查 Dashboard 是否存在
// GET /inter-api/supos/uns/dashboard/isExist
func (h *DashboardHandler) IsExistHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.URL.Query().Get("alias")
		if alias == "" {
			errors.Fail(w, http.StatusBadRequest, 400, "alias is required")
			return
		}

		dbJson, err := h.dashboardLogic.IsExist(alias)
		if err != nil {
			if buzErr, ok := err.(*errors.BuzError); ok {
				errors.Fail(w, http.StatusBadRequest, buzErr.Code, buzErr.Msg)
				return
			}
			errors.Fail(w, http.StatusInternalServerError, 500, err.Error())
			return
		}

		errors.SuccessWithData(w, dbJson)
	}
}

// 辅助函数

// getUserID 获取当前用户 ID
func getUserID(r *http.Request) string {
	if user := apiutil.GetUserFromContext(r); user != nil {
		return user.Sub
	}
	return "default-user"
}

// getUsername 获取当前用户名
func getUsername(r *http.Request) string {
	if user := apiutil.GetUserFromContext(r); user != nil {
		return user.PreferredUsername
	}
	return "default-user"
}

// extractPathParam 从路径中提取参数
func extractPathParam(path, prefix string) string {
	if len(path) <= len(prefix) {
		return ""
	}
	return path[len(prefix):]
}
