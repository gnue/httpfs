package templatefs

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

type File struct {
	engine       Engine
	pageTemplete *template.Template
	file         http.File
	finfo        os.FileInfo
	r            *bytes.Reader
	reExts       *regexp.Regexp
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

	if f.reExts != nil {
		output = reHref.ReplaceAllFunc(output, f.replaceHref)
	}

	var page bytes.Buffer
	err = f.pageTemplete.Execute(&page, &data{e.Title(b), string(output)})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(page.Bytes()), nil
}

var reHref = regexp.MustCompile(`(?i)\shref="[^"]+"`)

func (f *File) replaceHref(b []byte) []byte {
	sub := f.reExts.FindSubmatchIndex(b)
	if sub == nil {
		return b
	}

	s := string(b[sub[2]:sub[3]])
	u, err := url.Parse(s)
	if err == nil && u.Host != "" {
		return b
	}

	return append(b[:sub[4]], b[sub[5]:]...)
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
