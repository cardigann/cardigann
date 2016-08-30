package gohtml

import (
	"bytes"
	"testing"
)

func TestWriteLine(t *testing.T) {
	bf := &bytes.Buffer{}
	writeLine(bf, "test", 1)
	actual := bf.String()
	expected := defaultIndentString + "test"
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}

func TestWriteLineFeed(t *testing.T) {
	bf := &bytes.Buffer{}
	writeLine(bf, "test", 0)
	writeLineFeed(bf)
	actual := bf.String()
	expected := "test\n"
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}

func TestWriteIndent(t *testing.T) {
	bf := &bytes.Buffer{}
	writeIndent(bf, 1)
	actual := bf.String()
	expected := defaultIndentString
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}

func TestUnifyLineFeed(t *testing.T) {
	actual := unifyLineFeed("\r\n\n\r")
	expected := "\n\n\n"
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}
