package indexer

import "github.com/cardigann/cardigann/torznab"

type categoryMap map[string]torznab.Category

func (mapping categoryMap) Categories() torznab.Categories {
	cats := torznab.Categories{}
	added := map[int]bool{}

	for _, c := range mapping {
		if _, exists := added[c.ID]; exists {
			continue
		}
		cats = append(cats, c)
		added[c.ID] = true
	}

	return cats
}

func (mapping categoryMap) Resolve(cat torznab.Category) []string {
	var matched bool
	var results = []string{}

	// check for exact matches
	for localID, mappedCat := range mapping {
		if mappedCat.ID == cat.ID {
			results = append(results, localID)
			matched = true
		}
	}

	// check for matches on the parent categories of the mapped categories
	// e.g asked for Movies, but only had a more specific mapping for Movies/Blu-ray
	if !matched {
		for localID, mappedCat := range mapping {
			if torznab.ParentCategory(mappedCat).ID == cat.ID {
				results = append(results, localID)
				matched = true
			}
		}
	}

	// finally check for matches on the parent category of the requested cat
	// e.g. asked for Movies/Blu-ray but no mapping, so try Movies instead
	if !matched {
		parent := torznab.ParentCategory(cat)
		for localID, mappedCat := range mapping {
			if mappedCat.ID == parent.ID {
				results = append(results, localID)
				matched = true
			}
		}
	}

	return results
}

func (mapping categoryMap) ResolveAll(cats ...torznab.Category) []string {
	results := []string{}

	for _, cat := range cats {
		results = append(results, mapping.Resolve(cat)...)
	}

	return results
}
