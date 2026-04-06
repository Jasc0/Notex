package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
)


type NotexAttribute interface{
	String() string
}


type NotexUncompiledAttribute struct {
	name string
	parameters []string
}

func (nua NotexUncompiledAttribute) String() string{
	var strBuilder strings.Builder
	strBuilder.WriteString(string(nua.name))
	for _, p := range(nua.parameters){
		strBuilder.WriteString(p)
	}
	return strBuilder.String()
}

type NotexRawAttribute string

func (nra NotexRawAttribute) String() string{
	return string(nra)
}

func (na NotexUncompiledAttribute) Compile (defs, style_map  map[string]string) (NotexRawAttribute, error){
	name := fmt.Sprintf("!%s",na.name)
	params := []string{name}
	params = append(params, na.parameters...)
	cmd := exec.Command(defs[na.name], params...)

	outb, err := cmd.Output()
	if err != nil{ 
		return NotexRawAttribute(""), err
	}
	out := string(outb)
	// fill in style_map
	params = append(params, "--style")
	st_cmd := exec.Command(defs[na.name], params...)
	st_outb, err := st_cmd.Output()
	if err != nil{ 
		return NotexRawAttribute(""), err
	}
	st_out := string(st_outb)
	style_map[na.String()] = st_out

	return NotexRawAttribute(strings.TrimRight(out, "\n")), nil
}

func validAttrName(r rune) bool{
	return validFuncName(r)
}

func compileAttrSlice(attr_defs map[string]string,attrs []NotexAttribute, style_map map[string]string) []NotexAttribute{
	for i, a := range attrs{
		switch ta := a.(type){
		case NotexUncompiledAttribute:
			rawA, err := ta.Compile(attr_defs, style_map)
			if err != nil{
				log.Fatal("error compiling:",ta.name, ": ",err)
			}
			attrs[i] = rawA
		}
	}
	return attrs
}

func compileAttributes(attr_defs map[string]string, group *NotexGroup, style_map map[string]string){
	for _, s := range group.sequence{
		switch ts := s.(type){
		case NotexGroup:
			compileAttributes(attr_defs, &ts, style_map)
		}
	}
	group.Attributes = compileAttrSlice(attr_defs, group.Attributes, style_map)
}

type AttributeParseState int 
const (
	AttrBegin AttributeParseState = iota
	AttrPotential 
	AttrReadName 
	AttrReadParams
)

func scanAttributes(inp_str string) ([]NotexUncompiledAttribute, error){
	
	inp := []rune(inp_str);
	state := AttrBegin
	attrs := []NotexUncompiledAttribute{}
	var cur_attr *NotexUncompiledAttribute
	var strBuilder strings.Builder
	for _, c := range inp{
		switch(state){
		case AttrBegin:
			if c == '!'{
				state = AttrPotential
			}
		case AttrPotential:
			if c == '@'{
				cur_attr = new(NotexUncompiledAttribute)
				state = AttrReadName
			}else{
				state =AttrBegin
			}

		case AttrReadName:
			if validAttrName(c){
				strBuilder.WriteRune(c)
			} else{
				if c == '('{
					state = AttrReadParams
					cur_attr.name = strBuilder.String()
					strBuilder.Reset()
				} else{
					return []NotexUncompiledAttribute{}, fmt.Errorf("Invalid character in attribute name")
				}
			}
		case AttrReadParams:
			switch(c){
			case ',':
				cur_attr.parameters = append(cur_attr.parameters, strBuilder.String())
				strBuilder.Reset()
			case ')':
				cur_attr.parameters = append(cur_attr.parameters, strBuilder.String())
				attrs = append(attrs, *cur_attr)
				cur_attr = nil
				state = AttrBegin
			default:
				strBuilder.WriteRune(c)
			}
		}
	}

	return attrs, nil

}

func containsAttribute(inp_str string) bool{
	return strings.HasPrefix(strings.Trim(inp_str,"\n \t"), "!@")
}



