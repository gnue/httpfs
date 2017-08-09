package zipfs

import (
	"regexp"
	"strings"
)

func NewIgnore(pattern []string) (*regexp.Regexp, error) {
	a := make([]string, 0, len(pattern))

	for _, s := range pattern {
		a = append(a, regexp.QuoteMeta(s))
	}

	return regexp.Compile("(?i)(^|/)(" + strings.Join(a, "|") + ")(/|$)")
}
