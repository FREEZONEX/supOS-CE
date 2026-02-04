// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package importExport

import (
	"backend/internal/logic/supos/uns/importExport/service"
	"backend/share/spring"
	"context"
	"net/http"

	"backend/internal/types"
)

type ExportLogic struct {
}

func (l ExportLogic) Export(ctx context.Context, w http.ResponseWriter, req *types.ExportReq) (*types.BaseResult, error) {
	serv := spring.GetBean[*service.UnsImportExportService]()
	return serv.Export(ctx, w, req)
}
func (l ExportLogic) ExportGlobal(ctx context.Context, w http.ResponseWriter, req *types.GlobalExportParam) {
	serv := spring.GetBean[*service.UnsImportExportService]()
	serv.ExportGlobal(ctx, w, req)
}
