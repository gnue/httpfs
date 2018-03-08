package gitfs

import (
	"strconv"
	"strings"
)

type Object struct {
	Repo *Repo
	Id   string
}

func (obj *Object) Read() ([]byte, error) {
	return obj.Repo.Exec("cat-file", "-p", obj.Id)
}

func (obj *Object) Size() (int64, error) {
	b, err := obj.Repo.Exec("cat-file", "-s", obj.Id)
	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(b))

	return strconv.ParseInt(s, 10, 64)
}

func (obj *Object) Type() (string, error) {
	b, err := obj.Repo.Exec("cat-file", "-s", obj.Id)
	if err != nil {
		return "", err
	}

	s := strings.TrimSpace(string(b))

	return s, nil
}
