package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

func main (){
	args := os.Args
	if len(args) > 1 && args[1] == "--scope"{
		fmt.Print("in_line")
		os.Exit(0)
	}
	if len(args) > 1 && args[1] == "--supplies"{
		fmt.Print("now")
		os.Exit(0)
	}

	inputBytes, err := io.ReadAll(os.Stdin)
   if err != nil {
       log.Fatal(err)
   }
	input := string(inputBytes)
	out := input


	if strings.Contains(input,"/@now"){
		out = input[0:strings.Index(input,"/@now")]
		out += time.Now().Format(time.RFC1123Z)
		offset := strings.Index(input,"/@now") + utf8.RuneCountInString("/@now")
		out += input[offset:]

	}
	fmt.Printf(out)
}
