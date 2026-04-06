package main

import (
	"bufio"
	"log"
	"strings"
	"fmt"
)

type NotexHead struct{
	styles map[string]string;
	title, description, raw string;
}

func (nh NotexHead) String() (string, error){
	var strBuilder strings.Builder
	strBuilder.WriteString("<head>\n")
	strBuilder.WriteString("<meta charset=\"UTF-8\">\n")
	strBuilder.WriteString("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	fmt.Fprintf(&strBuilder, "<title>%s</title>\n", nh.title)
	if len(nh.description) > 0 {
		fmt.Fprintf(&strBuilder,"<meta name=\"description\" content=\"%s\">\n", nh.description)
	}
	for _, style := range nh.styles {
		fmt.Fprintf(&strBuilder,"<style>%s</style>\n", style)
	}
	if len(nh.raw) > 0 {
		strBuilder.WriteString(nh.raw + "\n")
	}
	strBuilder.WriteString("</head>")
	return strBuilder.String(), nil
}


func separateHeader(inp string) (body, head string){
	if !strings.HasPrefix(strings.TrimLeft(inp, " \t\n"), "["){
		return inp, ""
	}
	parts := strings.SplitN(inp, "]", 2)
	if len(parts) < 2{
		return inp, ""
	}

	if strings.HasPrefix(parts[0], "["){
		parts[0] = strings.TrimLeft(parts[0], "[")
	} 
	return parts[1], parts[0]
}

func handleHead(inp string) NotexHead{
	head := NotexHead{title: "Notex Document", description: ""}
	scanner := bufio.NewScanner(strings.NewReader(inp))
	for scanner.Scan(){
		line := strings.Trim(scanner.Text(), " \t")
		if len(line) == 0{
			continue
		}
		line_runes := []rune(line)
		switch line_runes[0]{
		case '#':
			head.title = line[1:]
			// title
		case '!':
			log.Fatal(fmt.Errorf("Attributes not allowed in header"))
		case '<':
			head.raw += line + "\n"
			
		default:
			head.description += line + " "
		}
	}


	return head
}
