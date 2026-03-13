package main

import (
	"fmt"
	"log"
	"strings"
)

type NotexHeading struct{
	strength int
	text string
}

func (h NotexHeading) String() (string, error){
	switch(h.strength){
	case 0,1,2,3:
		return fmt.Sprintf("<h%d>%s</h%d>",h.strength+1, h.text, h.strength+1), nil
	case 4:
		return fmt.Sprintf("<u><b>%s</b></u>", h.text ), nil
	default:
		return fmt.Sprintf("<b>%s</b>", h.text ), nil
	}
	

}

func handleTitle (inp_str string, group_depth int) NotexHeading{
	var ret NotexHeading
	var strBuilder strings.Builder
	inp := []rune(inp_str)
	for i, c := range inp{
		if i == 0 {
			continue
		}
		if c == '{'{
			log.Fatal(fmt.Errorf("Headings cannot take groups as an argument"))
		}
		strBuilder.WriteRune(c)
	}
	ret.strength = group_depth
	ret.text = strBuilder.String()
	return ret
}
