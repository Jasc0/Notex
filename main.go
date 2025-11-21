package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

func main() {

	args := os.Args

	if len(args) == 1 {
		printHelp()
		os.Exit(0)
	}

	writeToFile   := false
	readFromFile  := false
	printToSTDOUT := false
	readFromSTDIN := true
	readFrom      := ""
	writeTo		  := ""
	printTokens	  := false
	for i := 1; i < len(args); i++{
		switch args[i] {
		case "-i", "--input":
			readFromFile = true
			readFromSTDIN = false
			if i+1 < len(args) && !strings.HasPrefix(args[i+1],"-"){
				readFrom     = args[i+1]	
			} else{
				log.Fatal(errors.New("No/invalid argument for input"))
			}

			i++
		case "-o", "--output":
			writeToFile = true
			if i+1 < len(args) && !strings.HasPrefix(args[i+1],"-"){
				writeTo     = args[i+1]
			} else{
				log.Fatal(errors.New("No/invalid argument for output"))
			}
			i++
		case "-p", "--print":
			printToSTDOUT = true

		case "-t", "--tokens":
			printTokens = true

		}
	}

	var tokens []string
	if readFromFile{
		Temptokens, err := tokenizeFile(readFrom)		
		if err != nil{
			log.Fatal(err)
		}
		tokens = Temptokens
	} else if readFromSTDIN{
		inputBytes, err := io.ReadAll(os.Stdin)
   	if err != nil {
   	    log.Fatal(err)
   	}
		tokens = tokenizeString(string(inputBytes))
	} else{
		log.Fatal(errors.New("Unknown error reading input"))
	}

	if printTokens {
		fmt.Println("\n\nBEFORE PLUGINS")
		for i, t := range tokens{
			fmt.Printf("%d: %s\n", i, t)
		}
	}
	plugins := make(map[string]string)
	var plugins_order []string
	path := os.Getenv("NOTEX_PATH")
	for dir := range strings.SplitSeq(path,":"){
		files, err := os.ReadDir(dir)
		if err != nil{
			log.Fatal(err)
		}
		for _, file := range files {
			if !file.IsDir(){
				full_path := dir+"/"+ file.Name()
				supplies:= getSupplies(full_path)
				plugins[supplies] = full_path
				if !contains(plugins_order,supplies){
					plugins_order = append(plugins_order, getSupplies(full_path))
				}
			}

		}
	}
	// order plugins (longest name to shortest)
	sort.Slice(plugins_order, func(i, j int) bool {
		return getPluginNameLength(plugins_order[i]) > getPluginNameLength(plugins_order[j])
	})
	fmt.Printf("PLUGINS ORDER: %r" ,plugins_order)
	usedPlugins := getUsedPlugins(tokens)
	for _, up := range getUsedPlugins(tokens){
		fmt.Printf("%s : %r \n", up, getUsedPluginLines(tokens, up))
	}


	for _, p := range usedPlugins{
		if !contains(plugins_order, p){
			log.Printf("Warning: plugin %s referenced but not supplied, skipping", p)
		}
	}

	for _ , p := range plugins_order{
		// skip unused plugins
		if !contains(usedPlugins, p){
			continue
		}

		tokens = runPlugin(tokens, plugins[p], getUsedPluginLines(tokens, p))

	}

	if printTokens {
		fmt.Println("\n\nAFTER PLUGINS")
		for i, t := range tokens{
			fmt.Printf("%d: %s\n", i, t)
		}
	}

	if printToSTDOUT {
		fmt.Print(generateHTMLHead(tokens))
		fmt.Print(generateHTMLBody(tokens))
	}

	if writeToFile {
		outputWrite(tokens, writeTo)
	}

	


}

func printHelp(){
	fmt.Println(`Usage: notex -i example.ntx -o example.html
	or: cat example.ntx | notex -o example.html
	or: notex -i example.ntx -p > example.html
	
	Arguments:

	--input, -i 	Supplies .ntx file to be used for input, 
	takes precedence over stdin input if that is also supplied
	
	--output, -o	Provides file to be written to

	--print, -P 	Prints output to stdout instead of writing to a file 

	--tokens	Prints tokens parsed before and after plugin substitutions
	`)
}

func outputWrite (tokens []string, path string){
	writestr := ""
	writestr += generateHTMLHead(tokens) + "\n"
	writestr += generateHTMLBody(tokens) + "\n"

	err := os.WriteFile(path, []byte(writestr), 0644)
   if err != nil {
       log.Fatal("Error writing output file", err)
   }

	if err != nil{
	   log.Fatal(err)
	}
}

func contains(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

