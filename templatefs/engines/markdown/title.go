package markdown

import (
	"bufio"
	"bytes"
	"regexp"
)

const DEFAULT_TITLE = ""

var (
	header = regexp.MustCompile("^#{1,6}\\s+(.*?)(?:\\s+#+)?\\s*$")
	under  = regexp.MustCompile("^(-+|=+)\\s*$")
	blank  = regexp.MustCompile("^\\s*$")
)

func getTitle(input []byte) string {
	reader := bytes.NewReader(input)
	scanner := bufio.NewScanner(reader)

	var line []byte

	for scanner.Scan() {
		b := scanner.Bytes()
		if !blank.Match(b) {
			line = b
			break
		}
	}

	r := header.FindSubmatch(line)
	if 1 < len(r) {
		return string(r[1])
	}

	if scanner.Scan() {
		r := under.FindSubmatch(scanner.Bytes())

		if 1 < len(r) {
			return string(line)
		}
	}

	return DEFAULT_TITLE
}
