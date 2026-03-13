package main

import (
	"fmt"
	"strings"
)

type listType int 
const(
	unorderedList listType = iota
	orderedList
)

type listItemBehavior int 
const (
	inlineListItem listItemBehavior = iota
	groupListItem  
)

type NotexListItem struct {
	Type listType
	Behavior listItemBehavior
	inline string

}

func (li NotexListItem) String() (string, error){
	if li.Behavior == inlineListItem{
		return fmt.Sprintf("<li>%s</li>", li.inline), nil
	}
	// allow NotexGroup iteration to add to the compilation queue
	return "DO NOT USE", nil
}



func handleItem(inp_str string) NotexListItem{
	var ret NotexListItem
	var strBuilder strings.Builder
	if strings.HasPrefix(inp_str, "-"){
		ret.Type = unorderedList
	} else{
		ret.Type = orderedList
	}
	if inp_str == "-" || inp_str == "."{
		ret.Behavior = groupListItem
		return ret
	}
	inp := []rune(inp_str)
	for i, c := range inp{
		if i == 0 {
			continue
		}
		strBuilder.WriteRune(c)
	}
	ret.Behavior = inlineListItem
	ret.inline = strBuilder.String()
	return ret
}
