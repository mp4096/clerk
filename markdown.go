package clerk

import (
	"io/ioutil"
	"regexp"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// \b - word boundary
// \@ - as literal
// \w+ - at least one word character, greedily
var userRegexp = regexp.MustCompile(`(\W)(\@\w+)`)

const colorHighlighting = `$1<b><font color="#ff9933">$2</font></b>`

func MarkdownToHTML(filename string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	unsafe := blackfriday.MarkdownCommon(data)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	html = []byte(highlightUsers(string(html)))

	return html, nil
}

func highlightUsers(html string) string {
	return userRegexp.ReplaceAllString(html, colorHighlighting)
}
