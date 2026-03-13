package main

import "fmt"


type NotexPara struct{
	text string
}

func (p NotexPara) String() (string, error){
	return fmt.Sprintf("<p>%s</p>", p.text), nil
}

func handlePara (inp_str string) NotexPara{
	var ret NotexPara
	ret.text = inp_str
	return ret
}
