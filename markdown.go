package clerk

import (
	"io/ioutil"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// \W - not a word character
// \@ - as literal
// \p{L}+ - at least one Unicode word character, greedily
var userRegexp = regexp.MustCompile(`(\W)(\@\p{L}+)`)

const colorHighlighting = `$1<span style="font-weight: bold; color: #e34234;">$2</span>`

func markdownToHTML(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	unsafe := blackfriday.MarkdownCommon(data)
	html := string(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
	html = highlightUsers(html)

	return html, nil
}

func highlightUsers(html string) string {
	return userRegexp.ReplaceAllString(html, colorHighlighting)
}
