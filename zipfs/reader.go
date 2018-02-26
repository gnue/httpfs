package zipfs

import (
	"archive/zip"
	"io"
)

// Reader
type Reader interface {
	File() []*zip.File
	Comment() string
	RegisterDecompressor(method uint16, dcomp zip.Decompressor)
}

type reader struct {
	r *zip.Reader
}

func (r *reader) File() []*zip.File {
	return r.r.File
}

func (r *reader) Comment() string {
	return r.r.Comment
}

func (r *reader) RegisterDecompressor(method uint16, dcomp zip.Decompressor) {
	r.r.RegisterDecompressor(method, dcomp)
}

// ReadCloser
type ReadCloser interface {
	Reader
	io.Closer
}

type readCloser struct {
	rc *zip.ReadCloser
}

func (r *readCloser) File() []*zip.File {
	return r.rc.File
}

func (r *readCloser) Comment() string {
	return r.rc.Comment
}

func (r *readCloser) RegisterDecompressor(method uint16, dcomp zip.Decompressor) {
	r.rc.RegisterDecompressor(method, dcomp)
}

func (r *readCloser) Close() error {
	return r.rc.Close()
}
