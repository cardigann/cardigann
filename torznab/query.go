package torznab

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Query map[string]interface{}

func (q Query) Keywords() string {
	k, ok := q["q"]
	if ok {
		return k.(string)
	} else {
		return ""
	}
}

func ParseQuery(r *http.Request) (Query, error) {
	query := Query{}

	for k, vals := range r.URL.Query() {
		switch k {
		case "t":
			continue

		case "q", "ep", "season", "apikey", "offset", "limit":
			query[k] = vals[0]

		case "cat":
			catInts, err := splitInts(vals[0], ",")
			if err != nil {
				return Query{}, fmt.Errorf("Unable to parse cats %q", vals[0])
			}
			query["cat"] = catInts

		default:
			log.Printf("Unknown torznab request key %q", k)
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
