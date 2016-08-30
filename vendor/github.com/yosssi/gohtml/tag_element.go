package gohtml

import "bytes"

// A tagElement represents a tag element of an HTML document.
type tagElement struct {
	tagName     string
	startTagRaw string
	endTagRaw   string
	children    []element
}

// write writes a tag to the buffer.
func (e *tagElement) write(bf *bytes.Buffer, indent int) {
	writeLine(bf, e.startTagRaw, indent)
	for _, c := range e.children {
		var childIndent int
		if e.endTagRaw != "" {
			childIndent = indent + 1
		} else {
			childIndent = indent
		}
		c.write(bf, childIndent)
	}
	if e.endTagRaw != "" {
		writeLine(bf, e.endTagRaw, indent)
	}
}

// appendChild append an element to the element's children.
func (e *tagElement) appendChild(child element) {
	e.children = append(e.children, child)
}
