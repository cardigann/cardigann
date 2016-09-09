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

	for localID, mappedCat := range mapping {
		if mappedCat.ID == cat.ID {
			results = append(results, localID)
			matched = true
		}
	}

	if !matched {
		parent := torznab.ParentCategory(cat)
		for localID, mappedCat := range mapping {
			if mappedCat.ID == parent.ID {
				results = append(results, localID)
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
