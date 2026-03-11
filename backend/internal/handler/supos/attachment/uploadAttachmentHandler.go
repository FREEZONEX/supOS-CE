package attachment

import (
	"net/http"

	"backend/internal/logic/supos/attachment"
	"backend/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// UploadAttachmentHandler 上传附件
func UploadAttachmentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 解析multipart表单，支持最大2GB文件
		// 注意：这里使用32MB内存缓存，大文件会写入临时文件
		if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		// 获取上传的文件
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer file.Close()

		// 获取其他表单参数
		category := r.FormValue("category")
		description := r.FormValue("description")
		tags := r.Form["tags"] // 支持多个tag

		l := attachment.NewUploadAttachmentLogic(r.Context(), svcCtx)
		resp, err := l.UploadAttachment(fileHeader, category, description, tags)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
