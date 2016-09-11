package torznab

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// Query represents a torznab query
type Query struct {
	Type          string
	Q, Ep, Season string
	Limit, Offset int
	Extended      bool
	Categories    []int
	APIKey        string
}

// Episode returns either the season + episode in the format S00E00 or just the season as S00 if
// no episode has been specified.
func (query Query) Episode() (s string) {
	if query.Season != "" {
		s += fmt.Sprintf("S%s", padLeft(query.Season, "0", 2))
	}
	if query.Ep != "" {
		s += fmt.Sprintf("E%s", padLeft(query.Ep, "0", 2))
	}
	return s
}

// Keywords returns the query formatted as search keywords
func (query Query) Keywords() string {
	tokens := []string{}

	if query.Q != "" {
		tokens = append(tokens, query.Q)
	}

	if query.Season != "" || query.Ep != "" {
		tokens = append(tokens, query.Episode())
	}

	return strings.Join(tokens, " ")
}

// Encode returns the query as a url query string
func (query Query) Encode() string {
	v := url.Values{}

	if query.Type != "" {
		v.Set("t", query.Type)
	} else {
		v.Set("t", "search")
	}

	if query.Q != "" {
		v.Set("q", query.Q)
	}

	if query.Ep != "" {
		v.Set("ep", query.Ep)
	}

	if query.Season != "" {
		v.Set("season", query.Season)
	}

	if query.Offset != 0 {
		v.Set("offset", strconv.Itoa(query.Offset))
	}

	if query.Limit != 0 {
		v.Set("limit", strconv.Itoa(query.Limit))
	}

	if query.Extended {
		v.Set("extended", "1")
	}

	if query.APIKey != "" {
		v.Set("apikey", query.APIKey)
	}

	if len(query.Categories) > 0 {
		cats := []string{}

		for _, cat := range query.Categories {
			cats = append(cats, strconv.Itoa(cat))
		}

		v.Set("cat", strings.Join(cats, ","))
	}

	return v.Encode()
}

// ParseQuery takes the query string parameters for a torznab query and parses them
func ParseQuery(v url.Values) (Query, error) {
	query := Query{}

	for k, vals := range v {
		switch k {
		case "t":
			if len(vals) > 1 {
				return query, errors.New("Multiple t parameters not allowed")
			}
			query.Type = vals[0]

		case "q":
			query.Q = strings.Join(vals, " ")

		case "ep":
			if len(vals) > 1 {
				return query, errors.New("Multiple ep parameters not allowed")
			}
			query.Ep = vals[0]

		case "season":
			if len(vals) > 1 {
				return query, errors.New("Multiple season parameters not allowed")
			}
			query.Season = vals[0]

		case "apikey":
			if len(vals) > 1 {
				return query, errors.New("Multiple apikey parameters not allowed")
			}
			query.APIKey = vals[0]

		case "limit":
			if len(vals) > 1 {
				return query, errors.New("Multiple limit parameters not allowed")
			}
			limit, err := strconv.Atoi(vals[0])
			if err != nil {
				return query, err
			}
			query.Limit = limit

		case "offset":
			if len(vals) > 1 {
				return query, errors.New("Multiple offset parameters not allowed")
			}
			offset, err := strconv.Atoi(vals[0])
			if err != nil {
				return query, err
			}
			query.Offset = offset

		case "extended":
			if len(vals) > 1 {
				return query, errors.New("Multiple extended parameters not allowed")
			}
			extended, err := strconv.ParseBool(vals[0])
			if err != nil {
				return query, err
			}
			query.Extended = extended

		case "cat":
			query.Categories = []int{}
			for _, val := range vals {
				ints, err := splitInts(val, ",")
				if err != nil {
					return Query{}, fmt.Errorf("Unable to parse cats %q", vals[0])
				}
				query.Categories = append(query.Categories, ints...)
			}

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
