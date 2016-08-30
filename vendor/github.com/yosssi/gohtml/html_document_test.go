package gohtml

import "testing"

func TestHTMLDocumentHTML(t *testing.T) {
	s := `<!DOCTYPE html><html><head><title>This is a title.</title></head><body><p>Line1<br>Line2</p><br/></body></html><!-- aaa -->`
	htmlDoc := parse(s)
	actual := htmlDoc.html()
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

func TestHTMLDocumentAppend(t *testing.T) {
	htmlDoc := &htmlDocument{}
	textElem := &textElement{text: "TestText"}
	htmlDoc.append(textElem)
	if len(htmlDoc.elements) != 1 || htmlDoc.elements[0] != textElem {
		t.Errorf("htmlDocument.elements is invalid. [expected: %+v][actual: %+v]", []element{textElem}, htmlDoc.elements)
	}
}
