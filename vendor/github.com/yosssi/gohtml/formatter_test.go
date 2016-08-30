package gohtml

import "testing"

func TestFormat(t *testing.T) {
	s := `<!DOCTYPE html><html><head><title>This is a title.</title></head><body><p>Line1<br>Line2</p><br/></body></html> <!-- aaa -->`
	actual := Format(s)
	expected := `<!DOCTYPE html>
<html>
  <head>
    <title>
      This is a title.
    </title>
  </head>
  <body>
    <p>
      Line1
      <br>
      Line2
    </p>
    <br/>
  </body>
</html>
<!-- aaa -->`
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}

func TestFormatWithLineNo(t *testing.T) {
	s := `<!DOCTYPE html><html><head><title>This is a title.</title></head><body><p>Line1<br>Line2</p><br/></body></html> <!-- aaa -->`
	actual := FormatWithLineNo(s)
	expected := ` 1  <!DOCTYPE html>
 2  <html>
 3    <head>
 4      <title>
 5        This is a title.
 6      </title>
 7    </head>
 8    <body>
 9      <p>
10        Line1
11        <br>
12        Line2
13      </p>
14      <br/>
15    </body>
16  </html>
17  <!-- aaa -->`
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}
