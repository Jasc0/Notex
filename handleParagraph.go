package main

import (
	"fmt"
	"html"
	"regexp"
	"strings"
)

var htmlTagRe = regexp.MustCompile(`(<[^>]+>|&lt;|&gt;)`)

func renderText(text string) string {
	parts := htmlTagRe.Split(text, -1)
	tags := htmlTagRe.FindAllString(text, -1)
	var sb strings.Builder
	for i, part := range parts {
		sb.WriteString(html.EscapeString(part))
		if i < len(tags) {
			sb.WriteString(tags[i])
		}
	}
	return sb.String()
}

type NotexPara struct{
	text string
}
func (p NotexPara) String() (string, error) {
    return fmt.Sprintf("<p>%s</p>", renderText(p.text)), nil
}

func handlePara (inp_str string) NotexPara{
	var ret NotexPara
	ret.text = inp_str
	return ret
}
