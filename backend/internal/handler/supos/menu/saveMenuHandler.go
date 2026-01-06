// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"mime/multipart"
	"net/http"
	"strconv"

	"backend/internal/adapters/kong/dto"
	"backend/internal/adapters/kong/validator"
	"backend/internal/common/errors"
	"backend/internal/logic/supos/menu"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 保存菜单
func SaveMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
			httpx.ErrorCtx(r.Context(), w, errors.NewBuzError(500, "request.parse.failed"))
			return
		}

		openType, err := strconv.Atoi(r.FormValue("openType"))
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.NewBuzError(500, "menu.opentype.invalid"))
			return
		}

		if err := validator.ValidateOpenType(openType); err != nil {
			httpx.ErrorCtx(r.Context(), w, errors.NewBuzError(500, err.Error()))
			return
		}

		file, header, err := r.FormFile("file")
		var icon *multipart.FileHeader
		if err == nil {
			defer file.Close()
			icon = header
		} else if err != http.ErrMissingFile {
			httpx.ErrorCtx(r.Context(), w, errors.NewBuzError(500, "menu.icon.read.failed"))
			return
		}

		var req types.SaveMenuReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := menu.NewSaveMenuLogic(r.Context(), svcCtx)
		resp, err := l.SaveMenu(&dto.MenuDto{
			ServiceName: r.FormValue("serviceName"),
			Name:        r.FormValue("name"),
			ShowName:    r.FormValue("showName"),
			Description: r.FormValue("description"),
			BaseURL:     r.FormValue("baseUrl"),
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
