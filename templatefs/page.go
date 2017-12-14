package templatefs

type Page struct {
	Title      string
	CSS        string
	Layout     string
	Permalink  string
	Published  string
	Categories []string
	Tags       []string
}
