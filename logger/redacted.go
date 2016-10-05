package logger

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	redactedMarker = "r̶e̶d̶a̶c̶t̶e̶d̶"
)

var (
	redactedRegexps = []*regexp.Regexp{
		regexp.MustCompile("(?i)(torrent_pass|pass|authkey|token|apikey)(=)([^&$]+)"),
		regexp.MustCompile("(?i)(cookie|password)(:)([^\\s$\\]]+)"),
	}
)

type redactedLogFormatter struct {
	logrus.Formatter
	secrets      map[string]struct{}
	secretsRegex *regexp.Regexp
}

func (f *redactedLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.secrets == nil {
		f.secrets = map[string]struct{}{}
	}

	for k, v := range entry.Data {
		s := fmt.Sprintf("%v", v)

		for _, re := range redactedRegexps {
			var found bool
			if matches := re.FindAllStringSubmatch(s, -1); len(matches) > 0 {
				for _, match := range matches {
					if match[2] == redactedMarker {
						continue
					}
					s = strings.Replace(s, match[0], match[1]+match[2]+redactedMarker, -1)
					f.secrets[match[3]] = struct{}{}
					found = true
				}
			}
			if found {
				entry.Data[k] = s
				if err := f.updateRegexp(); err != nil {
					return nil, err
				}
			}
		}

		if f.secretsRegex != nil && f.secretsRegex.MatchString(s) {
			s = f.secretsRegex.ReplaceAllLiteralString(s, redactedMarker)
			entry.Data[k] = s
		}
	}

	return f.Formatter.Format(entry)
}

func (f *redactedLogFormatter) updateRegexp() error {
	quoted := []string{}
	for secret := range f.secrets {
		quoted = append(quoted, regexp.QuoteMeta(secret))
	}

	re, err := regexp.Compile(fmt.Sprintf("(?i)(%s)", strings.Join(quoted, "|")))
	if err != nil {
		return err
	}

	f.secretsRegex = re
	return nil
}
