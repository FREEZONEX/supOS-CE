package httpUtils

import (
	"io"
	"net/http"
)

type emptyFlusher struct{}

func (emptyFlusher) Flush() {}
func HttpFlusher(w io.Writer) http.Flusher {
	flusher, isFlusher := w.(http.Flusher)
	if !isFlusher {
		flusher = emptyFlusher{}
	}
	return flusher
}
func Flush(w io.Writer) {
	flusher, isFlusher := w.(http.Flusher)
	if isFlusher {
		flusher.Flush()
	}
}
