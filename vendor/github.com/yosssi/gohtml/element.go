package gohtml

import "bytes"

// An element represents an HTML element.
type element interface {
	write(*bytes.Buffer, int)
}
