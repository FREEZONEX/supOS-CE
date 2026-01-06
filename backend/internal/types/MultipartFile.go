package types

import (
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type MultipartFile struct {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer

	FileName    string
	ContentType string
	Size        int64
}

func FormFile(r *http.Request, name string) (*MultipartFile, error) {
	err := r.ParseMultipartForm(8 * 1024) //8k
	fileData, handler, err := r.FormFile(name)
	if err != nil && r.MultipartForm != nil && r.MultipartForm.File != nil {
		for k, fhs := range r.MultipartForm.File {
			if len(fhs) > 0 {
				f, er := fhs[0].Open()
				if er == nil {
					fileData, handler = f, fhs[0]
					logx.Info("FormFile", k, fhs[0].Filename, ", instead of :", name)
					break
				} else {
					err = er
				}
			}
		}
	}
	if err != nil {
		return nil, err
	}
	return &MultipartFile{
		Reader:      fileData,
		ReaderAt:    fileData,
		Seeker:      fileData,
		Closer:      fileData,
		FileName:    handler.Filename,
		ContentType: handler.Header.Get("Content-Type"),
		Size:        handler.Size,
	}, nil
}
