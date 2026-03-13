package main

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func read_input() string {
	var input_str string 
	switch(len(os.Args)){
	case 1:
		b, err := io.ReadAll(os.Stdin)
		if err != nil{
			log.Fatal(err)
		}
		input_str = string(b)

	case 2,3:
		b, err := os.ReadFile(os.Args[1])
		if err != nil{
			log.Fatal(err)
		}
		input_str = string(b)
	}
	return input_str
}


func getSupplies(path string) []string{
	var stdout bytes.Buffer
	cmd := exec.Command(path, "--supplies")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
				log.Fatal(err)
	}
	return strings.Fields(stdout.String())
}

type OutputDocument struct{
	head Stringable
	body string
}

func (o OutputDocument) String() string{
	var strBuilder strings.Builder
	hs, err := o.head.String()
	if err != nil{
		log.Fatal(err)
	}
	strBuilder.WriteString("<!doctype html>")
	strBuilder.WriteString("<html lang=\"en-US\">")
	strBuilder.WriteString(hs)
	strBuilder.WriteString(o.body)
	strBuilder.WriteString("</html>")
	return strBuilder.String()
}


func main() {
	input := read_input()
	func_defs, attr_defs := init_functions()
	no_coms := removeComments(input)
	escaped := html.EscapeString(no_coms)
	func_applied := applyFunction(func_defs, escaped)
	body_str, head_str := separateHeader(func_applied)
	head := handleHead(head_str)
	root, err := handleGroups(body_str)
	if err != nil{
		log.Fatal(err)
	}
	style_map := make(map[string]string)
	compileAttributes(attr_defs, &root, style_map )
	head.styles= style_map
	body := parse(root)
	doc := OutputDocument{head: head, body: body}
	fmt.Println(doc.String())

}


func init_functions()	(funcs_p map[string]string, attrs_p map[string]string){
	funcs := make(map[string]string)
	attrs := make(map[string]string)
	path := os.Getenv("NOTEX_TEST_PATH")

	for dir := range strings.SplitSeq(path,":"){
		files, err := os.ReadDir(dir)
		if err != nil{
			log.Fatal(err)
		}
		for _, file := range files {
			if !file.IsDir(){
				full_path := dir+"/"+ file.Name()
				supplies:= getSupplies(full_path)
				for _, f := range supplies{
					if strings.HasPrefix(f, "!"){
						attrs[strings.Trim(f,"!")] = full_path
					} else{
					funcs[f] = full_path
					}
				}
			}
		}
	}
	return funcs, attrs
}


