package gitfs

import (
	"net/http"
)

type GitFS struct {
	Repo   *Repo
	Branch string
}

func New(repo string, branch string) *GitFS {
	if branch == "" {
		branch = "master"
	}

	return &GitFS{&Repo{repo}, branch}
}

func (g *GitFS) Open(name string) (http.File, error) {
	return g.Repo.Open(name, g.Branch)
}
