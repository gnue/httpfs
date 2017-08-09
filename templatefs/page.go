package templatefs

type Page struct {
	Title      string
	Layout     string
	Permalink  string
	Published  string
	Categories []string
	Tags       []string
}
