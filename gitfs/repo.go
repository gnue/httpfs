package gitfs

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const TIME_ISO = "2006-01-02 15:04:05 -0700"

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

	return f, nil
}

func (repo *Repo) Stat(fname string, treeish string) (*FileInfo, error) {
	if fname == "/" {
		b, err := repo.Exec("rev-parse", treeish)
		if err != nil {
			return nil, err
		}
		object := strings.TrimSpace(string(b))
		return &FileInfo{mode: "040000", objtype: "tree", object: object, path: "/"}, nil
	}

	fname = strings.TrimPrefix(fname, "/")
	b, err := repo.Exec("ls-tree", "-l", treeish, fname)
	if err != nil {
		return nil, err
	}

	s := strings.TrimSpace(string(b))
	finfo, err := parseInfo(s)
	if err != nil {
		return nil, err
	}

	t, err := repo.modTime(fname, treeish)
	if err != nil {
		return nil, err
	}

	finfo.modTime = t

	return finfo, nil
}

func (repo *Repo) modTime(fname string, treeish string) (t time.Time, err error) {
	b, err := repo.Exec("log", "--pretty=%ad", "--date=iso", "-1", treeish, "--", fname)
	if err != nil {
		return
	}

	s := strings.TrimSpace(string(b))
	return time.Parse(TIME_ISO, s)
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
