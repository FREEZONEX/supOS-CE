package base

import "io"

type ReadCloserWrapper struct {
	tar io.Reader
}

func NewReadCloserWrapper(tar io.Reader) io.ReadCloser {
	return &ReadCloserWrapper{tar: tar}
}
func (bc *ReadCloserWrapper) Read(p []byte) (n int, err error) {
	return bc.tar.Read(p)
}
func (bc *ReadCloserWrapper) Close() error {
	if cc, is := bc.tar.(io.ReadCloser); is {
		return cc.Close()
	}
	return nil
}
