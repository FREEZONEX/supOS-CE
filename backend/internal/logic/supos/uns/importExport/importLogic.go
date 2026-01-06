package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/share/spring"
	"io"
)

func Import(fileName string, fileSize int64, resp io.Writer) (w io.Writer, waiter func()) {
	unsImportService := spring.GetBean[*service.UnsImportExportService]()
	return unsImportService.ImportUns(fileName, fileSize, resp)
}
func ImportUnsByReader(fileName string, fileSize int64, respWriter io.Writer, reader io.Reader) {
	unsImportService := spring.GetBean[*service.UnsImportExportService]()
	unsImportService.ImportUnsDirect(fileName, fileSize, respWriter, reader)
}
