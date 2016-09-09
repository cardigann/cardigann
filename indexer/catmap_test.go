package indexer

import (
	"reflect"
	"sort"
	"testing"

	"github.com/cardigann/cardigann/torznab"
)

func TestCategoryMap(t *testing.T) {
	cats := categoryMap{
		"1":   torznab.CategoryTV_Anime,
		"2":   torznab.CategoryMovies_BluRay,
		"3":   torznab.CategoryTV_Documentary,
		"4":   torznab.CategoryTV_Sport,
		"5":   torznab.CategoryAudio,
		"6":   torznab.CategoryAudio,
		"7":   torznab.CategoryMovies_HD,
		"8":   torznab.CategoryAudio_Video,
		"10":  torznab.CategoryTV,
		"12":  torznab.CategoryTV,
		"xyz": torznab.CategoryAudio_Foreign,
	}

	for _, test := range []struct {
		tCats     []torznab.Category
		localCats []string
	}{
		{tCats: []torznab.Category{torznab.CategoryTV_Anime, torznab.CategoryTV_SD}, localCats: []string{"1", "10", "12"}},
		{tCats: []torznab.Category{torznab.CategoryAudio_Foreign}, localCats: []string{"xyz"}},
		{tCats: []torznab.Category{torznab.CategoryMovies}, localCats: []string{"2", "7"}},
	} {
		r := cats.ResolveAll(test.tCats...)
		sort.Sort(sort.StringSlice(r))

		if !reflect.DeepEqual(r, test.localCats) {
			t.Fatalf("Expected to resolve %#v to %#v, instead got %#v",
				test.tCats, test.localCats, r)
		}
	}
}
