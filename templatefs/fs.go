// Go lang template http.FileSystem
package templatefs

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Engine interface {
	Render(input []byte) []byte
	PageInfo(input []byte) *Page
	Exts() []string
}

type FileSystem struct {
	Engines      map[string]Engine
	PageTemplete *template.Template
	FileSystem   http.FileSystem
	reExts       *regexp.Regexp
}

func New(fs http.FileSystem, layout *template.Template, engines ...Engine) *FileSystem {
	if layout == nil {
		layout = defaultLayout
	}

	t := &FileSystem{
		FileSystem:   fs,
		Engines:      make(map[string]Engine),
		PageTemplete: layout,
	}

	for _, e := range engines {
		t.RegEngine(e)
	}

	t.compileExts()

	return t
}

func (t *FileSystem) Open(name string) (http.File, error) {
	f, err := t.open(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()

	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if finfo.IsDir() {
		return f, nil
	}

	e := t.FindEngine(finfo.Name())
	if e != nil {
		file := &File{
			file:   f,
			finfo:  finfo,
			render: func(b []byte) ([]byte, error) { return t.render(e, b, finfo) },
		}

		return file, nil
	}

	return f, nil
}

func (t *FileSystem) open(name string) (http.File, error) {
	fs := t.FileSystem
	f, err := fs.Open(name)
	if err == nil {
		return f, nil
	}

	for ext, _ := range t.Engines {
		f, err = fs.Open(name + ext)
		if err == nil {
			return f, nil
		}
	}

	return nil, os.ErrNotExist
}

func (t *FileSystem) RegEngine(e Engine, exts ...string) {
	if exts == nil {
		exts = e.Exts()
	}

	for _, ext := range exts {
		t.Engines[ext] = e
	}

	if t.reExts != nil {
		t.compileExts()
	}
}

func (t *FileSystem) FindEngine(name string) Engine {
	ext := filepath.Ext(name)
	return t.Engines[ext]
}

type data struct {
	FileInfo os.FileInfo
	Page     *Page
	Title    string
	CSS      string
	Body     string
}

func (t *FileSystem) render(e Engine, b []byte, finfo os.FileInfo) ([]byte, error) {
	output := e.Render(b)
	output = t.postRender(output)

	pinfo := e.PageInfo(b)
	d := &data{FileInfo: finfo, Page: pinfo, Title: pinfo.Title, CSS: pinfo.CSS, Body: string(output)}

	tmpl := t.PageTemplete.Lookup(pinfo.Layout)
	if tmpl == nil {
		tmpl = t.PageTemplete
		pinfo.Layout = tmpl.Name()
	}

	var page bytes.Buffer
	err := tmpl.Execute(&page, d)
	if err != nil {
		return nil, err
	}

	return page.Bytes(), nil
}

var reHref = regexp.MustCompile(`(?i)\shref="[^"]+"`)

func (t *FileSystem) postRender(b []byte) []byte {
	fn := func(b []byte) []byte {
		sub := t.reExts.FindSubmatchIndex(b)
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

	return reHref.ReplaceAllFunc(b, fn)
}

func (t *FileSystem) compileExts() {
	exts := make([]string, 0, len(t.Engines))

	for ext, _ := range t.Engines {
		exts = append(exts, ext)
	}

	t.reExts = regexp.MustCompile(`(?i)\shref="(.+)(` + strings.Join(exts, "|") + `)"`)
}
