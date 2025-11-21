package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)



func main (){
	args := os.Args
	if len(args) > 1 && args[1] == "--scope"{
		fmt.Print("in_line")
		os.Exit(0)
	}
	if len(args) > 1 && args[1] == "--supplies"{
		fmt.Print("link")
		os.Exit(0)
	}

	inputBytes, err := io.ReadAll(os.Stdin)
   if err != nil {
       log.Fatal(err)
   }
	input := string(inputBytes)
	out := input

	prefix := ""
	for ; strings.Contains(input,"@link") ; {
		prefix = "+"
		index := strings.Index(input,"@link")
		out = input[0:index]
		link := ""
		text := ""
		readingLink := true
		readingText:= false
		offset := index + strLen("@link:")
		for i := offset; i < strLen(input); i++{
			if (readingLink){
				if input[i] == '('{
					readingLink = false
					readingText = true
					continue
				}
				link += string(input[i])
			}
			if (readingText){
				if input[i] == ')'{
					readingText = false
					offset = i+1
					break
				}
				text += string(input[i])

			}
		}
		
		out += fmt.Sprintf("<a href=\"%s\">%s</a>",link,text)
		out += input[offset:]
		input = out

	}
	
	fmt.Printf("%s%s",prefix,out)
}

func strLen(src string) int{
	return len([]rune(src))
}
