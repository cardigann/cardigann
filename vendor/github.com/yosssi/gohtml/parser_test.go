package gohtml

import "testing"

func TestParse(t *testing.T) {
	s := `<!DOCTYPE html><html><head><title>This is a title.</title></head><body><p>Line1<br>Line2</p><br/></body></html><!-- aaa --><a>`
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
<!-- aaa -->
<a>`
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}

func TestAppendElement(t *testing.T) {
	htmlDoc := &htmlDocument{}
	tagElem := &tagElement{}
	textElem := &textElement{text: "test text"}
	appendElement(htmlDoc, tagElem, textElem)
	if len(tagElem.children) != 1 || tagElem.children[0] != textElem {
		t.Errorf("tagElement.children is invalid. [expected: %+v][actual: %+v]", []element{textElem}, tagElem.children)
	}
	appendElement(htmlDoc, nil, textElem)
	if len(htmlDoc.elements) != 1 || htmlDoc.elements[0] != textElem {
		t.Errorf("htmlDocument.elements is invalid. [expected: %+v][actual: %+v]", []element{textElem}, htmlDoc.elements)
	}
}
