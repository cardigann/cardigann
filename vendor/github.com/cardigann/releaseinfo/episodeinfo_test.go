package releaseinfo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsPossibleSpecialEpisode(t *testing.T) {
	for idx, test := range []string{
		"Under.the.Dome.S02.Special-Inside.Chesters.Mill.HDTV.x264-BAJSKORV",
		"Under.the.Dome.S02.Special-Inside.Chesters.Mill.720p.HDTV.x264-BAJSKORV",
		"Rookie.Blue.Behind.the.Badge.S05.Special.HDTV.x264-2HD",
	} {
		result, err := Parse(test)

		require.NoError(t, err,
			fmt.Sprintf("Row %d should have no parsing error", idx+1))
		require.True(t, result.IsPossibleSpecialEpisode(),
			fmt.Sprintf("Row %d should be a possible special episode", idx+1))
	}
}

func TestNormalizeTitle(t *testing.T) {
	for idx, test := range []struct {
		title, expected string
	}{
		{"Under.the.Dome.S02.Special-Inside.Chesters.Mill.HDTV.x264-BAJSKORV", "underthedome"},
		{"Under.the.Dome.S02.Special-Inside.Chesters.Mill.720p.HDTV.x264-BAJSKORV", "underthedome"},
		{"Rookie.Blue.Behind.the.Badge.S05.Special.HDTV.x264-2HD", "rookiebluebehindthebadge"},
	} {
		result, err := Parse(test.title)

		require.NoError(t, err,
			fmt.Sprintf("Row %d should have no parsing error", idx+1))
		require.Equal(t, test.expected, result.SeriesTitleInfo.Normalize(),
			fmt.Sprintf("Row %d should match normalized title", idx+1))
		require.True(t, result.SeriesTitleInfo.Equal(test.expected),
			fmt.Sprintf("Row %d should be equal to normalized title", idx+1))
	}
}
