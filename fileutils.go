package clerk

import (
	"errors"
	"io/ioutil"
)

func FindLatestFile() (string, error) {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		return "", err
	}

	latestTimestamp := "0000-00-00"
	latestFilename := ""

	for _, f := range files {
		timestamp, _ := dateFromMDFilename(f.Name())
		if timestamp > latestTimestamp {
			latestTimestamp = timestamp
			latestFilename = f.Name()
		}
	}

	if len(latestFilename) == 0 {
		return "", errors.New("Could not find a suitable file")
	}
	return latestFilename, nil
}
