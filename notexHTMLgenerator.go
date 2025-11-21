package main

import (
	"fmt"
	"strings"
	"unicode"
)

func generateDefaultCSS(sectionDepth int) string {
	ret := ""
	ret += ` body {
		background-color: #101010;
		color: #f9e0a1;
	}
	p {
		background-color : inherit;
	}
	`
	
	for i := 1; i <= sectionDepth; i++{
		red := int(10*3.5*i)
		green := int(10*3.5*i)
		blue := int(10*3.5*i)

		ret += fmt.Sprintf(".section_depth\\=%d {\n", i)
		ret += fmt.Sprintf("background-color : #%d%d%d;\n",red,green,blue)
		ret += "margin : 2px;"
		ret += "padding : 2px;"

		ret += "}\n"

	}
	return ret
}

func getMaxSectionDepth(tokens []string) int{
	maxDepth := 0
	curDepth := 0
	for _, token := range tokens{
		if token == "{"{
			curDepth++
		}
		if token == "}"{
			curDepth--
		}
		if maxDepth < curDepth{
			maxDepth = curDepth
		}
	}
	return maxDepth
}

func generateHTMLHead(tokens []string) string{
	varMap := make(map[string]string)
	for _, token := range tokens{
		if len(token) > 0 && token[0] == '!'{
			parts := strings.SplitN(token, "=", 2)	
			key := parts[0][1:]
			value := parts[1]
			varMap[key] = value
		}
	}
	ret := "<!DOCTYPE html> <html lang=en> <head>\n"
	
	// head

	head_val, ok := varMap["head"]
	if ok{
		ret += head_val
	}

	// style
	style_val, ok := varMap["style"]
	if !ok{
		style_val = generateDefaultCSS(getMaxSectionDepth(tokens))
	}else{
		ret += "\n<style>\n"
		ret += style_val
		ret += "\n</style>\n"
	}
	ret += "\n</head>\n"
	return ret

}

func generateHTMLBody(tokens []string) string {
	ret := "<body>\n"

	sectionDepth:= 0
	openPTag := false

	for i := 0; i < len(tokens); i++ {

		if len(tokens[i]) > 1 && tokens[i][0] == '!'{
			continue
		}

		// always append
		if (len(tokens[i]) > 0 && tokens[i][0] == '+'){
			ret += tokens[i][1:]
			continue
		}

		// paragraphs
		if startsAlphanumeric(string(tokens[i])) {
			if openPTag{
				ret += " " + tokens[i]  
				continue
			}else{
				ret += "<p>\n" + tokens[i]
				openPTag = true
			}
		} else if openPTag{
			ret += "\n</p>\n"
			openPTag = false
		}


		// headings
		if strings.HasPrefix(tokens[i], "#"){
			headingLevel := len([]rune(tokens[i]))
			ret += fmt.Sprintf("<h%d>%s</h%d>\n",headingLevel,tokens[i+1],headingLevel)
			i++ // subsequent token already handled
			continue
		}


		// lists
		if tokens[i] == "-" {
			ret += "<ul>\n"
			for ; tokens[i] == "-"; i += 2 {
				ret += fmt.Sprintf("<li> %s </li>\n", tokens[i+1])
			}
			i--
			ret += "</ul>\n"
			continue
		}

		if tokens[i] == "." {
			ret += "<ol>\n"
			for ; tokens[i] == "."; i += 2 {
				ret += fmt.Sprintf("<li> %s </li>\n", tokens[i+1])
			}
			ret += "</ol>\n"
			i--
			continue
		}


		//sections
		if tokens[i] == "{" {
			sectionDepth++
			ret += fmt.Sprintf("<div class=\"section_depth=%d\">\n",sectionDepth)
			continue
		}

		if tokens[i] == "}" {
			ret += "</div>"
			sectionDepth--
			continue
		}

		// don't handles
		if strings.HasPrefix(tokens[i],"\\") {
			nonWSIndex := 1
			for ; unicode.IsSpace(rune(tokens[i][nonWSIndex])); nonWSIndex++{continue}
			if tokens[i][nonWSIndex] != '<' && !openPTag{
				openPTag = true
				ret += "\n<p>\n"
			}
			ret += tokens[i][1:]
			if !openPTag{
				ret += "\n"
			}
			continue
		}
	} 

	ret += "</body>\n</html>"
	return ret
}

func startsAlphanumeric(s string) bool {
    if len(s) == 0 {
        return false // Empty string case
    }
    firstChar := rune(s[0]) // Convert to rune (char)
    return unicode.IsLetter(firstChar) || unicode.IsDigit(firstChar) || s[0] ==  '"' || s[0] == '\''
}
