package main

import (
	"bytes"
	"fmt"
//	"html"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func read_input() (string, *string) {
	var input_str string 
	var output *string = nil
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
		if len(os.Args) == 3{
			output = &os.Args[2]
		}
	}
	return input_str, output
}


func getSupplies(path string) []string{
	var stdout bytes.Buffer
	cmd := exec.Command(path, "--supplies")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
				log.Fatal("Failure reading path:",path, err)
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
	strBuilder.WriteString("<!doctype html>\n")
	strBuilder.WriteString("<html lang=\"en-US\">\n")
	strBuilder.WriteString(hs)
	strBuilder.WriteString(o.body)
	strBuilder.WriteString("</html>")
	return strBuilder.String()
}


func main() {
	input, output := read_input()
	func_defs, attr_defs := init_functions()
	no_coms := removeComments(input)
	func_applied := applyFunction(func_defs, no_coms)
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
	if output == nil{
		fmt.Println(doc.String())
	} else{
		f, err := os.Create(*output)
		if err != nil{
			log.Fatal(err)
		}
		defer f.Close()
		_, err = fmt.Fprint(f, doc.String())
		if err != nil{
			log.Fatal(err)
		}

	}

}


func init_functions()	(funcs_p map[string]string, attrs_p map[string]string){
	funcs := make(map[string]string)
	attrs := make(map[string]string)
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
	//fmt.Fprintf(os.Stderr, "Found functions: %v\n", funcs)
	//fmt.Fprintf(os.Stderr, "Found attributes: %v\n", attrs)
	return funcs, attrs
}


