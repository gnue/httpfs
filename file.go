package templatefs

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
)

type File struct {
	engine       Engine
	pageTemplete *template.Template
	file         http.File
	finfo        os.FileInfo
	r            *bytes.Reader
	postRender   func([]byte) []byte
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) Read(p []byte) (int, error) {
	if f.r == nil {
		r, err := f.newReader()
		if err != nil {
			return 0, err
		}
		f.r = r
	}

	i, err := f.r.Read(p)
	return i, err

	return f.r.Read(p)
}

type data struct {
	Page  *Page
	Title string
	Body  string
}

func (f *File) newReader() (*bytes.Reader, error) {
	if f.finfo.IsDir() {
		return nil, os.ErrInvalid
	}

	b, err := ioutil.ReadAll(f.file)
	if err != nil {
		return nil, err
	}
	f.file.Seek(0, os.SEEK_CUR)

	e := f.engine
	output := e.Render(b)

	if f.postRender != nil {
		output = f.postRender(output)
	}

	pinfo := e.PageInfo(b)
	d := &data{Page: pinfo, Title: pinfo.Title, Body: string(output)}

	t := f.pageTemplete.Lookup(pinfo.Layout)
	if t == nil {
		t = f.pageTemplete
		pinfo.Layout = t.Name()
	}

	var page bytes.Buffer
	err = t.Execute(&page, d)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(page.Bytes()), nil
}

func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, os.ErrInvalid
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.r == nil {
		r, err := f.newReader()
		if err != nil {
			return 0, err
		}
		f.r = r
	}

	return f.r.Seek(offset, whence)
}

func (f *File) Stat() (os.FileInfo, error) {
	if f.r == nil {
		r, err := f.newReader()
		if err != nil {
			return nil, err
		}
		f.r = r
	}

	return &FileInfo{finfo: f.finfo, size: f.r.Size()}, nil
}
