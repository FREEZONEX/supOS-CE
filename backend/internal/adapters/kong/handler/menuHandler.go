package handler

import (
	"backend/internal/adapters/kong/dto"
	"backend/internal/adapters/kong/logic"
	"backend/internal/adapters/kong/validator"
	"backend/internal/adapters/kong/vo"
	"backend/internal/common/errors"
	"net/http"
	"strconv"

	"mime/multipart"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type MenuHandler struct {
	menuLogic *logic.MenuLogic
}

func NewMenuHandler(menuLogic *logic.MenuLogic) *MenuHandler {
	return &MenuHandler{menuLogic: menuLogic}
}

// SaveMenuHandler 对应 POST /open-api/menu
func (h *MenuHandler) SaveMenuHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
			httpx.Error(w, errors.NewBuzError(500, "request.parse.failed"))
			return
		}

		openType, err := strconv.Atoi(r.FormValue("openType"))
		if err != nil {
			httpx.Error(w, errors.NewBuzError(500, "menu.opentype.invalid"))
			return
		}

		if err := validator.ValidateOpenType(openType); err != nil {
			httpx.Error(w, errors.NewBuzError(500, err.Error()))
			return
		}

		file, header, err := r.FormFile("file") // 注意：Java 中参数名是 "file"
		var icon *multipart.FileHeader
		if err == nil {
			defer file.Close()
			icon = header
		} else if err != http.ErrMissingFile {
			httpx.Error(w, errors.NewBuzError(500, "menu.icon.read.failed"))
			return
		}

		menuDto := &dto.MenuDto{
			ServiceName: r.FormValue("serviceName"),
			Name:        r.FormValue("name"),
			ShowName:    r.FormValue("showName"),
			Description: r.FormValue("description"),
			BaseURL:     r.FormValue("baseUrl"),
			OpenType:    openType,
			IsMenu:      true,
			Icon:        icon,
		}

		if err := h.menuLogic.CreateMenu(menuDto, true); err != nil {
			httpx.Error(w, err)
			return
		}

		httpx.OkJson(w, vo.Success[any](nil))
	}
}
