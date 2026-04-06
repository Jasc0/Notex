package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"unicode"
)


type NotexFunction struct {

	name string
	parameters []string
	object string
	start, end int

}

type FunctionParseState int 

const (
	FuncNormalState FunctionParseState = iota
	FuncIgnoreState // \ before function
	FuncPotentialFunction // /
	FuncBeginFunction // @ after / 
	FuncParseName
	FuncParseParams
	FuncParseQuotedParam
	FuncParseQuotedParamEscape
	FuncParseObject
)


func validFuncName(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_'
}
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

func scanFunction( inp_str string) ([]NotexFunction, error){
	NFs := []NotexFunction{}
	cur_index := 0
	state := FuncNormalState

	var cur_func* NotexFunction
	var strBuilder strings.Builder 
	groupDepth := -1

	inp := []rune(inp_str)

	for ;cur_index < len(inp);{
		c := inp[cur_index]
		switch(state){
		case FuncNormalState:
			switch(c){
			case '/':
				state = FuncPotentialFunction
			case '\\':
				state = FuncIgnoreState
			}
		case FuncIgnoreState:
			state = FuncNormalState
		case FuncPotentialFunction:
			if (c == '@'){
				state = FuncBeginFunction
			} else{
				state = FuncNormalState
			}
		case FuncBeginFunction:
			cur_func = new(NotexFunction)
			cur_func.start = cur_index - 2
			strBuilder.Reset()
			strBuilder.WriteRune(c)
			state = FuncParseName
			groupDepth = -1
		case FuncParseName:
			if validFuncName(c){
				strBuilder.WriteRune(c)
			} else if (c == '('){
				state = FuncParseParams
				cur_func.name = strBuilder.String()
				strBuilder.Reset()
			} else if (c == '{'){
				cur_index--
				state = FuncParseObject
				cur_func.name = strBuilder.String()
				strBuilder.Reset()
			} else{
				if !isWhitespace(c){
					return []NotexFunction{}, fmt.Errorf("Invalid character '%c' expected '(' or '{' ", c)
				}
			}
		case FuncParseParams:
			switch(c){
			case '"':
        state = FuncParseQuotedParam
			case ',':
				cur_func.parameters = append(cur_func.parameters, strBuilder.String())
				//fmt.Printf("adding param: %s", strBuilder.String())
				strBuilder.Reset()
			case ')':
				cur_func.parameters = append(cur_func.parameters, strBuilder.String())
				//fmt.Printf("adding param: %s", strBuilder.String())
				strBuilder.Reset()
				state = FuncParseObject
			default:
				strBuilder.WriteRune(c)
			}

		case FuncParseQuotedParam:
			switch(c){
			case '"':
				state = FuncParseParams  // closing quote, back to normal
			case '\\':
				state = FuncParseQuotedParamEscape
			default:
				strBuilder.WriteRune(c)
			}
		case FuncParseQuotedParamEscape:
			strBuilder.WriteRune(c)
			state = FuncParseQuotedParam

		case FuncParseObject:
			switch (c){
			case '{':
				groupDepth++
				if groupDepth == 0{
					cur_index++
					continue
				}
			case '}':
				if (groupDepth == 0){
					cur_func.end = cur_index+1
					cur_func.object = strBuilder.String()
					NFs = append(NFs, *cur_func)
					cur_func = nil
					state = FuncNormalState
					strBuilder.Reset()
					cur_index++
					groupDepth = -1
					continue
				} else{
					groupDepth--
				}
			}
			strBuilder.WriteRune(c)
		}
		cur_index++
	}

	if cur_func != nil{
		return NFs, fmt.Errorf("Reached end of line before end of function")
	}

	

	return NFs, nil	
}

func applyFunction(funcs map[string]string, inp_str string) string {
	var strBuilder strings.Builder
	var ret string 
	fns, err := scanFunction(inp_str)
	if err != nil{
		log.Fatal(err)
	}
	fnLen := len(fns)
	
	if fnLen == 0{
		return inp_str
	}

	f := fns[0]
	sf, err := scanFunction(f.object)
	if err != nil{
		log.Fatal(err)
	}
	for ; len(sf) > 0; sf, err = scanFunction(f.object){
		if err != nil{
			log.Fatal(err)
		}
		f.object = applyFunction(funcs, f.object)
	}

	params := []string{f.name}
	params = append(params, f.parameters...)
	params = append(params, f.object)
	cmd := exec.Command(funcs[f.name], params...)

	out, err := cmd.Output()
	if err != nil{ 
		fmt.Fprintf(os.Stderr, "%s:%s is undefined or unrunnable\n",
		f.name, funcs[f.name])
		log.Fatal(err) 
	}
	runes := []rune(inp_str)
	strBuilder.WriteString(string(runes[:f.start]))
	strBuilder.WriteString(strings.TrimRight(string(out), "\n"))
	strBuilder.WriteString(string(runes[f.end:]))
	
	newInpStr := strBuilder.String()
	strBuilder.Reset()
	ret = applyFunction(funcs, newInpStr)



	return ret
} 
