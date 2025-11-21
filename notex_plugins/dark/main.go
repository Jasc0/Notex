package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

const BG_INIT int = 0x10
const SECTION_SCALER float32 = 1.5
const FG_COLOR = "#f9e0a1"

func main(){
	args := os.Args
	if len(args) > 1 && args[1] == "--scope"{
		fmt.Print("document")
		os.Exit(0)
	}
	if len(args) > 1 && args[1] == "--supplies"{
		fmt.Print("dark")
		os.Exit(0)
	}

	inputBytes, err := io.ReadAll(os.Stdin)
   if err != nil {
       log.Fatal(err)
   }
	tokens := []string{}

	// Compile regex to match text between <token> and </token>
   re := regexp.MustCompile(`<token>(.*?)</token>`)

   // Find all matches
   matches := re.FindAllStringSubmatch(string(inputBytes), -1)

   // Extract only the inner text (group 1)
   for _, match := range matches {
       if len(match) > 1 {
           tokens = append(tokens, match[1])
       }
   }
	
	replaceIndex := -1
	maxDepth:= 0
	curDepth:= 0
	for i := 0; i < len(tokens); i++ {
		if strings.Contains(tokens[i], "@dark"){
			if strings.HasPrefix(tokens[i],"\\"){
				continue
			}
			replaceIndex = i
		}
		if tokens[i] == "{"{
			curDepth++
		}
		if tokens[i] == "}"{
			curDepth--
		}
		if curDepth > maxDepth {
			maxDepth = curDepth
		}
		
	}
	if replaceIndex == -1{
		fmt.Print(string(inputBytes))
		os.Exit(0)
	}
	parts := strings.SplitN(tokens[replaceIndex],"=",2)
	rep := parts[0] + "="
	rep += fmt.Sprintf(` body {
		background-color: #%x%x%x;
		color: %s;
	}
	p {
		background-color : inherit;
	}
	img {
		height : 400px;
	}
	h1 {
		color: #326416;
	}
	`, BG_INIT,BG_INIT,BG_INIT,FG_COLOR)

	rep += generateSections(maxDepth)

	replaced := tokens
	replaced[replaceIndex] = rep

	for _, tok := range replaced{
		fmt.Printf("<token>%s</token>", tok)
	}

	
}

func generateSections (depth int) string {
	ret := ""

	for i := 1 ; i <= depth; i++{
		scaled := float32(BG_INIT) * SECTION_SCALER
		background := (int(scaled) * i )% 255
		ret += fmt.Sprintf(".section_depth\\=%d {\n", i)
		ret += fmt.Sprintf("background-color : #%x%x%x;\n",background,background,background)
		ret += "margin-top : 2px;\n"
		ret += "margin-bottom : 2px;\n"
		ret += "margin-left : 0px;\n"
		ret += "padding-left : 0px;\n"
		ret += "padding-top : 2px;\n"
		ret += "padding-bottom : 2px;\n"
		ret += "}\n"
	}
	return ret
}
