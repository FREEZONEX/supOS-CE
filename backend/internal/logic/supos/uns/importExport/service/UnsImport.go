package service

import (
	"backend/internal/common"
	"backend/internal/common/I18nUtils"
	"backend/internal/common/constants"
	"backend/internal/common/utils/datetimeutils"
	"backend/internal/common/utils/fileutil"
	"backend/internal/common/utils/integerutil"
	"backend/internal/logic/supos/uns/importExport/service/jsonstream"
	"backend/internal/logic/supos/uns/uns/bo"
	"backend/internal/types"
	"backend/share/base"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

func (l *UnsImportExportService) ImportUns(fileName string, fileSize int64, respWriter io.Writer) (w io.Writer, waiter func()) {

	l.log.Infof("UNS导入: %s (size=%d)\n", fileName, fileSize)
	//
	pipeReader, pipeWriter := io.Pipe()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		//
		l.doImport(fileName, fileSize, pipeReader, respWriter)
		wg.Done()
		_, _ = io.Copy(io.Discard, pipeReader) //防止坏的文件或其他意外让json没有读完导致外面的 w.Write()卡死
	}()
	return &noBufferWriter{tar: pipeWriter}, func() {
		wg.Wait()
	}
}
func (l *UnsImportExportService) ImportUnsDirect(fileName string, fileSize int64, respWriter io.Writer, reader io.Reader) {
	l.log.Infof("UNS导入ByReader: %s (size=%d)\n", fileName, fileSize)
	l.doImport(fileName, fileSize, reader, respWriter)
}

type noBufferWriter struct {
	tar io.Writer
}

func (fr *noBufferWriter) Write(p []byte) (int, error) {
	return fr.tar.Write(p)
}
func (fr *noBufferWriter) ReadFrom(r io.Reader) (n int64, err error) {
	return io.Copy(struct{ io.Writer }{fr.tar}, r)
}

