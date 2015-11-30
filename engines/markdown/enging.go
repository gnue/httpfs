package markdown

import (
	"github.com/gnue/templatefs"
	"github.com/russross/blackfriday"
)

type Engine struct {
	Extensions int
	HtmlFlags  int
	CSS        string
}

func (e *Engine) Render(input []byte) []byte {
	r := blackfriday.HtmlRenderer(e.HtmlFlags, "", e.CSS)
	return blackfriday.Markdown(input, r, e.Extensions)
}

func (e *Engine) PageInfo(input []byte) *templatefs.Page {
	return &templatefs.Page{Title: getTitle(input)}
}

func (e *Engine) Exts() []string {
	return []string{".md", ".markdown"}
}
