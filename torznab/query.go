package torznab

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// Query represents a torznab query
type Query map[string]interface{}

// Episode returns either the season + episode in the format S00E00 or just the season as S00 if
// no episode has been specified.
func (query Query) Episode() (s string) {
	if season, ok := query["season"].(string); ok {
		s += fmt.Sprintf("S%s", padLeft(season, "0", 2))
	}
	if ep, ok := query["ep"].(string); ok {
		s += fmt.Sprintf("E%s", padLeft(ep, "0", 2))
	}
	return s
}

func (query Query) Limit() (int, bool) {
	if limit, hasLimit := query["limit"].(int); hasLimit {
		return limit, true
	}
	if limit, hasLimit := query["limit"].(string); hasLimit {
		if limitInt, err := strconv.Atoi(limit); err == nil {
			return limitInt, true
		}
	}
	return 0, false
}

// Keywords returns a combination of the q, ep and season parameters formatted for text search
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

// ParseQuery takes the query string parameters for a torznab query and parses them
func ParseQuery(v url.Values) (Query, error) {
	query := Query{}

	for k, vals := range v {
		switch k {
		case "t":
			continue

		case "q", "ep", "season", "apikey", "offset", "limit", "extended":
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
