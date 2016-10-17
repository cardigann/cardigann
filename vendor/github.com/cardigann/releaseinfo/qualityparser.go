package releaseinfo

import (
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
)

var (
	sourceRegex = regexp2.MustCompile(
		`\b(?:(?<bluray>BluRay|Blu-Ray|HDDVD|BD)|(?<webdl>WEB[-_. ]DL|WEBDL|WebRip|iTunesHD|WebHD|[. ]WEB[. ](?:[xh]26[45]|DD5[. ]1)|\d+0p[. ]WEB[. ])|(?<hdtv>HDTV)|(?<bdrip>BDRip)|(?<brrip>BRRip)|(?<dvd>DVD|DVDRip|NTSC|PAL|xvidvd)|(?<dsr>WS[-_. ]DSR|DSR)|(?<pdtv>PDTV)|(?<sdtv>SDTV)|(?<tvrip>TVRip))\b`,
		regexp2.Compiled|regexp2.IgnoreCase)

	rawHDRegex = regexp2.MustCompile(
		`\b(?<rawhd>RawHD|1080i[-_. ]HDTV|Raw[-_. ]HD|MPEG[-_. ]?2)\b`,
		regexp2.Compiled|regexp2.IgnoreCase)

	versionRegex = regexp2.MustCompile(
		`\dv(?<version>\d)\b|\[v(?<version>\d)\]`,
		regexp2.Compiled|regexp2.IgnoreCase)

	properRegex = regexp2.MustCompile(
		`\b(?<real>real[-_. ])*(?<proper>proper|repack|rerip)(?<real>[-_. ]real)*\b`,
		regexp2.Compiled|regexp2.IgnoreCase)

	resolutionRegex = regexp2.MustCompile(
		`\b(?:(?<_480p>480p|640x480|848x480)|(?<_576p>576p)|(?<_720p>720p|1280x720)|(?<_1080p>1080p|1920x1080)|(?<_2160p>2160p))\b`,
		regexp2.Compiled|regexp2.IgnoreCase)

	codecRegex = regexp2.MustCompile(
		`\b(?:(?<x264>x264)|(?<h264>h264)|(?<xvidhd>XvidHD)|(?<xvid>Xvid)|(?<divx>divx))\b`,
		regexp2.Compiled|regexp2.IgnoreCase)

	otherSourceRegex = regexp2.MustCompile(
		`(?<hdtv>HD[-_. ]TV)|(?<sdtv>SD[-_. ]TV)`,
		regexp2.Compiled|regexp2.IgnoreCase)

	animeBlurayRegex = regexp2.MustCompile(
		`bd(?:720|1080)|(?<=[-_. (\[])bd(?=[-_. )\]])`,
		regexp2.Compiled|regexp2.IgnoreCase)

	highDefPdtvRegex = regexp2.MustCompile(`hr[-_. ]ws`,
		regexp2.Compiled|regexp2.IgnoreCase)
)

type sourceMatches map[string][]string

func (s sourceMatches) Has(key string) bool {
	return len(s[key]) > 0
}

func findSourceMatches(name string) sourceMatches {
	matches := sourceMatches{}
	sourceMatch, _ := sourceRegex.FindStringMatch(name)

	for sourceMatch != nil {
		for _, group := range sourceMatch.Groups() {
			switch group.Name {
			case "bluray", "webdl", "hdtv", "bdrip", "brrip", "dvd", "dsr", "pdtv", "sdtv", "tvrip":
				for _, capture := range group.Captures {
					if _, exists := matches[group.Name]; !exists {
						matches[group.Name] = []string{}
					}
					matches[group.Name] = append(matches[group.Name], capture.String())
				}
			}
		}
		sourceMatch, _ = sourceRegex.FindNextMatch(sourceMatch)
	}

	return matches
}

