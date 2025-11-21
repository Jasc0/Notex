package main

import (
	"fmt"
	"log"
	"os"
	"io"
	"strings"
	"errors"
)

func main (){
	args := os.Args
	if len(args) > 1 && args[1] == "--scope"{
		fmt.Print("single_line")
		os.Exit(0)
	}
	if len(args) > 1 && args[1] == "--supplies"{
		fmt.Print("img")
		os.Exit(0)
	}

	inputBytes, err := io.ReadAll(os.Stdin)
   if err != nil {
       log.Fatal(err)
   }
	input := string(inputBytes)
	// images
	out := input
	if strings.HasPrefix(input, "@img"){
		out = "\\"
	   parts := strings.SplitN(input,":", 2)
	   if len(parts) < 2{
	  	 log.Fatal(errors.New("Image missing src"))
	   }
	   parts = strings.Split(parts[1], ",")
	   image_src := parts[0]
	   out += "\n<center>\n <img src=\"" + image_src +"\""
	   for j := 1; j < len(parts); j++ {
	  	 out += " " + parts[j]
	   }
	   out += ">\n</center>\n"

	}
	fmt.Printf(out)
}
