package torznab

type Category struct {
	ID          int
	Name        string
	Description string
}

// Categories are the predefined categories from the nZEDb spec
// See https://github.com/nZEDb/nZEDb/blob/0.x/docs/newznab_api_specification.txt#L608
var (
	Categories = struct {
		TV,
		TV_WebDL,
		TV_Foreign,
		TV_StandardDef,
		TV_HighDef,
		TV_Other,
		TV_Sport,
		TV_Anime,
		TV_Documentary Category
	}{
		Category{5000, "TV", "All of TV"},
		Category{5010, "TV/WEB-DL", "WEB-DL TV"},
		Category{5020, "TV/FOREIGN", "FOREIGN TV"},
		Category{5030, "TV/SD", "SD TV"},
		Category{5040, "TV/HD", "HD TV"},
		Category{5999, "TV/OTHER", "Other TV Content"},
		Category{5060, "TV/Sport", "Sports"},
		Category{5070, "TV/Anime", "Anime"},
		Category{5080, "TV/Documentary", "Documentaries"},
	}
)

func CategoryById(id int) (Category, bool) {
	return Category{}, false
}

type CategoryMapping map[int]Category
