package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/internal/types"
	"backend/share/spring"
	"io"
)

func Import(file *types.MultipartFile, w io.Writer) {
	unsImportService := spring.GetBean[*service.UnsImportExportService]()
	unsImportService.ImportUns(file, w)
}
