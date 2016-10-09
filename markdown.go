package clerk

import (
	"io/ioutil"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func MarkdownToHTML(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	unsafe := blackfriday.MarkdownCommon(data)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	return html, nil
}
