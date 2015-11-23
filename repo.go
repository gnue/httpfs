package gitfs

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Repo struct {
	Path string
}

func (repo *Repo) Exec(args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.Path

	return cmd.Output()
}

func (repo *Repo) Open(fname string, treeish string) (*File, error) {
	finfo, err := repo.Stat(fname, treeish)
	if err != nil {
		return nil, err
	}

	f := &File{repo: repo, finfo: finfo}

	if !finfo.IsDir() {
		object := &Object{repo, finfo.object}
		b, err := object.Read()
		if err != nil {
			return nil, err
		}
		f.r = bytes.NewReader(b)
	}

	return f, nil
}

func (repo *Repo) Stat(fname string, treeish string) (*FileInfo, error) {
	if fname == "/" {
		b, err := repo.Exec("rev-parse", treeish)
		if err != nil {
			return nil, err
		}
		object := strings.TrimRight(string(b), "\r\n")
		return &FileInfo{mode: "040000", objtype: "tree", object: object, path: "/"}, nil
	}

	fname = strings.TrimPrefix(fname, "/")
	b, err := repo.Exec("ls-tree", "-l", treeish, fname)
	if err != nil {
		return nil, err
	}

	s := strings.TrimRight(string(b), "\r\n")
	return parseInfo(s)
}

func parseInfo(s string) (*FileInfo, error) {
	line := strings.SplitN(s, "\t", 2)
	if len(line) < 2 {
		return nil, os.ErrNotExist
	}

	info := strings.Fields(line[0])
	if len(info) < 4 {
		return nil, os.ErrNotExist
	}

	size, _ := strconv.ParseInt(info[3], 10, 64)

	finfo := &FileInfo{
		mode:    info[0],
		objtype: info[1],
		object:  info[2],
		size:    size,
		path:    line[1],
	}

	return finfo, nil
}
