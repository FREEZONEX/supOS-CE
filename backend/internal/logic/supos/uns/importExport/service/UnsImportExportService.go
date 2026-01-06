package service

import (
	"backend/internal/common/constants"
	"backend/internal/config"
	labelServ "backend/internal/logic/supos/uns/label/service"
	"backend/internal/logic/supos/uns/uns/service"
	dao "backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"
	"backend/share/spring"
	"context"
	"embed"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnsImportExportService struct {
	log           logx.Logger
	unsMapper     dao.UnsNamespaceRepo
	labelMapper   dao.UnsLabelRepo
	unsAddService *service.UnsAddService
	labelService  *labelServ.UnsLabelService
	// 导出认定为小文件的最大数据行数
	exportConfig config.ExportConfig
}

//go:embed templates/*
var templates embed.FS

func init() {
	spring.RegisterLazy[*UnsImportExportService](func() *UnsImportExportService {
		return &UnsImportExportService{
			log:           logx.WithContext(context.Background()),
			unsAddService: spring.GetBean[*service.UnsAddService](),
			labelService:  spring.GetBean[*labelServ.UnsLabelService](),
			exportConfig:  spring.GetBean[*svc.ServiceContext]().Config.Export,
		}
	})
}
func (l *UnsImportExportService) TemplateDownload(req *types.TemplateDownloadReq, r *http.Request, w http.ResponseWriter) error {
	if req.FileType == "json" {
		var httpRequest = *r
		httpRequest.Method = http.MethodGet
		path := constants.JSONTemplatePath
		httpRequest.URL, _ = url.ParseRequestURI(path)
		fileName := filepath.Base(path)
		w.Header().Set("Content-Type", "application/octet-stream;charset=UTF-8")
		w.Header().Set("Content-disposition", "attachment;filename="+fileName)
		http.ServeFileFS(w, &httpRequest, templates, path)
	}
	return nil
}
func (l *UnsImportExportService) FileDownload(req *types.FileDownloadReq, r *http.Request, w http.ResponseWriter) error {
	req.Path = strings.Replace(req.Path, "\\", "/", -1)
	l.log.Info("下载：", req.Path)
	var httpRequest = *r
	httpRequest.Method = http.MethodGet
	path := filepath.Join(constants.RootPath, req.Path)
	path = strings.Replace(path, "\\", "/", -1)
	var err error
	httpRequest.URL, err = url.ParseRequestURI(path)
	if err != nil {
		return err
	}
	fileName := filepath.Base(path)
	w.Header().Set("Content-Type", "application/octet-stream;charset=UTF-8")
	w.Header().Set("Content-disposition", "attachment;filename="+fileName)
	http.ServeFile(w, &httpRequest, path)
	return nil
}

/*
func destFile(fileName string, size int64) (targetPath, relativePath string) {
	relativePath = filepath.Join(constants.Upload, fmt.Sprintf("%s%X_%s", datetimeutils.DateSimple(), size, fileName))
	targetPath = filepath.Join(fileutil.GetFileRootPath(), relativePath)
	return
}
func (l *UnsImportExportService) UploadFile(req *types.MultipartFile) (resp *types.StringResult, err error) {
	extName := filepath.Ext(req.FileName)
	resp = &types.StringResult{}
	resp.Code, resp.Msg = 200, "ok"
	if extName != ".json" {
		resp.Code, resp.Msg = 400, I18nUtils.GetMessage("uns.import.not.json")
		return
	}
	dstPath, relativePath := destFile(req.FileName, req.Size)
	_ = os.MkdirAll(filepath.Dir(dstPath), os.ModeDir)
	dstFile, er := os.Create(dstPath)
	if er != nil {
		err = er
	} else if req.Size > 0 {
		_, err = io.CopyN(dstFile, req.Reader, req.Size)
	} else {
		_, err = io.Copy(dstFile, req.Reader)
	}
	if err == nil {
		resp.Data = relativePath
	}
	return
}

*/
