package main

import (
	"os"
)


func tokenizeFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil{
		return []string{}, err
	}
	
	contents := string(data)
	return tokenizeString(contents), nil

}

func tokenizeString(source string) []string {
	runes := []rune(source)
	tokens := []string{}
	token := ""
	for i := 0; i < len(runes); i++{
		switch(runes[i]){
		case '!':
			if len([]rune(token)) > 0{
				token += string(runes[i])
			}
			for ; runes[i] != '\n'; i++{
				token += string(runes[i])
			}
			tokens = append(tokens, token)
			token = ""
		case '#':
			if len([]rune(token)) > 0{
				tokens = append(tokens, token)
			}
			token = "#"
			for ;  i+1 < len(runes) && runes[i+1] == '#'; i++{
				token += "#"
			}
			tokens = append(tokens, token)
			token = ""

		case ' ', '\t':
			if len([]rune(token)) > 0{
				token += " "
			}

		case '\n':
					tokens = append(tokens, token)
					token = ""
		case '%':
			if i-1 >= 0 && runes[i-1] != '\\'{
				for ; i+1 < len(runes) && runes[i+1] != '\n' ; i++ {
					continue
				}
			} 
		
		case '\\':
			if len(token) > 0{
				token += string(runes[i+1])
				continue
			}

			for ; runes[i] != '\n' ; i++ {
				token += string(runes[i])
			}
			tokens = append(tokens, token, "") // extra "" for new line 
			token = ""

		case '{', '}', '.', '-':
			if len([]rune(token)) == 0 {
				tokens = append(tokens,  string(runes[i]))
				token = ""
			} else{
				token += string(runes[i])
			}

		case '/':
			if i+1 < len(runes) && runes[i+1] == '@'{

				if len([]rune(token)) > 0{
					tokens = append(tokens, token)
				}
					token = ""
			} else{
				token += string(runes[i])
			}

		default:
			token = token + string(runes[i])
		}

	}
	tokens = append(tokens, token)

	return tokens
}

