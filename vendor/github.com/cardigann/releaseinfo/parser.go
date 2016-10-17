package releaseinfo

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

const (
	airDateFormat = "2006-01-02"
)

var titleRegex = []*regexp2.Regexp{
	//1. Multi-Part episodes without a title (S01E05.S01E06)
	regexp2.MustCompile(`^(?:\W*S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:(?:[ex]){1,2}(?<episode>\d{1,3}(?!\d+)))+){2,}`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//2. Episodes without a title, Single (S01E05, 1x05) AND Multi (S01E04E05, 1x04x05, etc)
	regexp2.MustCompile(`^(?:S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{2,3}(?!\d+)))+)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//3. Anime - [SubGroup] Title Absolute Episode Number + Season+Episode
	regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\](?:_|-|\s|\.)?)(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+(?<absoluteepisode>\d{2,3}))+(?:_|-|\s|\.)+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]){1,2}(?<episode>\d{2}(?!\d+)))+).*?(?<hash>[(\[]\w{8}[)\]])?(?:$|\.)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//4. Anime - [SubGroup] Title Season+Episode + Absolute Episode Number
	regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\](?:_|-|\s|\.)?)(?<title>.+?)(?:[-_\W](?<![()\[!]))+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]){1,2}(?<episode>\d{2}(?!\d+)))+)(?:(?:_|-|\s|\.)+(?<absoluteepisode>(?<!\d+)\d{2,3}(?!\d+)))+.*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//5. Anime - [SubGroup] Title Season+Episode
	regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\](?:_|-|\s|\.)?)(?<title>.+?)(?:[-_\W](?<![()\[!]))+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:[ex]|\W[ex]){1,2}(?<episode>\d{2}(?!\d+)))+)(?:\s|\.).*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//6. Anime - [SubGroup] Title with trailing number Absolute Episode Number
	regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>[^-]+?\d+?)[-_. ]+(?:[-_. ]?(?<absoluteepisode>\d{3}(?!\d+)))+(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//7. Anime - [SubGroup] Title - Absolute Episode Number
	regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>.+?)(?:[. ]-[. ](?<absoluteepisode>\d{2,3}(?!\d+|[-])))+(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//8. Anime - [SubGroup] Title Absolute Episode Number
	regexp2.MustCompile(`^\[(?<subgroup>.+?)\][-_. ]?(?<title>.+?)[-_. ]+(?:[-_. ]?(?<absoluteepisode>\d{2,3}(?!\d+)))+(?:[-_. ]+(?<special>special|ova|ovd))?.*?(?<hash>\[\w{8}\])?(?:$|\.mkv)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//9. Anime - Title Season EpisodeNumber + Absolute Episode Number [SubGroup]
	regexp2.MustCompile(`^(?<title>.+?)(?:[-_\W](?<![()\[!]))+(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:[ex]|\W[ex]){1,2}(?<episode>\d{2}(?!\d+)))+).+?(?:[-_. ]?(?<absoluteepisode>\d{3}(?!\d+)))+.+?\[(?<subgroup>.+?)\](?:$|\.mkv)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//10. Anime - Title Absolute Episode Number [SubGroup]
	regexp2.MustCompile(`^(?<title>.+?)(?:(?:_|-|\s|\.)+(?<absoluteepisode>\d{3}(?!\d+)))+(?:.+?)\[(?<subgroup>.+?)\].*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//11. Anime - Title Absolute Episode Number [Hash]
	regexp2.MustCompile(`^(?<title>.+?)(?:(?:_|-|\s|\.)+(?<absoluteepisode>\d{2,3}(?!\d+)))+(?:[-_. ]+(?<special>special|ova|ovd))?[-_. ]+.*?(?<hash>\[\w{8}\])(?:$|\.)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//12. Episodes with airdate AND season/episode number, capture season/epsiode only
	regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airdate>\d{4}\W+[0-1][0-9]\W+[0-3][0-9])(?!\W+[0-3][0-9])[-_. ](?:s?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+)))(?:[ex](?<episode>(?<!\d+)(?:\d{1,3})(?!\d+)))`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//13. Episodes with airdate AND season/episode number
	regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airyear>\d{4})\W+(?<airmonth>[0-1][0-9])\W+(?<airday>[0-3][0-9])(?!\W+[0-3][0-9]).+?(?:s?(?<season>(?<!\d+)(?:\d{1,2})(?!\d+)))(?:[ex](?<episode>(?<!\d+)(?:\d{1,3})(?!\d+)))`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//14. Multi-episode Repeated (S01E05 - S01E06, 1x05 - 1x06, etc)
	regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:(?:[ex]|[-_. ]e){1,2}(?<episode>\d{1,3}(?!\d+)))+){2,}`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//15. Episodes with a title, Single episodes (S01E05, 1x05, etc) & Multi-episode (S01E05E06, S01E05-06, S01E05 E06, etc) **
	regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:[ex]|\W[ex]|_){1,2}(?<episode>\d{2,3}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{2,3}(?!\d+)))*)\W?(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//16. Mini-Series, treated as season 1, episodes are labelled as Part01, Part 01, Part.1
	regexp2.MustCompile(`^(?<title>.+?)(?:\W+(?:(?:Part\W?|(?<!\d+\W+)e)(?<episode>\d{1,2}(?!\d+)))+)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//17. Mini-Series, treated as season 1, episodes are labelled as Part One/Two/Three/...Nine, Part.One, Part_One
	regexp2.MustCompile(`^(?<title>.+?)(?:\W+(?:Part[-._ ](?<episode>One|Two|Three|Four|Five|Six|Seven|Eight|Nine)(?>[-._ ])))`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//18. Mini-Series, treated as season 1, episodes are labelled as XofY
	regexp2.MustCompile(`^(?<title>.+?)(?:\W+(?:(?<episode>(?<!\d+)\d{1,2}(?!\d+))of\d+)+)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//19. Supports Season 01 Episode 03
	regexp2.MustCompile(`(?:.*(?:\""|^))(?<title>.*?)(?:[-_\W](?<![()\[]))+(?:\W?Season\W?)(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:\W|_)+(?:Episode\W)(?:[-_. ]?(?<episode>(?<!\d+)\d{1,2}(?!\d+)))+`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//20. Multi-episode release with no space between series title and season (S01E11E12)
	regexp2.MustCompile(`(?:.*(?:^))(?<title>.*?)(?:\W?|_)S(?<season>(?<!\d+)\d{2}(?!\d+))(?:E(?<episode>(?<!\d+)\d{2}(?!\d+)))+`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//21. Multi-episode with single episode numbers (S6.E1-E2, S6.E1E2, S6E1E2, etc)
	regexp2.MustCompile(`^(?<title>.+?)[-_. ]S(?<season>(?<!\d+)(?:\d{1,2}|\d{4})(?!\d+))(?:[-_. ]?[ex]?(?<episode>(?<!\d+)\d{1,2}(?!\d+)))+`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//22. Single episode season or episode S1E1 or S1-E1
	regexp2.MustCompile(`(?:.*(?:\""|^))(?<title>.*?)(?:\W?|_)S(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:\W|_)?E(?<episode>(?<!\d+)\d{1,2}(?!\d+))`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//23. 3 digit season S010E05
	regexp2.MustCompile(`(?:.*(?:\""|^))(?<title>.*?)(?:\W?|_)S(?<season>(?<!\d+)\d{3}(?!\d+))(?:\W|_)?E(?<episode>(?<!\d+)\d{1,2}(?!\d+))`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//24. 5 digit episode number with a title
	regexp2.MustCompile(`^(?:(?<title>.+?)(?:_|-|\s|\.)+)(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+)))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>(?<!\d+)\d{5}(?!\d+)))`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//25. 5 digit multi-episode with a title
	regexp2.MustCompile(`^(?:(?<title>.+?)(?:_|-|\s|\.)+)(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+)))(?:(?:[-_. ]{1,3}ep){1,2}(?<episode>(?<!\d+)\d{5}(?!\d+)))+`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//26. Separated season and episode numbers S01 - E01
	regexp2.MustCompile(`^(?<title>.+?)(?:_|-|\s|\.)+S(?<season>\d{2}(?!\d+))(\W-\W)E(?<episode>(?<!\d+)\d{2}(?!\d+))(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//27. Season only releases
	regexp2.MustCompile(`^(?<title>.+?)\W(?:S|Season)\W?(?<season>\d{1,2}(?!\d+))(\W+|_|$)(?<extras>EXTRAS|SUBPACK)?(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//28. digit season only releases
	regexp2.MustCompile(`^(?<title>.+?)\W(?:S|Season)\W?(?<season>\d{4}(?!\d+))(\W+|_|$)(?<extras>EXTRAS|SUBPACK)?(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//29. Episodes with a title and season/episode in square brackets
	regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+\[S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>(?<!\d+)\d{2}(?!\d+|i|p)))+\])\W?(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//30. Supports 103/113 naming
	regexp2.MustCompile(`^(?<title>.+?)?(?:(?:[-_\W](?<![()\[!]))+(?<season>(?<!\d+)[1-9])(?<episode>[1-9][0-9]|[0][1-9])(?![a-z]|\d+))+`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//31. Episodes with airdate
	regexp2.MustCompile(`^(?<title>.+?)?\W*(?<airyear>\d{4})\W+(?<airmonth>[0-1][0-9])\W+(?<airday>[0-3][0-9])(?!\W+[0-3][0-9])`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//32. Supports 1103/1113 naming
	regexp2.MustCompile(`^(?<title>.+?)?(?:(?:[-_\W](?<![()\[!]))*(?<season>(?<!\d+|\(|\[|e|x)\d{2})(?<episode>(?<!e|x)\d{2}(?!p|i|\d+|\)|\]|\W\d+)))+(\W+|_|$)(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//33. 4 digit episode number
	//Episodes without a title, Single (S01E05, 1x05) AND Multi (S01E04E05, 1x04x05, etc)
	regexp2.MustCompile(`^(?:S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{4}(?!\d+|i|p)))+)(\W+|_|$)(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//34. 4 digit episode number
	//Episodes with a title, Single episodes (S01E05, 1x05, etc) & Multi-episode (S01E05E06, S01E05-06, S01E05 E06, etc)
	regexp2.MustCompile(`^(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]|\W[ex]|_){1,2}(?<episode>\d{4}(?!\d+|i|p)))+)\W?(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//35. Episodes with single digit episode number (S01E1, S01E5E6, etc)
	regexp2.MustCompile(`^(?<title>.*?)(?:(?:[-_\W](?<![()\[!]))+S?(?<season>(?<!\d+)\d{1,2}(?!\d+))(?:(?:\-|[ex]){1,2}(?<episode>\d{1}))+)+(\W+|_|$)(?!\\)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//36. iTunes Season 1\05 Title (Quality).ext
	regexp2.MustCompile(`^(?:Season(?:_|-|\s|\.)(?<season>(?<!\d+)\d{1,2}(?!\d+)))(?:_|-|\s|\.)(?<episode>(?<!\d+)\d{1,2}(?!\d+))`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//37. Anime - Title Absolute Episode Number (e66)
	regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:(?:_|-|\s|\.)+(?:e|ep)(?<absoluteepisode>\d{2,3}))+.*?(?<hash>\[\w{8}\])?(?:$|\.)`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//38. Anime - [SubGroup] Title Episode Absolute Episode Number ([SubGroup] Series Title Episode 01)
	regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)[-_. ](?:Episode)(?:[-_. ]+(?<absoluteepisode>(?<!\d+)\d{2,3}(?!\d+)))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//39. Anime - Title Absolute Episode Number
	regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:[-_. ]+(?<absoluteepisode>(?<!\d+)\d{2,3}(?!\d+)))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//40. Anime - Title {Absolute Episode Number}
	regexp2.MustCompile(`^(?:\[(?<subgroup>.+?)\][-_. ]?)?(?<title>.+?)(?:(?:[-_\W](?<![()\[!]))+(?<absoluteepisode>(?<!\d+)\d{2,3}(?!\d+)))+(?:_|-|\s|\.)*?(?<hash>\[.{8}\])?(?:$|\.)?`,
		regexp2.IgnoreCase|regexp2.Compiled),

	//41. Extant, terrible multi-episode naming (extant.10708.hdtv-lol.mp4)
	regexp2.MustCompile(`^(?<title>.+?)[-_. ](?<season>[0]?\d?)(?:(?<episode>\d{2}){2}(?!\d+))[-_. ]`,
		regexp2.IgnoreCase|regexp2.Compiled),
}

var rejectHashedReleasesRegex = []*regexp2.Regexp{
	// Generic match for md5 and mixed-case hashes.
	regexp2.MustCompile(`^[0-9a-zA-Z]{32}`, regexp2.Compiled),

	// Generic match for shorter lower-case hashes.
	regexp2.MustCompile(`^[a-z0-9]{24}$`, regexp2.Compiled),

	// Format seen on some NZBGeek releases
	// Be very strict with these coz they are very close to the valid 101 ep numbering.
	regexp2.MustCompile(`^[A-Z]{11}\d{3}$`, regexp2.Compiled),
	regexp2.MustCompile(`^[a-z]{12}\d{3}$`, regexp2.Compiled),

	//Backup filename (Unknown origins)
	regexp2.MustCompile(`^Backup_\d{5,}S\d{2}-\d{2}$`, regexp2.Compiled),

	//123 - Started appearing December 2014
	regexp2.MustCompile(`^123$`, regexp2.Compiled),

	//abc - Started appearing January 2015
	regexp2.MustCompile(`^abc$`, regexp2.Compiled|regexp2.IgnoreCase),

	//b00bs - Started appearing January 2015
	regexp2.MustCompile(`^b00bs$`, regexp2.Compiled|regexp2.IgnoreCase),
}

var (
	reversedTitleRegex = regexp2.MustCompile(`[-._ ](p027|p0801|\d{2}E\d{2}S)[-._ ]`,
		regexp2.Compiled)

	normalizeRegex = regexp2.MustCompile(`((?:\b|_)(?<!^)(a(?!$)|an|the|and|or|of)(?:\b|_))|\W|_`,
		regexp2.IgnoreCase|regexp2.Compiled)

	simpleTitleRegex = regexp2.MustCompile(`(?:480[ip]|720[ip]|1080[ip]|[xh][\W_]?26[45]|DD\W?5\W1|[<>?*:|]|848x480|1280x720|1920x1080|(8|10)b(it)?)\s*`,
		regexp2.IgnoreCase|regexp2.Compiled)

	websitePrefixRegex = regexp2.MustCompile(`^\[\s*[a-z]+(\.[a-z]+)+\s*\][- ]*`,
		regexp2.IgnoreCase|regexp2.Compiled)

	airDateRegex = regexp2.MustCompile(`^(.*?)(?<!\d)((?<airyear>\d{4})[_.-](?<airmonth>[0-1][0-9])[_.-](?<airday>[0-3][0-9])|(?<airmonth>[0-1][0-9])[_.-](?<airday>[0-3][0-9])[_.-](?<airyear>\d{4}))(?!\d)`,
		regexp2.IgnoreCase|regexp2.Compiled)

	sixDigitAirDateRegex = regexp2.MustCompile(`(?<=[_.-])(?<airdate>(?<!\d)(?<airyear>[1-9]\d{1})(?<airmonth>[0-1][0-9])(?<airday>[0-3][0-9]))(?=[_.-])`,
		regexp2.IgnoreCase|regexp2.Compiled)

	cleanReleaseGroupRegex = regexp2.MustCompile(`^(.*?[-._ ](S\d+E\d+)[-._ ])|(-(RP|1|NZBGeek|Obfuscated|sample))+$`,
		regexp2.IgnoreCase|regexp2.Compiled)

	cleanTorrentSuffixRegex = regexp2.MustCompile(`\[(?:ettv|rartv|rarbg|cttv)\]$`,
		regexp2.IgnoreCase|regexp2.Compiled)

	releaseGroupRegex = regexp2.MustCompile(`-(?<releasegroup>[a-z0-9]+)(?<!WEB-DL|480p|720p|1080p|2160p)(?:\b|[-._ ])`,
		regexp2.IgnoreCase|regexp2.Compiled)

	animeReleaseGroupRegex = regexp2.MustCompile(`^(?:\[(?<subgroup>(?!\s).+?(?<!\s))\](?:_|-|\s|\.)?)`,
		regexp2.IgnoreCase|regexp2.Compiled)

	yearInTitleRegex = regexp2.MustCompile(`^(?<title>.+?)(?:\W|_)?(?<year>\d{4})`,
		regexp2.IgnoreCase|regexp2.Compiled)

	requestInfoRegex = regexp2.MustCompile(`\[.+?\]`,
		regexp2.Compiled)
)

func Parse(title string) (*EpisodeInfo, error) {
	if !validateBeforeParsing(title) {
		return nil, errors.New("Title failed to validate before parsing")
	}

	if match, _ := reversedTitleRegex.MatchString(title); match {
		titleWithoutExtension := removeFileExtension(title)

		title = reverseString(titleWithoutExtension) + strings.TrimPrefix(title, titleWithoutExtension)
	}

	simpleTitle, err := simpleTitleRegex.Replace(title, "", 0, -1)
	if err != nil {
		return nil, err
	}

	simpleTitle = removeFileExtension(simpleTitle)

	// TODO: Quick fix stripping [url] - prefixes.
	simpleTitle, err = websitePrefixRegex.Replace(simpleTitle, "", 0, -1)
	if err != nil {
		return nil, err
	}

	simpleTitle, err = cleanTorrentSuffixRegex.Replace(simpleTitle, "", 0, -1)
	if err != nil {
		return nil, err
	}

	if airDateMatch, _ := airDateRegex.FindStringMatch(simpleTitle); airDateMatch != nil {
		simpleTitle = airDateMatch.GroupByNumber(1).String() +
			getMatchGroupString(airDateMatch, "airyear") + "." +
			getMatchGroupString(airDateMatch, "airmonth") + "." +
			getMatchGroupString(airDateMatch, "airday")
	}

	if sixDigitAirDateMatch, _ := sixDigitAirDateRegex.FindStringMatch(simpleTitle); sixDigitAirDateMatch != nil {
		var airYear = getMatchGroupString(sixDigitAirDateMatch, "airyear")
		var airMonth = getMatchGroupString(sixDigitAirDateMatch, "airmonth")
		var airDay = getMatchGroupString(sixDigitAirDateMatch, "airday")

		if airMonth != "00" || airDay != "00" {
			var fixedDate = fmt.Sprintf("20%s.%s.%s", airYear, airMonth, airDay)
			simpleTitle = strings.Replace(simpleTitle,
				getMatchGroupString(sixDigitAirDateMatch, "airdate"), fixedDate, -1)
		}
	}

	for _, regex := range titleRegex {
		match, err := regex.FindStringMatch(simpleTitle)
		if match == nil {
			continue
		}

		result, err := parseMatchCollection(&matchCollection{regex, match})
		if err != nil {
			return nil, err
		} else if result == nil {
			continue
		}

		if result.FullSeason && containsIgnoreCase(title, "Special") {
			result.FullSeason = false
			result.Special = true
		}

		result.Language = ParseLanguage(title)
		// log.Printf("Language parsed: %q", result.Language)

		result.Quality = ParseQuality(title)
		// log.Printf("Quality parsed: %s", result.Quality)

		result.ReleaseGroup = ParseReleaseGroup(title)

		var subGroup = getSubGroup(match)
		if !isNullOrWhiteSpace(subGroup) {
			result.ReleaseGroup = subGroup
		}

		result.ReleaseHash = getReleaseHash(match)
		// if !isNullOrWhiteSpace(result.ReleaseHash) {
		// 	log.Printf("Release Hash parsed: %q", result.ReleaseHash)
		// }

		return result, nil
	}

	return nil, fmt.Errorf("Unable to parse %q", title)
}

// ParsePath extracts episode info from a full path
func ParsePath(path string) (*EpisodeInfo, error) {
	path = normalizePath(path)
	name := filepath.Base(path)
	dir := filepath.Base(filepath.Dir(path))
	ext := filepath.Ext(path)

	for _, test := range []string{
		name,             // path name
		dir + " " + name, // dirname + name of the file
		dir + ext,        // dirname + extentions
	} {
		result, err := Parse(test)
		if err == nil {
			return result, err
		}
	}

	return nil, fmt.Errorf("Unable to parse path %s", path)
}

// ParseSeriesName parses just the name of the series from the title. If Parse fails internally
// then the passed in title is cleaned up and returned as-is.
func ParseSeriesName(title string) string {
	// log.Printf("Parsing series name from %q", title)

	parseResult, err := Parse(title)
	if err != nil {
		return CleanSeriesTitle(title)
	}

	return parseResult.SeriesTitle
}

func CleanSeriesTitle(title string) string {
	//If Title only contains numbers return it as is.
	if _, err := strconv.Atoi(title); err == nil {
		return title
	}

	title = optionalReplace(normalizeRegex, title, "")
	return removeSpace(removeAccent(strings.ToLower(title)))
}

func ParseReleaseGroup(title string) string {
	title = removeSpace(title)
	title = removeFileExtension(title)
	title = optionalReplace(duplicateSpacesRegex, title, " ")
	title = optionalReplace(websitePrefixRegex, title, " ")

	if animeMatch, _ := animeReleaseGroupRegex.FindStringMatch(title); animeMatch != nil {
		return animeMatch.GroupByName("subgroup").String()
	}

	title = optionalReplace(cleanReleaseGroupRegex, title, "")

	if match, _ := releaseGroupRegex.FindStringMatch(title); match != nil {
		var releaseGroupCaptures = getMatchGroupCaptures(match, "releasegroup")

		if len(releaseGroupCaptures) > 0 {
			lastGroup := releaseGroupCaptures[len(releaseGroupCaptures)-1].String()

			if _, err := strconv.Atoi(lastGroup); err == nil {
				return ""
			}

			return lastGroup
		}
	}

	return ""
}

func getSeriesTitleInfo(title string) SeriesTitleInfo {
	seriesTitleInfo := SeriesTitleInfo{
		Title: title,
	}

	match, _ := yearInTitleRegex.FindStringMatch(title)
	if match == nil {
		seriesTitleInfo.TitleWithoutYear = title
	} else {
		seriesTitleInfo.TitleWithoutYear = match.GroupByName("title").String()
		seriesTitleInfo.Year, _ = strconv.Atoi(match.GroupByName("year").String())
	}

	return seriesTitleInfo
}

func getMatchGroupString(m *regexp2.Match, group string) string {
	if result := m.GroupByName(group); result != nil {
		return result.String()
	}
	return ""
}

func getMatchGroupCaptures(m *regexp2.Match, group string) []regexp2.Capture {
	if result := m.GroupByName(group); result != nil {
		return result.Captures
	}
	return nil
}

func hasGroup(m *regexp2.Match, group string) bool {
	if m == nil {
		return false
	}
	if result := m.GroupByName(group); result != nil && result.Length > 0 {
		return true
	}
	return false
}

func dumpGroups(m *regexp2.Match) {
	if m == nil {
		return
	}
	for _, group := range m.Groups() {
		captures := []string{}

		for _, capture := range group.Captures {
			captures = append(captures, capture.String())
		}

		log.Printf("Group %s [%d]: %s",
			group.Name, len(group.Captures), strings.Join(captures, ", "))
	}
}

type matchCollection struct {
	Regexp *regexp2.Regexp
	Match  *regexp2.Match
}

func (mc *matchCollection) Next() error {
	m, err := mc.Regexp.FindNextMatch(mc.Match)
	if err != nil {
		return err
	}
	mc.Match = m
	return nil
}

func newAirDate(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
}

func parseMatchCollection(col *matchCollection) (*EpisodeInfo, error) {
	seriesName := getMatchGroupString(col.Match, "title")
	seriesName = strings.Replace(seriesName, "_", " ", -1)
	seriesName = strings.Replace(seriesName, ".", " ", -1)
	seriesName = optionalReplace(requestInfoRegex, seriesName, "")
	seriesName = removeSpace(seriesName)

	var result *EpisodeInfo
	var airYear int

	if match := col.Match.GroupByName("airyear"); match != nil {
		airYear, _ = strconv.Atoi(match.String())
	}

	if airYear < 1900 {
		var distinctSeasons = map[int]struct{}{}
		var seasons = []int{}

		for _, seasonCapture := range getMatchGroupCaptures(col.Match, "season") {
			if parsedSeason, err := strconv.Atoi(seasonCapture.String()); err == nil {
				distinctSeasons[parsedSeason] = struct{}{}
				seasons = append(seasons, parsedSeason)
			}
		}

		//If no season was found it should be treated as a mini series and season 1
		if len(seasons) == 0 {
			distinctSeasons[1] = struct{}{}
			seasons = append(seasons, 1)
		}

		//If more than 1 season was parsed go to the next REGEX (A multi-season release is unlikely)
		if len(distinctSeasons) > 1 {
			return nil, nil
		}

		result = &EpisodeInfo{
			SeasonNumber:           seasons[0],
			EpisodeNumbers:         []int{},
			AbsoluteEpisodeNumbers: []int{},
		}

		for col.Match != nil {
			var episodeCaptures = getMatchGroupCaptures(col.Match, "episode")
			var absoluteEpisodeCaptures = getMatchGroupCaptures(col.Match, "absoluteepisode")

			//Allows use to return a list of 0 episodes (We can handle that as a full season release)
			if len(episodeCaptures) > 0 {
				first, _ := parseNumber(episodeCaptures[0].String())
				last, _ := parseNumber(episodeCaptures[len(episodeCaptures)-1].String())

				if first > last {
					return nil, nil
				}

				for i := first; i <= last; i++ {
					result.EpisodeNumbers = append(result.EpisodeNumbers, i)
				}
			}

			if len(absoluteEpisodeCaptures) > 0 {
				first, _ := parseNumber(absoluteEpisodeCaptures[0].String())
				last, _ := parseNumber(absoluteEpisodeCaptures[len(absoluteEpisodeCaptures)-1].String())

				if first > last {
					return nil, nil
				}

				for i := first; i <= last; i++ {
					result.AbsoluteEpisodeNumbers = append(result.AbsoluteEpisodeNumbers, i)
				}

				if hasGroup(col.Match, "special") {
					result.Special = true
				}
			}

			if len(episodeCaptures) == 0 && len(absoluteEpisodeCaptures) == 0 {
				//Check to see if this is an "Extras" or "SUBPACK" release, if it is, return NULL
				//Todo: Set a "Extras" flag in EpisodeParseResult if we want to download them ever
				if !isNullOrWhiteSpace(getMatchGroupString(col.Match, "extras")) {
					return nil, nil
				}

				result.FullSeason = true
			}

			if len(result.AbsoluteEpisodeNumbers) > 0 && len(result.EpisodeNumbers) == 0 {
				result.SeasonNumber = 0
			}

			if err := col.Next(); err != nil {
				return nil, nil
			}
		}
	} else {
		//Try to Parse as a daily show
		airMonth, _ := strconv.Atoi(getMatchGroupString(col.Match, "airmonth"))
		airDay, _ := strconv.Atoi(getMatchGroupString(col.Match, "airday"))

		//Swap day and month if month is bigger than 12 (scene fail)
		if airMonth > 12 {
			var tempDay = airDay
			airDay = airMonth
			airMonth = tempDay
		}

		airDate := newAirDate(airYear, airMonth, airDay)

		//Check if episode is in the future (most likely a parse error)
		if airDate.After(time.Now()) {
			return nil, fmt.Errorf("Invalid date %d-%d-%d", airYear, airMonth, airDay)
		}

		result = &EpisodeInfo{
			EpisodeNumbers:         []int{},
			AbsoluteEpisodeNumbers: []int{},
			AirDate:                airDate.Format(airDateFormat),
		}
	}

	result.SeriesTitle = seriesName
	result.SeriesTitleInfo = getSeriesTitleInfo(result.SeriesTitle)

	return result, nil
}

func validateBeforeParsing(title string) bool {
	var titleWithoutExtension = removeFileExtension(title)

	for _, regex := range rejectHashedReleasesRegex {
		if m, _ := regex.FindStringMatch(titleWithoutExtension); m != nil {
			// log.Printf("Rejected Hashed Release Title: " + title)
			return false
		}
	}

	return true
}

func getSubGroup(m *regexp2.Match) string {
	if subGroup := m.GroupByName("subgroup"); subGroup != nil {
		return subGroup.String()
	}

	return ""
}

func getReleaseHash(m *regexp2.Match) string {
	if hash := m.GroupByName("hash"); hash != nil {
		hashValue := strings.Trim(hash.String(), "[]")

		if hashValue == "1280x720" {
			return ""
		}

		return hashValue
	}

	return ""
}
