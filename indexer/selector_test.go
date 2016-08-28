package indexer

import "testing"

func TestSelectorIsEmpty(t *testing.T) {
	for idx, test := range []struct {
		block    selectorBlock
		expected bool
	}{
		{selectorBlock{}, true},
		{selectorBlock{Selector: "*"}, false},
		{selectorBlock{TextVal: "llamas"}, false},
	} {
		result := test.block.IsEmpty()
		if result != test.expected {
			t.Fatalf("Row #%d expected %#v, got %#v", idx+1, test.expected, result)
		}
	}
}
