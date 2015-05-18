package zipfs

import (
	"sort"
)

type Ignore []string

func NewIgnore(pattern []string) Ignore {
	ig := make([]string, len(pattern))
	copy(ig, pattern)
	sort.Strings(ig)

	return ig
}

func (ig Ignore) Match(name string) bool {
	i := sort.SearchStrings(ig, name)
	if i < len(ig) && ig[i] == name {
		return true
	}

	return false
}
