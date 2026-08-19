package handlerx

import (
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	domainasset "backend/internal/domain/asset"
	respx "backend/internal/httpx"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
)

type AssetMultipartOptions struct {
	OwnerType string
	OwnerID   int64
}

func IsMultipart(r *http.Request) bool {
	return strings.HasPrefix(strings.ToLower(r.Header.Get("Content-Type")), "multipart/form-data")
}

func UploadMultipartAssets(r *http.Request, svcCtx *svc.ServiceContext, opts AssetMultipartOptions) (map[string]any, error) {
	if err := r.ParseMultipartForm(64 << 20); err != nil {
		return nil, respx.NewHTTPError(http.StatusBadRequest, "invalid multipart request: "+err.Error())
	}
	defer r.MultipartForm.RemoveAll()

	files := r.MultipartForm.File["files"]
	if len(files) == 0 {
		files = r.MultipartForm.File["file"]
	}
	if len(files) == 0 {
		return nil, respx.NewHTTPError(http.StatusBadRequest, "file is required")
	}

	ownerType := strings.TrimSpace(opts.OwnerType)
	if ownerType == "" {
		ownerType = strings.TrimSpace(r.FormValue("ownerType"))
	}
	ownerID := opts.OwnerID
	if ownerID == 0 {
		ownerID = domainasset.ParseOwnerID(r.FormValue("ownerId"))
	}

	list := make([]map[string]any, 0, len(files))
	for _, header := range files {
		file, err := header.Open()
		if err != nil {
			return nil, err
		}
		contentType := header.Header.Get("Content-Type")
		if contentType == "" {
			contentType = mime.TypeByExtension(strings.ToLower(filepath.Ext(header.Filename)))
		}
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		data, err := svcCtx.App.Asset.UploadContent(r.Context(), domainasset.UploadContentCommand{
			FileName:    header.Filename,
			ContentType: contentType,
			Size:        header.Size,
			Reader:      file,
			UserID:      logicx.UserID(r.Context()),
		})
		_ = file.Close()
		if err != nil {
			return nil, err
		}
		if ownerType != "" && ownerID != 0 {
			bound, err := svcCtx.App.Asset.Bind(r.Context(), domainasset.BindCommand{
				FileID:    domainasset.ParseOwnerID(assetValueToString(data["fileId"])),
				OwnerType: ownerType,
				OwnerID:   ownerID,
				UserID:    logicx.UserID(r.Context()),
			})
			if err != nil {
				return nil, err
			}
			if binding, ok := bound["binding"]; ok {
				data["binding"] = binding
			}
		}
		list = append(list, data)
	}
	return map[string]any{"list": list, "total": len(list)}, nil
}

func assetValueToString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int64:
		return strconv.FormatInt(v, 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case float64:
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
	}
}
