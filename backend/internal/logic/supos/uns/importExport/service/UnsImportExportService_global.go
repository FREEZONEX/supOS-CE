package service

import (
	"archive/zip"
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/types"
	"backend/share/base"
	"backend/share/spring"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

var exporters []ExportService

type exportServiceSlice []ExportService

func (x exportServiceSlice) Len() int           { return len(x) }
func (x exportServiceSlice) Less(i, j int) bool { return x[i].Order() < x[j].Order() }
func (x exportServiceSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

var exportLock sync.Mutex

func _initExporters() {
	if exporters == nil {
		exportLock.Lock()
		if exporters == nil {
			exporters = spring.GetBeansOfType[ExportService]()
			sort.Sort(exportServiceSlice(exporters))
		}
		exportLock.Unlock()
	}
}
func getExporterByFileName(name string) ExportService {
	i := base.BinarySearch(len(exporters), func(i int) int {
		return cmp.Compare(exporters[i].FileName(), name)
	})
	if i >= 0 {
		return exporters[i]
	} else {
		return nil
	}
}
func (l *UnsImportExportService) FileName() string {
	return "uns.json"
}
func (l *UnsImportExportService) Order() int {
	return 9000
}

func (l *UnsImportExportService) ExportGlobal(ctx context.Context, w http.ResponseWriter, req *types.GlobalExportParam) {
	_initExporters()
	log := logx.WithContext(ctx)
	// 设置附件下载头
	attachFileName := `GlobalUNS_` + datetimeutils.DateSimple() + ".zip"
	log.Info("全局导出：", attachFileName)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename=`+attachFileName)
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	zw := zip.NewWriter(w)
	defer zw.Close()
	zw.SetComment("UNS global export: " + attachFileName)
	for _, exportService := range exporters {
		exporter := exportService.ExportStream(ctx, req)
		if fileName := exportService.FileName(); exporter != nil {
			fw, err := zw.Create(fileName)
			if err == nil {
				t0 := time.Now()
				exporter(fw)
				log.Debugf("导出 %s 耗时 %d ms", fileName, time.Since(t0).Milliseconds())
				_ = zw.Flush()
			} else {
				log.Errorf("创建压缩文件失败: %s, %v", fileName, err.Error())
			}
		}
	}

}
func (l *UnsImportExportService) ImportGlobal(ctx context.Context, zipFileName string, fileSize int64, respWriter io.Writer, readAt io.ReaderAt) {
	_initExporters()
	log := logx.WithContext(ctx)

	r, err := zip.NewReader(readAt, fileSize)
	if err != nil {
		log.Errorf("fail to open zip %v, %s", err, zipFileName)
		writeErrorFinished("openZip", respWriter, err, 100)
		return
	}
	// 遍历所有文件
	TotalTasks := len(r.File)
	for i, f := range r.File {
		fileName := filepath.Base(f.Name)
		exportService := getExporterByFileName(fileName)
		if exportService == nil {
			dir := strings.ToLower(filepath.Dir(f.Name))
			if len(dir) > 0 {
				for _, export := range exporters {
					fName := strings.ToLower(export.FileName())
					if strings.HasPrefix(fName, dir) {
						exportService = export
						break
					}
				}
			}
		}
		module := fileName
		if dot := strings.Index(module, "."); dot > 0 {
			module = module[:dot]
		}
		if moduleMsg := I18nUtils.GetMessage("import." + module); !strings.HasPrefix(moduleMsg, "import.") {
			module = moduleMsg
		}
		if exportService == nil {
			writeError(module, respWriter, fmt.Errorf("unknown file name: %s", fileName))
			continue
		}
		fr, er := f.Open()
		if er != nil {
			writeError(module, respWriter, er)
		} else {
			statusWriter := func(status *common.RunningStatus) {
				status.Module = module
				if status.Progress != nil {
					// 单个任务的进度转为整体进度
					I := common.Float3(i)
					N := common.Float3(TotalTasks)
					*status.Progress = (I + *status.Progress*(I+1)) / N
					if *status.Progress > 100 {
						*status.Progress = 100
					}
				}
				tsJson, _ := json.Marshal(status)
				_, Er := respWriter.Write(append(tsJson, '\n', '\n'))
				respWriter.(http.Flusher).Flush()
				if Er != nil {
					log.Error("导入进度发送失败:", Er, fileName)
				}
			}
			exportService.ImportStream(ctx, fileName, int64(f.UncompressedSize64), fr, statusWriter)
			_ = fr.Close()
		}
	}
}
func writeError(module string, respWriter io.Writer, err error) {
	writeErrorFinished(module, respWriter, err, 0)
}
func writeErrorFinished(module string, respWriter io.Writer, err error, progress common.Float3) {
	status := &common.RunningStatus{Module: module, Code: 500, Msg: err.Error(), Finished: base.OptionalTrue}
	if progress > 0 {
		status.Progress = &progress
	}
	responseWriterStatusConsumer{respWriter: respWriter}.write(status)
}