func ParseQuality(name string) QualityModel {
	// log.Printf("Parsing quality for %q", name)

	normalizedName := removeSpace(name)
	normalizedName = strings.Replace(normalizedName, "_", " ", -1)
	normalizedName = removeSpace(normalizedName)
	normalizedName = strings.ToLower(normalizedName)

	result := parseQualityModifiers(name, normalizedName)

	if match, _ := rawHDRegex.FindStringMatch(normalizedName); match != nil {
		result.Quality = QualityRAWHD
		return result
	}

	sourceMatches := findSourceMatches(normalizedName)
	resolution := ParseResolution(normalizedName)
	codecRegex, _ := codecRegex.FindStringMatch(normalizedName)

	if len(sourceMatches) > 0 {
		if sourceMatches.Has("bluray") {
			if codecRegex != nil {
				if hasGroup(codecRegex, "xvid") || hasGroup(codecRegex, "divx") {
					result.Quality = QualityDVD
					return result
				}
			}

			if resolution == Resolution2160p {
				result.Quality = QualityBluray2160p
				return result
			}

			if resolution == Resolution1080p {
				result.Quality = QualityBluray1080p
				return result
			}

			if resolution == Resolution480p || resolution == Resolution576p {
				result.Quality = QualityDVD
				return result
			}

			result.Quality = QualityBluray720p
			return result
		}

		if sourceMatches.Has("webdl") {
			if resolution == Resolution2160p {
				result.Quality = QualityWEBDL2160p
				return result
			}

			if resolution == Resolution1080p {
				result.Quality = QualityWEBDL1080p
				return result
			}

			if resolution == Resolution720p {
				result.Quality = QualityWEBDL720p
				return result
			}

			if strings.Contains(name, "[WEBDL]") {
				result.Quality = QualityWEBDL720p
				return result
			}

			result.Quality = QualityWEBDL480p
			return result
		}

		if sourceMatches.Has("hdtv") {
			if resolution == Resolution2160p {
				result.Quality = QualityHDTV2160p
				return result
			}

			if resolution == Resolution1080p {
				result.Quality = QualityHDTV1080p
				return result
			}

			if resolution == Resolution720p {
				result.Quality = QualityHDTV720p
				return result
			}

			if strings.Contains(name, "[HDTV]") {
				result.Quality = QualityHDTV720p
				return result
			}

			result.Quality = QualitySDTV
			return result
		}

		if sourceMatches.Has("bdrip") || sourceMatches.Has("brrip") {
			switch resolution {
			case Resolution720p:
				result.Quality = QualityBluray720p
				return result
			case Resolution1080p:
				result.Quality = QualityBluray1080p
				return result
			default:
				result.Quality = QualityDVD
				return result
			}
		}

		if sourceMatches.Has("dvd") {
			result.Quality = QualityDVD
			return result
		}

		if sourceMatches.Has("pdtv") ||
			sourceMatches.Has("sdtv") ||
			sourceMatches.Has("dsr") ||
			sourceMatches.Has("tvrip") {
			if match, _ := highDefPdtvRegex.FindStringMatch(normalizedName); match != nil {
				result.Quality = QualityHDTV720p
				return result
			}

			result.Quality = QualitySDTV
			return result
		}
	}

	//Anime Bluray matching
	if match, _ := animeBlurayRegex.FindStringMatch(normalizedName); match != nil {
		if resolution == Resolution480p || resolution == Resolution576p || strings.Contains(normalizedName, "480p") {
			result.Quality = QualityDVD
			return result
		}

		if resolution == Resolution1080p || strings.Contains(normalizedName, "1080p") {
			result.Quality = QualityBluray1080p
			return result
		}

		result.Quality = QualityBluray720p
		return result
	}

	if resolution == Resolution2160p {
		result.Quality = QualityHDTV2160p
		return result
	}

	if resolution == Resolution1080p {
		result.Quality = QualityHDTV1080p
		return result
	}

	if resolution == Resolution720p {
		result.Quality = QualityHDTV720p
		return result
	}

	if resolution == Resolution480p {
		result.Quality = QualitySDTV
		return result
	}

	if codecRegex != nil && hasGroup(codecRegex, "x264") {
		result.Quality = QualitySDTV
		return result
	}

	if strings.Contains(normalizedName, "848x480") {
		if strings.Contains(normalizedName, "dvd") {
			result.Quality = QualityDVD
		}

		result.Quality = QualitySDTV
	}

	if strings.Contains(normalizedName, "1280x720") {
		if strings.Contains(normalizedName, "bluray") {
			result.Quality = QualityBluray720p
		}

		result.Quality = QualityHDTV720p
	}

	if strings.Contains(normalizedName, "1920x1080") {
		if strings.Contains(normalizedName, "bluray") {
			result.Quality = QualityBluray1080p
		}

		result.Quality = QualityHDTV1080p
	}

	if strings.Contains(normalizedName, "bluray720p") {
		result.Quality = QualityBluray720p
	}

	if strings.Contains(normalizedName, "bluray1080p") {
		result.Quality = QualityBluray1080p
	}

	if otherSourceMatch := otherSourceMatch(normalizedName); otherSourceMatch != QualityUnknown {
		result.Quality = otherSourceMatch
	}

	//Based on extension
	if result.Quality == QualityUnknown {
		result.Quality = getQualityForExtension(filepath.Ext(normalizedName))
		result.QualitySource = "extension"
	}

	return result
}

func ParseResolution(name string) Resolution {
	match, _ := resolutionRegex.FindStringMatch(name)

	switch {
	case hasGroup(match, "_480p"):
		return Resolution480p
	case hasGroup(match, "_576p"):
		return Resolution576p
	case hasGroup(match, "_720p"):
		return Resolution720p
	case hasGroup(match, "_1080p"):
		return Resolution1080p
	case hasGroup(match, "_2160p"):
		return Resolution2160p
	}

	return ResolutionUnknown
}

func otherSourceMatch(name string) Quality {
	match, _ := otherSourceRegex.FindStringMatch(name)

	switch {
	case hasGroup(match, "sdtv"):
		return QualitySDTV
	case hasGroup(match, "hdtv"):
		return QualityHDTV720p
	}

	return QualityUnknown
}

func parseQualityModifiers(name string, normalizedName string) QualityModel {
	result := QualityModel{Quality: QualityUnknown, QualitySource: "name"}

	versionRegexResult, _ := versionRegex.FindStringMatch(normalizedName)
	if versionRegexResult != nil {
		version, _ := strconv.Atoi(versionRegexResult.GroupByName("version").String())
		result.Revision = version - 1
	}

	properRegexResult, _ := properRegex.FindStringMatch(normalizedName)
	if properRegexResult != nil {
		result.Revision += len(properRegexResult.GroupByName("proper").Captures)
		result.Revision += len(properRegexResult.GroupByName("real").Captures)
	}

	return result
}

type Resolution string

const (
	Resolution480p    Resolution = "480p"
	Resolution576p    Resolution = "576p"
	Resolution720p    Resolution = "720p"
	Resolution1080p   Resolution = "1080p"
	Resolution2160p   Resolution = "2160p"
	ResolutionUnknown Resolution = "unknown"
)
