package clerk

import (
	"errors"
	"regexp"
)

var isoDateRegexp = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}`)

func dateFromMDFilename(filename string) (string, error) {
	m := isoDateRegexp.FindString(filename)
	if m == "" {
		return "", errors.New("Markdown filename does not begin with a timestamp")
	}
	return string(m), nil
}
