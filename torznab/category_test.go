package torznab

import (
	"reflect"
	"testing"
)

func TestCategoryParent(t *testing.T) {
	for _, test := range []struct {
		cat, parent Category
	}{
		{CategoryTV_Anime, CategoryTV},
		{CategoryTV_HD, CategoryTV},
		{CategoryPC_PhoneAndroid, CategoryPC},
		{CategoryOther_Hashed, CategoryOther},
	} {
		c := ParentCategory(test.cat)
		if c != test.parent {
			t.Fatalf("Expected to resolve %s to %s, instead got %s", test.cat, test.parent, c)
		}
	}
}

func TestCategorySubset(t *testing.T) {
	s := AllCategories.Subset(5030, 5040)
	expected := Categories{CategoryTV_SD, CategoryTV_HD}

	if !reflect.DeepEqual(s, expected) {
		t.Fatalf("Expected to resolve to %s, instead got %s", expected, s)
	}
}
