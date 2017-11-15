package util

import (
	"html"
	"regexp"
	"strings"

	"mvdan.cc/xurls"
)

var paragraph = regexp.MustCompile(`\n+(.*)\n`)
var codeBlock = regexp.MustCompile("^```(.*)^```")
var digitLink = regexp.MustCompile(`&gt;&gt;([0-9]+)`)

// Constructs raw html based on a user message. The output must be escaped.
//
// - Blank lines separate paragraphs
//
// - ``` is used to fence code blocks
//
// - Urls become links
//
// - Post id's prefixed with `>>` become thread local links
//
// TODO: This is slow and could be done in a single pass.
func FormatPost(msg string) string {
	split := strings.Split(msg, "\n")

	// Blank lines to paragraphs
	var pp string = ""
	for _, s := range split {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		pp = pp + "<p>" + html.EscapeString(s) + "</p>"
	}

	msg = pp

	// TODO: Code blocks require a pass at the same time as paragraph checking

	// Urls become links
	urls := xurls.Strict().FindAllString(msg, -1)
	for _, url := range urls {
		msg = strings.Replace(msg, url, `<a href="`+url+`">`+url+`</a>`, -1)
	}

	// Replace >>:digit with a link.
	msg = digitLink.ReplaceAllString(msg, `<a href="https://gorum.xyz/post/$1">&gt;&gt;$1</a>`)

	return msg
}