func (l *UnsImportExportService) doImport(fileName string, fileSize int64, pipeReader io.Reader, respWriter io.Writer) {
	var errFile *os.File
	var errBufWriter *bufio.Writer
	errFileRelativePath := ""
	var errJsonEncoder *json.Encoder

	pushStatus := func(status *common.RunningStatus) {
		task := status.Task
		segments := strings.Split(task, ".")
		if sz := len(segments); sz > 1 {
			status.Task = I18nUtils.GetMessage(segments[sz-2])
		} else if sz == 1 {
			status.Task = I18nUtils.GetMessage(status.Task)
		}
		tsJson, _ := json.Marshal(status)
		_, er := respWriter.Write(append(tsJson, '\n', '\n'))
		respWriter.(http.Flusher).Flush()
		if er != nil {
			l.log.Error("导入进度发送失败:", er)
		}
	}
	createErrorFile := func() bool {
		if errFile != nil {
			return false
		}
		var tarPath string
		var err error
		errFileRelativePath = filepath.Join(constants.ImportErr, fmt.Sprintf("err_%s_%s", datetimeutils.DateSimple(), fileName))
		tarPath = filepath.Join(fileutil.GetFileRootPath(), errFileRelativePath)

		_ = os.MkdirAll(filepath.Dir(tarPath), os.ModeDir)
		errFile, err = os.Create(tarPath)
		if err != nil {
			l.log.Error("创建错误提示文件失败", err, tarPath)
		} else {
			errBufWriter = bufio.NewWriter(errFile)
			_ = errBufWriter.WriteByte('[')
			errJsonEncoder = json.NewEncoder(errBufWriter)
		}
		return errFile != nil
	}

	FILE_SIZE := fileSize
	var TOTAL_SIZE = float64(FILE_SIZE)
	var prevReadSize int64 = 0
	var progress float64 = 0
	prevTask := ""

	countUns, countErr := 0, 0
	er := jsonstream.DecodeJsonTreeToFlat(pipeReader, l.exportConfig.BatchSize, node2vo, func(readSize int64, propName string, nodes []*types.CreateTopicDto) {
		if prevReadSize < readSize {
			newProgress := 20 * float64(readSize) / TOTAL_SIZE
			if newProgress <= progress {
				if progress < 80 {
					if propName == prevTask {
						progress += 1
					} else {
						if progress < 60 {
							progress += 10
						} else {
							progress += 1
						}
						prevTask = propName
					}
				} else if progress < 90 {
					progress += 0.01
				}
			} else if newProgress < 90 {
				progress = newProgress
			}
			status := &common.RunningStatus{Code: 200, Task: propName}
			status.SetProgress(progress)
			pushStatus(status)
			prevReadSize = readSize
		}
		switch propName {
		case Label:
			labelNames := base.Map[*types.CreateTopicDto, string](nodes, func(e *types.CreateTopicDto) string {
				return e.Name
			})
			_, er := l.labelService.CreateBatch(context.Background(), labelNames)
			if er != nil {
				l.log.Error("创建标签失败", er)
			}
		case Template, UNS:
			if propName == UNS {
				var insertTemplates []*types.CreateTopicDto
				for _, n := range nodes {
					if template := n.Template; template != nil {
						insertTemplates = append(insertTemplates, template)
					}
				}
				if len(insertTemplates) > 0 {
					nodes = append(insertTemplates, nodes...)
				}
			}
			countUns += len(nodes)
			errTipMap := l.unsAddService.CreateModelAndInstancesInner(context.Background(), bo.CreateModelInstancesArgs{
				Topics:     nodes,
				FromImport: true,
				StatusConsumer: func(status *common.RunningStatus) {
					if progress < 95 {
						if status.Code > 0 {
							if status.N != nil && *status.N > 1 {
								progress += 1 / float64(*status.N)
							} else {
								progress += 0.1
							}
						}
						progressStatus := &common.RunningStatus{Code: 200, Msg: status.Msg, Task: status.Task, SpendMills: status.SpendMills}
						progressStatus.SetProgress(progress)
						pushStatus(progressStatus)
					}
				},
			})
			if len(errTipMap) > 0 {
				countErr += len(errTipMap)
				first := errFile == nil
				createErrorFile()
				logErrImports(errTipMap, nodes, first, errBufWriter, errJsonEncoder)
			}
		}
	}, func(errNode *FileData) {
		first := errFile == nil
		createErrorFile()
		if !first {
			_ = errBufWriter.WriteByte(',')
		}
		_ = errJsonEncoder.Encode(errNode)
	})
	if er != nil {
		l.log.Error("JsonDecodeError", er)
		first := errFile == nil
		createErrorFile()
		if !first {
			_ = errBufWriter.WriteByte(',')
		}
		_ = errJsonEncoder.Encode(er.Error())
	}
	status := &common.RunningStatus{Code: 200, Msg: I18nUtils.GetMessage("uns.import.rs.ok"), Task: I18nUtils.GetMessage("uns.create.task.name.final")}
	status.SetProgress(100)
	status.Finished = base.OptionalTrue
	if errFile != nil {
		_ = errBufWriter.WriteByte(']')
		_ = errBufWriter.Flush()

		er = errFile.Close()
		status.Code = 206
		if countUns == countErr {
			status.Msg = I18nUtils.GetMessage("global.import.rs.allErr")
		} else {
			status.Msg = I18nUtils.GetMessage("uns.import.rs.hasErr") + fmt.Sprintf(": %d/%d", countErr, countUns)
		}
		status.ErrTipFile = errFileRelativePath
	} else if countUns == 0 {
		status.Msg = I18nUtils.GetMessage("uns.noData")
	}
	pushStatus(status)
}

func logErrImports(errTipMap map[string]string, nodes []*types.CreateTopicDto, first bool, errBufWriter *bufio.Writer, errJsonEncoder *json.Encoder) {
	var indexMap = make(map[int]*FileData, len(errTipMap))
	for k, v := range errTipMap {
		if n, er := integerutil.ExtractTailNumbers(k); er == nil {
			i := int(n)
			fileData := uns2DataVo(nil, nodes[i])
			fileData.Error = v
			indexMap[i] = fileData
		}
	}
	indexes := base.MapKeys(indexMap)
	sort.Ints(indexes)
	for _, index := range indexes {
		fileData := indexMap[index]
		if !first {
			_ = errBufWriter.WriteByte(',')
		} else {
			first = false
		}
		_ = errJsonEncoder.Encode(fileData)
	}
	_ = errBufWriter.Flush()
}
