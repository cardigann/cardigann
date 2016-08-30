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

func TestCategoryMapping(t *testing.T) {
	cats := CategoryMapping{
		1:  CategoryTV_Anime,
		2:  CategoryMovies_BluRay,
		3:  CategoryTV_Documentary,
		4:  CategoryTV_Sport,
		5:  CategoryAudio,
		6:  CategoryAudio,
		7:  CategoryMovies,
		8:  CategoryAudio_Video,
		10: CategoryTV,
		12: CategoryTV,
	}

	for _, test := range []struct {
		tCats     []Category
		localCats []int
	}{
		{tCats: []Category{CategoryTV_SD, CategoryTV_HD}, localCats: []int{10, 12}},
	} {
		if r := cats.ResolveAll(test.tCats...); !reflect.DeepEqual(r, test.localCats) {
			t.Fatalf("Expected to resolve %#v to %#v, instead got %#v",
				test.tCats, test.localCats, r)
		}
	}
}
