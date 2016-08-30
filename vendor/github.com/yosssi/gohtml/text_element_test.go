package gohtml

import (
	"bytes"
	"testing"
)

func TestTextElementWrite(t *testing.T) {
	textElem := &textElement{text: "Test text"}
	bf := &bytes.Buffer{}
	textElem.write(bf, 0)
	actual := bf.String()
	expected := "Test text"
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}
