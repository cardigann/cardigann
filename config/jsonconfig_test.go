package config

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestJSONConfig(t *testing.T) {
	dir, err := ioutil.TempDir("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

	j := jsonConfig{
		dirs:       []string{dir},
		defaultDir: dir,
	}

	j.Set("section1", "my_key", "true")
	j.Set("section1", "another_key", "llamas1")
	j.Set("section2", "my_key", "true")
	j.Set("section2", "another_key", "llamas2")

	s1, err := j.Section("section1")
	if err != nil {
		t.Fatal(err)
	}

	s2, err := j.Section("section2")
	if err != nil {
		t.Fatal(err)
	}

	if s1["another_key"] != "llamas1" {
		t.Fatalf("section1[another_key] is %q, expected llamas1", s1["another_key"])
	}

	if s2["another_key"] != "llamas2" {
		t.Fatalf("section1[another_key] is %q, expected llamas2", s2["another_key"])
	}
}
