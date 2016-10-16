package releaseinfo

import "strings"

//This method should prefer false negatives over false positives.
//It's better not to use a title that might be scene than to use one that isn't scene
func IsSceneTitle(title string) (bool, error) {
	// log.Printf("Checking if %q is a scene title", title)

	if !strings.Contains(title, ".") {
		return false, nil
	}

	if strings.Contains(title, " ") {
		return false, nil
	}

	parsedTitle, err := Parse(title)
	if err != nil {
		return false, err
	}

	if parsedTitle == nil {
		return false, err
	}

	if parsedTitle != nil &&
		parsedTitle.ReleaseGroup != "" &&
		parsedTitle.Quality.Quality != QualityUnknown &&
		!isNullOrWhiteSpace(parsedTitle.SeriesTitle) {
		return true, nil
	}

	return false, nil
}
