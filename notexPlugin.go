package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

// Plugins are executable files that will accept any token that starts with an @ as input
// if the plugin supplies a matching implementation for the string following the @ symbol
// it will write to stdout the text that will take the place of the token as input before being sent
// to be compiled, if the plugin does not implement the text, it will write the original token to stdout

func plugNameAcceptable(r rune) bool{
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-'
}

func getUsedPlugins(tokens []string) []string {
	var ret []string

	for _ , token := range tokens{
		if strings.Contains(token, "@") || (strings.HasPrefix(token,"!") && strings.Contains(token,"@")){
			plug_part := []rune(token[strings.Index(token,"@")+1:])
			valid := make([]rune, 0, len(plug_part))
			for _, r := range plug_part{
				if plugNameAcceptable(r){
					valid = append(valid, r)
				} else{
					break
				}
			}
			plug_name := string(valid)
			plug_name = strings.TrimRight(plug_name,"-")

			if !contains(ret, plug_name){
				ret = append(ret, plug_name)
			}
		}
	}

	return ret
}

func getUsedPluginLines(tokens []string, name string) []int{
	var ret []int
	for i, token := range tokens{
		if strings.Contains(token, fmt.Sprintf("@%s",name) ){
			ret = append(ret, i)
		}
	}
	return ret
}


func runPlugin(tokens []string, path string, usedOn []int) []string{

	output := []string{}

	var stdout bytes.Buffer
	cmd := exec.Command(path, "--scope")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
				log.Fatal(err)
	}

	switch (stdout.String()){

	case "single_line":
		stdout.Reset()
		for _, index := range usedOn{
			token := tokens[index]
			if strings.HasPrefix(token, "@"){
				cmd = exec.Command(path)
				stdin := bytes.NewBufferString(token)

				cmd.Stdin = stdin
				cmd.Stdout = &stdout

				if err := cmd.Run(); err != nil {
					log.Fatal(err)
				}
				
				//output = append(output, stdout.String())
				tokens[index] = stdout.String()


			} 
			stdout.Reset()
		}
	case "in_line":
		stdout.Reset()
		for _, index := range usedOn{
			token := tokens[index]
			if strings.Contains(token, "@") && !strings.Contains(token, "\\/@"){
				cmd = exec.Command(path)
				stdin := bytes.NewBufferString(token)

				cmd.Stdin = stdin
				cmd.Stdout = &stdout

				if err := cmd.Run(); err != nil {
					log.Fatal(err)
				}
				
				//output = append(output, stdout.String())
				tokens[index] = stdout.String()

			} 
			stdout.Reset()
		}

	// must handle escaped \/@ itself
	case "document":
		stdout.Reset()
		input := ""
		for _, token := range tokens{
			input += "<token>" + token + "</token>"
		}
			cmd = exec.Command(path)
			stdin := bytes.NewBufferString(input)

			cmd.Stdin = stdin
			cmd.Stdout = &stdout

			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}

			// Compile regex to match text between <token> and </token>
    		re := regexp.MustCompile(`(?s)<token>(.*?)</token>`)

    		// Find all matches
    		matches := re.FindAllStringSubmatch(stdout.String(), -1)

    		// Extract only the inner text (group 1)
    		for _, match := range matches {
    		    if len(match) > 1 {
    		        output = append(output, match[1])
    		    }
    		}
			tokens = output
		}
	
return tokens

}

func getPluginName(path string) string{
	var stdout bytes.Buffer
	cmd := exec.Command(path, "--supplies")
	cmd.Stdout = &stdout
	ret := stdout.String()
	stdout.Reset()
	return ret
}

func getPluginNameLength(path string) int{
	str := getPluginName(path)
	return len([]rune(str))
}

func getSupplies(path string) string{
	var stdout bytes.Buffer
	cmd := exec.Command(path, "--supplies")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
				log.Fatal(err)
	}
	return stdout.String()
}
