package releaseinfo

import (
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
)

var fileExtensionRegex = regexp2.MustCompile(`\.[a-z0-9]{2,4}$`,
	regexp2.IgnoreCase|regexp2.Compiled)

func removeFileExtension(title string) string {
	result, err := fileExtensionRegex.ReplaceFunc(title, func(m regexp2.Match) string {
		ext := strings.ToLower(filepath.Ext(m.String()))
		if _, match := mediaFileExtensions[ext]; match {
			return ""
		}
		return m.String()
	}, 0, -1)

	if err != nil {
		return title
	}

	return result
}

func normalizePath(s string) string {
	return strings.Map(func(r rune) rune {
		if r == '/' || r == '\\' {
			return filepath.Separator
		}
		return r
	}, s)
}

var mediaFileExtensions = map[string]Quality{
	".webm":   QualityUnknown,
	".m4v":    QualitySDTV,
	".3gp":    QualitySDTV,
	".nsv":    QualitySDTV,
	".ty":     QualitySDTV,
	".strm":   QualitySDTV,
	".rm":     QualitySDTV,
	".rmvb":   QualitySDTV,
	".m3u":    QualitySDTV,
	".ifo":    QualitySDTV,
	".mov":    QualitySDTV,
	".qt":     QualitySDTV,
	".divx":   QualitySDTV,
	".xvid":   QualitySDTV,
	".bivx":   QualitySDTV,
	".nrg":    QualitySDTV,
	".pva":    QualitySDTV,
	".wmv":    QualitySDTV,
	".asf":    QualitySDTV,
	".asx":    QualitySDTV,
	".ogm":    QualitySDTV,
	".ogv":    QualitySDTV,
	".m2v":    QualitySDTV,
	".avi":    QualitySDTV,
	".bin":    QualitySDTV,
	".dat":    QualitySDTV,
	".dvr-ms": QualitySDTV,
	".mpg":    QualitySDTV,
	".mpeg":   QualitySDTV,
	".mp4":    QualitySDTV,
	".avc":    QualitySDTV,
	".vp3":    QualitySDTV,
	".svq3":   QualitySDTV,
	".nuv":    QualitySDTV,
	".viv":    QualitySDTV,
	".dv":     QualitySDTV,
	".fli":    QualitySDTV,
	".flv":    QualitySDTV,
	".wpl":    QualitySDTV,
	".img":    QualityDVD,
	".iso":    QualityDVD,
	".vob":    QualityDVD,
	".mkv":    QualityHDTV720p,
	".ts":     QualityHDTV720p,
	".wtv":    QualityHDTV720p,
	".m2ts":   QualityBluray720p,
}

func getQualityForExtension(ext string) Quality {
	q, ok := mediaFileExtensions[ext]
	if !ok {
		return QualityUnknown
	}
	return q
}
