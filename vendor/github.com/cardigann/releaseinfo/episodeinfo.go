package releaseinfo

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"

	"golang.org/x/text/language"
)

var (
	wordDelimiterRegex = regexp2.MustCompile(`(\s|\.|,|_|-|=|\|)+`,
		regexp2.Compiled)

	punctuationRegex = regexp2.MustCompile(`[^\w\s]`,
		regexp2.Compiled)

	commonWordRegex = regexp2.MustCompile(`\b(a|an|the|and|or|of)\b\s?`,
		regexp2.IgnoreCase|regexp2.Compiled)

	duplicateSpacesRegex = regexp2.MustCompile(`\s{2,}`,
		regexp2.Compiled)
)

// Normalize a series title, removing all spaces, punctation and whitespace
func NormalizeSeriesTitle(title string) string {
	title = optionalReplace(wordDelimiterRegex, title, "")
	title = optionalReplace(punctuationRegex, title, "")
	title = optionalReplace(commonWordRegex, title, "")
	title = optionalReplace(duplicateSpacesRegex, title, "")

	return strings.ToLower(removeSpace(title))
}

type EpisodeInfo struct {
	SeriesTitle            string
	SeriesTitleInfo        SeriesTitleInfo
	Quality                QualityModel
	SeasonNumber           int
	EpisodeNumbers         []int
	AbsoluteEpisodeNumbers []int
	AirDate                string
	Language               language.Tag
	FullSeason             bool
	Special                bool
	ReleaseGroup           string
	ReleaseHash            string
}

type SeriesTitleInfo struct {
	Title            string
	TitleWithoutYear string
	Year             int
}

func (i SeriesTitleInfo) Normalize() string {
	return NormalizeSeriesTitle(i.Title)
}

func (i SeriesTitleInfo) Equal(title string) bool {
	return i.Normalize() == NormalizeSeriesTitle(title)
}

func (i EpisodeInfo) IsDaily() bool {
	return removeSpace(i.AirDate) != ""
}

func (i EpisodeInfo) IsAbsoluteNumbering() bool {
	return len(i.AbsoluteEpisodeNumbers) > 0
}

func (i EpisodeInfo) IsPossibleSpecialEpisode() bool {
	return removeSpace(i.AirDate) != "" &&
		removeSpace(i.SeriesTitle) != "" &&
		(len(i.EpisodeNumbers) == 0 || i.SeasonNumber == 0) ||
		(removeSpace(i.SeriesTitle) != "" && i.Special)
}

func (i EpisodeInfo) String() string {
	episodeString := "[Unknown Episode]"

	if i.IsDaily() && len(i.EpisodeNumbers) == 0 {
		episodeString = fmt.Sprintf("%s", i.AirDate)
	} else if i.FullSeason {
		episodeString = fmt.Sprintf("S%02d", i.SeasonNumber)
	} else if len(i.EpisodeNumbers) > 0 {
		episodes := []string{}
		for _, episode := range i.EpisodeNumbers {
			episodes = append(episodes, fmt.Sprintf("%02d", episode))
		}
		episodeString = fmt.Sprintf("S%02dE%s", i.SeasonNumber, strings.Join(episodes, "-"))
	} else if len(i.AbsoluteEpisodeNumbers) > 0 {
		episodes := []string{}
		for _, episode := range i.AbsoluteEpisodeNumbers {
			episodes = append(episodes, fmt.Sprintf("%03d", episode))
		}
		episodeString = strings.Join(episodes, "-")
	}

	return fmt.Sprintf("%s - %s (%s)", i.SeriesTitle, episodeString, i.Quality)
}
