package gohtml

import (
	"bytes"
	"strings"
)

// A textElement represents a text element of an HTML document.
type textElement struct {
	text string
}

// write writes a text to the buffer.
func (e *textElement) write(bf *bytes.Buffer, indent int) {
	lines := strings.Split(strings.Trim(unifyLineFeed(e.text), "\n"), "\n")
	for _, line := range lines {
		writeLineFeed(bf)
		writeIndent(bf, indent)
		bf.WriteString(line)
	}
}
