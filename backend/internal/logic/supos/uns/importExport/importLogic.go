package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/share/spring"
	"context"
	"io"
	"strings"
)

func Import(ctx context.Context, fileName string, fileSize int64, respWriter io.Writer, reader io.Reader, readAt io.ReaderAt) {
	unsImportService := spring.GetBean[*service.UnsImportExportService]()
	if strings.HasSuffix(fileName, ".zip") {
		unsImportService.ImportGlobal(ctx, fileName, fileSize, respWriter, readAt)
	} else if strings.HasSuffix(fileName, ".json") {
		unsImportService.Import(ctx, fileName, fileSize, respWriter, reader)
	}
}
