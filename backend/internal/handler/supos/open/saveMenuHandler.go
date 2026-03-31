// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package open

import (
	"mime/multipart"
	"net/http"
	"regexp"
	"strconv"

	"backend/internal/adapters/kong/dto"
	"backend/internal/adapters/kong/validator"
	"backend/internal/common/errors"
	menu "backend/internal/logic/supos/open"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

var menuNamePattern = regexp.MustCompile(`^[\x20-\x7E]+$`)
var serviceNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_\-$]*$`)

// 保存菜单
func SaveMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
			httpx.ErrorCtx(r.Context(), w, errors.BadRequest(r.Context(), "request.parse.failed"))
			return
		}

		openType, err := strconv.Atoi(r.FormValue("openType"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.BadRequest(r.Context(), "menu.opentype.invalid"))
			return
		}

		if err := validator.ValidateOpenType(r.Context(), openType); err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.NewBuzError(r.Context(), 500, err.Error()))
			return
		}

		file, header, err := r.FormFile("file")
		var icon *multipart.FileHeader
		if err == nil {
			defer file.Close()
			icon = header
		} else if err != http.ErrMissingFile {
			httpx.ErrorCtx(r.Context(), w, errors.NewBuzError(r.Context(), 500, "menu.icon.read.failed"))
			return
		}

		var req types.SaveMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		name := r.FormValue("name")
		if len(name) < 1 || len(name) > 64 || !menuNamePattern.MatchString(name) {
			httpx.ErrorCtx(r.Context(), w, errors.BadRequest(r.Context(), "menu.name.invalid"))
			return
		}
		serviceName := r.FormValue("serviceName")
		if len(serviceName) > 64 || (len(serviceName) > 0 && !serviceNamePattern.MatchString(serviceName)) {
			httpx.ErrorCtx(r.Context(), w, errors.BadRequest(r.Context(), "menu.servicename.invalid"))
			return
		}
		showName, description, baseUrl := r.FormValue("showName"), r.FormValue("description"), r.FormValue("baseUrl")
		if len(showName) < 1 || len(showName) > 64 {
			httpx.ErrorCtx(r.Context(), w, errors.BadRequest(r.Context(), "menu.showname.length"))
			return
		}
		if len(description) > 512 {
			httpx.ErrorCtx(r.Context(), w, errors.BadRequest(r.Context(), "menu.description.length"))
			return
		}
		if len(baseUrl) > 1024 {
			httpx.ErrorCtx(r.Context(), w, errors.BadRequest(r.Context(), "menu.baseurl.length"))
			return
		}
		l := menu.NewSaveMenuLogic(r.Context(), svcCtx)
		resp, err := l.SaveMenu(&dto.MenuDto{
			ServiceName: serviceName,
			Name:        name,
			ShowName:    showName,
			Description: description,
			BaseURL:     baseUrl,
			OpenType:    openType,
			Icon:        icon,
			IsMenu:      true,
		})
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
