package gohtml

import (
	"bytes"
	"os"
	"testing"
)

func TestWriterSetLastElement(t *testing.T) {
	wr := NewWriter(os.Stdout)
	wr.SetLastElement("test")
	if wr.lastElement != "test" {
		t.Errorf("Invalid lastElement. [expected: %s][actual: %s]", "test", wr.lastElement)
	}
}

func TestWriterWrite(t *testing.T) {
	wr := NewWriter(os.Stdout)
	n, err := wr.Write([]byte("<html><head><title>This is a title.</title></head><body><p>test</p></body></html>"))
	if err != nil {
		t.Errorf("An error occurred. [error: %s]", err.Error())
	}
	expected := 129
	if n != expected {
		t.Errorf("Invalid return value. [expected: %d][actual: %d]", expected, n)
	}

	wr = NewWriter(os.Stdout)
	n, err = wr.Write([]byte(""))
	if err != nil {
		t.Errorf("An error occurred. [error: %s]", err.Error())
	}
	expected = 0
	if n != expected {
		t.Errorf("Invalid return value. [expected: %d][actual: %d]", expected, n)
	}
}

func TestNewWriter(t *testing.T) {
	wr := NewWriter(os.Stdout)
	if wr.writer != os.Stdout || wr.lastElement != defaultLastElement || wr.bf.Len() != 0 {
		t.Errorf("Invalid Writer. [expected: %+v][actual: %+v]", &Writer{writer: os.Stdout, lastElement: defaultLastElement, bf: &bytes.Buffer{}}, wr)
	}
}
