package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/share/spring"
	"io"
	"strings"
)

func Import(fileName string, fileSize int64, respWriter io.Writer, reader io.Reader) {
	unsImportService := spring.GetBean[*service.UnsImportExportService]()
	if strings.HasSuffix(fileName, ".zip") {
		unsImportService.ImportGlobal(fileName, fileSize, respWriter, reader)
	} else if strings.HasSuffix(fileName, ".json") {
		unsImportService.Import(fileName, fileSize, respWriter, reader)
	}
}
