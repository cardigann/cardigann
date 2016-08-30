package gohtml

import (
	"bytes"
	"io"
)

// A Writer represents a formatted HTML source codes writer.
type Writer struct {
	writer      io.Writer
	lastElement string
	bf          *bytes.Buffer
}

// SetLastElement set the lastElement to the Writer.
func (wr *Writer) SetLastElement(lastElement string) *Writer {
	wr.lastElement = lastElement
	return wr
}

// Write writes the parameter.
func (wr *Writer) Write(p []byte) (n int, err error) {
	wr.bf.Write(p)
	if bytes.HasSuffix(p, []byte(wr.lastElement)) {
		return wr.writer.Write([]byte(Format(wr.bf.String()) + "\n"))
	}
	return 0, nil
}

// NewWriter generates a Writer and returns it.
func NewWriter(wr io.Writer) *Writer {
	return &Writer{writer: wr, lastElement: defaultLastElement, bf: &bytes.Buffer{}}
}
