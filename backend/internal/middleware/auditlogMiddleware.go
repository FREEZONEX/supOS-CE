package middleware

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"strings"

	"backend/internal/domain/audit"
)

const auditRequestBodyLimit = 64 << 10

type AuditLogMiddleware struct {
	auditSvc *audit.Service
}

func NewAuditLogMiddleware(auditSvc *audit.Service) *AuditLogMiddleware {
	return &AuditLogMiddleware{auditSvc: auditSvc}
}

func (m *AuditLogMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if m.auditSvc == nil {
			next(w, r)
			return
		}
		recorder := &auditResponseWriter{ResponseWriter: w}
		r = r.WithContext(m.auditSvc.WithState(r.Context()))
		if payload := captureJSONRequestPayload(r); len(payload) > 0 {
			m.auditSvc.BindRequestPayload(r.Context(), payload)
		}
		next(recorder, r)
		statusCode := recorder.StatusCode()
		m.auditSvc.FlushRequest(r, audit.ResponseMeta{
			StatusCode: statusCode,
			Status:     http.StatusText(statusCode),
		})
	}
}

func captureJSONRequestPayload(r *http.Request) map[string]any {
	if r == nil || r.Body == nil || r.ContentLength <= 0 || r.ContentLength > auditRequestBodyLimit {
		return nil
	}
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if !strings.Contains(contentType, "application/json") {
		return nil
	}
	raw, err := io.ReadAll(r.Body)
	r.Body = io.NopCloser(bytes.NewReader(raw))
	if err != nil || len(raw) == 0 {
		return nil
	}
	var payload map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&payload); err != nil {
		return nil
	}
	return payload
}

type auditResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *auditResponseWriter) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}
	return w.statusCode
}

func (w *auditResponseWriter) WriteHeader(code int) {
	if w.statusCode != 0 {
		return
	}
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *auditResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(data)
}

func (w *auditResponseWriter) Flush() {
	if flusher, ok := w.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *auditResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}

func (w *auditResponseWriter) Push(target string, opts *http.PushOptions) error {
	pusher, ok := w.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return pusher.Push(target, opts)
}

func (w *auditResponseWriter) ReadFrom(r io.Reader) (int64, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	if readerFrom, ok := w.ResponseWriter.(io.ReaderFrom); ok {
		return readerFrom.ReadFrom(r)
	}
	return io.Copy(w.ResponseWriter, r)
}
