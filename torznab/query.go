package torznab

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type Query map[string]interface{}

func (query Query) Episode() (s string) {
	if season, ok := query["season"].(string); ok {
		s += fmt.Sprintf("S%s", padLeft(season, "0", 2))
	}
	if ep, ok := query["ep"].(string); ok {
		s += fmt.Sprintf("E%s", padLeft(ep, "0", 2))
	}
	return s
}

// Returns a combination of the q, ep and season parameters formatted for text search
func (query Query) Keywords() string {
	keywords := []string{}

	if q, hasQ := query["q"].(string); hasQ {
		keywords = append(keywords, q)
	}

	if ep := query.Episode(); ep != "" {
		keywords = append(keywords, ep)
	}

	return strings.Join(keywords, " ")
}

func ParseQuery(v url.Values) (Query, error) {
	query := Query{}

	for k, vals := range v {
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

func padLeft(str, pad string, length int) string {
	for len(str) < length {
		str = pad + str
	}
	return str
}
