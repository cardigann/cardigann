package torznab

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Query map[string]interface{}

func ParseQuery(r *http.Request) (Query, error) {
	query := make(Query)

	for k, vals := range r.URL.Query() {
		switch k {
		case "q", "ep", "season", "apikey":
			query[k] = vals[0]

		case "cat":
			catInts, err := splitInts(vals[0], ",")
			if err != nil {
				return Query{}, fmt.Errorf("Unable to parse cats %q", vals[0])
			}
			query["cat"] = catInts
		}
	}

	return query, nil
}

func splitInts(s, delim string) (i []int, err error) {
	for _, v := range strings.Split(s, delim) {
		vInt, err := strconv.Atoi(v)
		if err != nil {
			return i, err
		}
		i = append(i, vInt)
	}
	return i, err
}
