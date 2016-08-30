package gohtml

import (
	"bytes"
	"strconv"
	"strings"
)

// Format parses the input HTML string, formats it and returns the result.
func Format(s string) string {
	return parse(s).html()
}

// Format parses the input HTML string, formats it and returns the result with line no.
func FormatWithLineNo(s string) string {
	return AddLineNo(Format(s))
}

func AddLineNo(s string) string {
	lines := strings.Split(s, "\n")
	maxLineNoStrLen := len(strconv.Itoa(len(lines)))
	bf := &bytes.Buffer{}
	for i, line := range lines {
		lineNoStr := strconv.Itoa(i + 1)
		if i > 0 {
			bf.WriteString("\n")
		}
		bf.WriteString(strings.Repeat(" ", maxLineNoStrLen-len(lineNoStr)) + lineNoStr + "  " + line)
	}
	return bf.String()

}
