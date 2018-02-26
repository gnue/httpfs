// Go lang git http.FileSystem
package gitfs

import (
	"net/http"
)

type FileSystem struct {
	Repo   *Repo
	Branch string
}

func New(repo string, branch string) *FileSystem {
	if branch == "" {
		branch = "master"
	}

	return &FileSystem{&Repo{repo}, branch}
}

func (g *FileSystem) Open(name string) (http.File, error) {
	return g.Repo.Open(name, g.Branch)
}
