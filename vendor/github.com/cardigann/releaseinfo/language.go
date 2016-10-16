package releaseinfo

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
	"golang.org/x/text/language"
)

var languageRegex = regexp2.MustCompile(
	`(?:\W|_)(?<italian>\b(?:ita|italian)\b)|(?<german>german\b|videomann)|(?<flemish>flemish)|(?<greek>greek)|(?<french>(?:\W|_)(?:FR|VOSTFR)(?:\W|_))|(?<russian>\brus\b)|(?<dutch>nl\W?subs?)|(?<hungarian>\b(?:HUNDUB|HUN)\b)`,
	regexp2.IgnoreCase|regexp2.Compiled)

var subtitleLanguageRegex = regexp2.MustCompile(
	`.+?[-_. ](?<iso_code>[a-z]{2,3})$`,
	regexp2.Compiled|regexp2.IgnoreCase)

func ParseLanguage(title string) language.Tag {
	lowerTitle := strings.ToLower(title)

	// log.Printf("Parsing language from %s", title)

	if strings.Contains(lowerTitle, "english") {
		return language.English
	}

	if strings.Contains(lowerTitle, "french") {
		return language.French
	}

	if strings.Contains(lowerTitle, "spanish") {
		return language.Spanish
	}

	if strings.Contains(lowerTitle, "danish") {
		return language.Danish
	}

	if strings.Contains(lowerTitle, "dutch") {
		return language.Dutch
	}

	if strings.Contains(lowerTitle, "japanese") {
		return language.Japanese
	}

	if strings.Contains(lowerTitle, "cantonese") {
		return language.MustParse("yue")
	}

	if strings.Contains(lowerTitle, "mandarin") {
		return language.MustParse("cmn")
	}

	if strings.Contains(lowerTitle, "korean") {
		return language.Korean
	}

	if strings.Contains(lowerTitle, "russian") {
		return language.Russian
	}

	if strings.Contains(lowerTitle, "polish") {
		return language.Polish
	}

	if strings.Contains(lowerTitle, "vietnamese") {
		return language.Vietnamese
	}

	if strings.Contains(lowerTitle, "swedish") {
		return language.Swedish
	}

	if strings.Contains(lowerTitle, "norwegian") {
		return language.Norwegian
	}

	if strings.Contains(lowerTitle, "nordic") {
		return language.Norwegian
	}

	if strings.Contains(lowerTitle, "finnish") {
		return language.Finnish
	}

	if strings.Contains(lowerTitle, "turkish") {
		return language.Turkish
	}

	if strings.Contains(lowerTitle, "portuguese") {
		return language.Portuguese
	}

	if strings.Contains(lowerTitle, "hungarian") {
		return language.Hungarian
	}

	match, _ := languageRegex.FindStringMatch(title)

	if match == nil {
		return language.English
	}

	if hasGroup(match, "italian") {
		return language.Italian
	}

	if hasGroup(match, "german") {
		return language.German
	}

	if hasGroup(match, "flemish") {
		return language.MustParse("nl-BE")
	}

	if hasGroup(match, "greek") {
		return language.Greek
	}

	if hasGroup(match, "french") {
		return language.French
	}

	if hasGroup(match, "russian") {
		return language.Russian
	}

	if hasGroup(match, "dutch") {
		return language.Dutch
	}

	if hasGroup(match, "hungarian") {
		return language.Hungarian
	}

	return language.English
}

func ParseSubtitleLanguage(fileName string) (language.Tag, error) {
	// log.Printf("Parsing language from subtitle file: %s", fileName)

	ext := filepath.Ext(fileName)
	simpleFilename := strings.TrimSuffix(filepath.Base(fileName), ext)
	languageMatch, _ := subtitleLanguageRegex.FindStringMatch(simpleFilename)

	if hasGroup(languageMatch, "iso_code") {
		return language.Make(getMatchGroupString(languageMatch, "iso_code")), nil
	}

	return language.Tag{}, errors.New("Unable to find a subtitle language")
}
