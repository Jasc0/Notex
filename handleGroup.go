package main

import (
	"fmt"
	"log"
	"strings"
)

type Stringable interface{
	String() (s string, err error)
}

type NotexString string

func (s NotexString) String() (string, error){
	return string(s), nil
}
func (s NotexString) startsWith() rune{
	str, _ := s.String()
	runes := []rune(str)
	if len(runes) < 1{
		return 0
	}
	return runes[0]
}

type NotexGroup struct{
	parent *NotexGroup 
	sequence []Stringable
	depth int
	Attributes []NotexAttribute
}

func (g NotexGroup ) String() (string, error){
	var strBuilder strings.Builder

	for _, s := range g.sequence{
		str, err := s.String()
		if err != nil{
			return "", err
		}
		strBuilder.WriteString(str)
	}

	return strBuilder.String(), nil
}

func handleGroups(inp_str string) (NotexGroup, error) {
	var root NotexGroup 
	root.parent = nil
	root.depth = 0
	root.Attributes = append(root.Attributes, NotexRawAttribute("group_depth0"))
	var strBuilder strings.Builder
	var cur_group *NotexGroup 
	cur_group = &root
	literal := false
	flushString := func() {
		str := strings.Trim(strBuilder.String(), "\n \t")
		strBuilder.Reset()
		if len(str) > 0 {
			cur_group.sequence = append(cur_group.sequence, NotexString(str))
		}
	}

	inp := []rune(inp_str)
	for _, c := range inp{
		if literal{
			switch c {
			case '<':
				strBuilder.WriteString("&lt;")
			case '>':
				strBuilder.WriteString("&gt;")
			default:
				strBuilder.WriteRune(c)
			}
			literal = false
			continue
		}
		switch(c){
		case '\\':
			literal = true
		case '{':
			tmp := new(NotexGroup)
			tmp.parent = cur_group
			tmp.depth = cur_group.depth + 1
			gda := NotexRawAttribute(fmt.Sprintf("group_depth%d", tmp.depth))
			tmp.Attributes = append(tmp.Attributes, gda)
			flushString()
			strBuilder.Reset()
			cur_group = tmp
		case '}':
			if cur_group.depth <= 0{
				return root, fmt.Errorf("attempted to close root group before EOF")
			}
			flushString()
			cur_group.parent.sequence = append(cur_group.parent.sequence, *cur_group)
			cur_group = cur_group.parent
		case '\n':
			str := strings.Trim(strBuilder.String(), "\n \t")
			if len(str) < 1{
				strBuilder.Reset()
				continue
			}
			if containsAttribute(str){
				attrs, err := scanAttributes(str)
				if err != nil{
					log.Fatal(err)
				}
				for _, a := range(attrs){
					cur_group.Attributes = append(cur_group.Attributes, a)
				}
				
			}else{
			flushString()
			}
			strBuilder.Reset()
		default:
			strBuilder.WriteRune(c)
		}

	}
		flushString()
		if cur_group.depth > 0{
			return root, fmt.Errorf("reached EOF without closing group")
		}

	return root, nil	

}
