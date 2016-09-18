package fuzzytime

import (
	"sort"
)

// Span represents the range [Begin,End), used to indicate the part
// of a string from which time or date information was parsed.
type Span struct {
	Begin int
	End   int
}

type spanSlice []Span

// implement sort.Interface
func (l spanSlice) Len() int           { return len(l) }
func (l spanSlice) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l spanSlice) Less(i, j int) bool { return l[i].Begin < l[j].Begin }

// sorts and merges a set of spans
func tidySpans(in []Span) []Span {
	// cull empty spans
	tmp := make(spanSlice, 0, len(in))
	for _, foo := range in {
		if foo.Begin != foo.End {
			tmp = append(tmp, foo)
		}
	}

	// sort
	sort.Sort(tmp)

	// merge overlapping
	out := make([]Span, 0, len(tmp))
	for i := 0; i < len(tmp); i++ {
		foo := tmp[i]
		for i+1 < len(tmp) {
			if foo.End < tmp[i+1].Begin {
				break
			}
			// overlapping (or adjacent)
			foo.End = tmp[i+i].End
			i++
		}
		out = append(out, foo)
	}

	return out
}
