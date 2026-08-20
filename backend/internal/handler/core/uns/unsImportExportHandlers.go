package uns

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	domainuns "backend/internal/domain/uns"
	respx "backend/internal/httpx"
	"backend/internal/infra/eventstream"
	"backend/internal/logic/logicx"
	"backend/internal/svc"
)

type unsExportParam struct {
	ExportType string      `json:"exportType"`
	Folders    flexibleIDs `json:"folders"`
	Files      flexibleIDs `json:"files"`
}

type unsExportReq struct {
	ExportType      string          `json:"exportType"`
	Folders         flexibleIDs     `json:"folders"`
	Files           flexibleIDs     `json:"files"`
	UnsExportParam  *unsExportParam `json:"unsExportParam"`
	CheckSmallFile  *bool           `json:"checkSmallFile"`
	SourceFlowParam any             `json:"sourceFlowExportParam"`
	EventFlowParam  any             `json:"eventFlowExportParam"`
}

type flexibleIDs []int64

func (ids *flexibleIDs) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*ids = nil
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var raw []any
	if err := dec.Decode(&raw); err != nil {
		return err
	}
	out := make([]int64, 0, len(raw))
	for _, item := range raw {
		switch value := item.(type) {
		case json.Number:
			if n, err := value.Int64(); err == nil {
				out = append(out, n)
			}
		case float64:
			out = append(out, int64(value))
		case string:
			if n, err := strconv.ParseInt(strings.TrimSpace(value), 10, 64); err == nil {
				out = append(out, n)
			}
		}
	}
	*ids = out
	return nil
}

func UnsExportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseUnsExportReq(r)
		if err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}
		data, err := svcCtx.App.UNS.ExportData(r.Context(), req.toCommand())
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		_ = json.NewEncoder(w).Encode(data)
	}
}

func UnsExportGlobalHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req, err := parseUnsExportReq(r)
		if err != nil {
			respx.WriteError(w, respx.NewHTTPError(http.StatusBadRequest, "invalid request: "+err.Error()))
			return
		}
		data, err := svcCtx.App.UNS.ExportData(r.Context(), req.toCommand())
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		data["requested"] = map[string]any{
			"sourceFlowExportParam": req.SourceFlowParam,
			"eventFlowExportParam":  req.EventFlowParam,
		}
		var buf bytes.Buffer
		zw := zip.NewWriter(&buf)
		f, err := zw.Create("uns.json")
		if err != nil {
			respx.WriteError(w, err)
			return
		}
		encoded, _ := json.MarshalIndent(data, "", "  ")
		if _, err := f.Write(encoded); err != nil {
			respx.WriteError(w, err)
			return
		}
		if err := zw.Close(); err != nil {
			respx.WriteError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", `attachment; filename="global-export.zip"`)
		_, _ = w.Write(buf.Bytes())
	}
}

func UnsImportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeSSEHeaders(w)
		flusher, _ := w.(http.Flusher)
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			writeImportStatus(w, flusher, domainuns.ImportStatus{Code: 500, Msg: "invalid multipart request: " + err.Error(), Task: "Parse File", Finished: true, Progress: 100})
			return
		}
		file, header, err := r.FormFile("file")
		if err != nil {
			writeImportStatus(w, flusher, domainuns.ImportStatus{Code: 500, Msg: "file is required", Task: "Parse File", Finished: true, Progress: 100})
			return
		}
		defer file.Close()
		raw, err := io.ReadAll(file)
		if err != nil {
			writeImportStatus(w, flusher, domainuns.ImportStatus{Code: 500, Msg: err.Error(), Task: "Read File", Finished: true, Progress: 100})
			return
		}
		_, _ = svcCtx.App.UNS.ImportDataStream(r.Context(), header.Filename, raw, logicx.UserID(r.Context()), func(status domainuns.ImportStatus) {
			writeImportStatus(w, flusher, status)
		})
	}
}

func UnsNewMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeEventSourceHeaders(w)
		flusher, _ := w.(http.Flusher)
		id := strings.TrimSpace(r.URL.Query().Get("id"))
		global := strings.EqualFold(r.URL.Query().Get("globalTopology"), "true")
		if global {
			_, _ = fmt.Fprint(w, "data: Connected\n\n")
			if flusher != nil {
				flusher.Flush()
			}
			streamGlobalTopology(svcCtx, w, r, flusher)
			return
		}
		if id == "" {
			_, _ = fmt.Fprint(w, "data: Connected\n\n")
			if flusher != nil {
				flusher.Flush()
			}
			return
		}
		match := func(event eventstream.Event) bool {
			if event.Type != "uns.payload" {
				return false
			}
			return eventstream.ContainsKey(event.Keys, "id:"+id)
		}
		ch, cancel := eventstream.Subscribe(match)
		defer cancel()
		_, _ = fmt.Fprint(w, "data: Connected\n\n")
		if flusher != nil {
			flusher.Flush()
		}
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		_ = svcCtx
		for {
			select {
			case event, ok := <-ch:
				if !ok {
					return
				}
				raw, _ := json.Marshal(event.Data)
				_, _ = fmt.Fprintf(w, "data: %s\n\n", raw)
				if flusher != nil {
					flusher.Flush()
				}
			case <-ticker.C:
				_, _ = fmt.Fprint(w, ": keep-alive\n\n")
				if flusher != nil {
					flusher.Flush()
				}
			case <-r.Context().Done():
				return
			}
		}
	}
}

func writeSSEHeaders(w http.ResponseWriter) {
	w.Header().Del("Keep-Alive")
	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}

func writeEventSourceHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}

func writeImportStatus(w io.Writer, flusher http.Flusher, status domainuns.ImportStatus) {
	raw, _ := json.Marshal(status)
	_, _ = w.Write(append(raw, '\n', '\n'))
	if flusher != nil {
		flusher.Flush()
	}
}

func parseUnsExportReq(r *http.Request) (unsExportReq, error) {
	defer r.Body.Close()
	var req unsExportReq
	if r.Body == nil {
		return req, nil
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
		return req, err
	}
	return req, nil
}

func (r unsExportReq) toCommand() domainuns.ExportCommand {
	param := r.UnsExportParam
	if param == nil {
		param = &unsExportParam{ExportType: r.ExportType, Folders: r.Folders, Files: r.Files}
	}
	return domainuns.ExportCommand{
		ExportType: strings.ToUpper(strings.TrimSpace(param.ExportType)),
		Folders:    cleanIDs([]int64(param.Folders)),
		Files:      cleanIDs([]int64(param.Files)),
	}
}

func cleanIDs(ids []int64) []int64 {
	out := make([]int64, 0, len(ids))
	seen := map[int64]struct{}{}
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
