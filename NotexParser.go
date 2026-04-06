package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type NotexHTMLTag struct{
	tag string;
	isOpen bool;
	attributes []NotexAttribute;
} 
func (tag NotexHTMLTag) String() (string, error){
	if tag.isOpen{
		if len(tag.attributes) > 0{
			var strBuilder strings.Builder
			strBuilder.WriteString(fmt.Sprintf("<%s class=\"", tag.tag))
			for i, a := range tag.attributes{
				switch ta := a.(type){
				case NotexRawAttribute:
					strBuilder.WriteString(ta.String())
					if i +1 < len(tag.attributes){
					strBuilder.WriteString(" ")

					}
				default:
					return "", fmt.Errorf("Uncompiled attribute found: %v", ta)
				}
			}
			strBuilder.WriteString("\">")
			return strBuilder.String(), nil

		} else{
			return fmt.Sprintf("<%s>", tag.tag), nil
		}
	}
	return fmt.Sprintf("</%s>", tag.tag), nil

}


func readFile(path string) string{
	b, err := os.ReadFile(path)
	if err != nil{
		log.Fatal(err)
	}
	str := string(b)
	return str
}

func parse( root NotexGroup) string {
	var strBuilder strings.Builder
	queue := []Stringable{}

	opener := NotexHTMLTag{tag: "body",
	isOpen: true,
	attributes: root.Attributes}
	closer :=  NotexHTMLTag{tag: "body", isOpen: false}

	queue = append(queue, opener)
	_parse(root, &queue)
	queue = append(queue, closer)

	for _, s := range queue{
		st, err := s.String()
		if err != nil{
			log.Fatal(err)
		}
		strBuilder.WriteString(st)
		strBuilder.WriteString("\n")
	}
	return strBuilder.String()
}

type GroupState struct{
	queue *[]Stringable
	inOrderedList bool;
	inUnorderedList bool;
}
func (gs *GroupState) updateLists(isInOl, isInUl bool) {
	var tag NotexHTMLTag
	enque := true
	if gs.inOrderedList && !isInOl{
		gs.inOrderedList = false
		tag = NotexHTMLTag{tag: "ol", isOpen: false}
	} else if !gs.inOrderedList && isInOl{
		gs.inOrderedList = true
		tag = NotexHTMLTag{tag: "ol", isOpen: true}

	} else if gs.inUnorderedList && !isInUl{
		gs.inUnorderedList = false
		tag = NotexHTMLTag{tag: "ul", isOpen: false}

	} else if !gs.inUnorderedList && isInUl{
		gs.inUnorderedList = true
		tag = NotexHTMLTag{tag: "ul", isOpen: true}
	} else{
		enque = false
	}

	if enque{
		*gs.queue = append(*gs.queue, tag)
	}
}

func _parse(group NotexGroup, queue *[]Stringable){
	var groupAsArgument Stringable = nil
	state := GroupState{queue: queue, inOrderedList: false, inUnorderedList: false}
	for _, s := range group.sequence{
		switch ts := s.(type){
		case NotexGroup:
			var opener, closer NotexHTMLTag
			 if groupAsArgument != nil{
				 switch tg := groupAsArgument.(type){
				 case NotexListItem:
					 opener = NotexHTMLTag{tag: "li",
					 isOpen: true,
					 attributes: ts.Attributes}
					 _ = tg
					 closer = NotexHTMLTag{tag: "li",
					 isOpen: false}
				 }
			 } else{
					 opener = NotexHTMLTag{tag: "div",
					 isOpen: true,
					 attributes: ts.Attributes}
					 closer =  NotexHTMLTag{tag: "div", isOpen: false}
					 

			 } 
			 *queue = append(*queue, opener)
			 _parse(ts, queue)
			 *queue = append(*queue, closer)


		case NotexString:
			str, err := ts.String()
			if len(str) < 1{
				continue
			}
			if err != nil{
				log.Fatal(err)
			}
			switch(ts.startsWith()){
			case '#':
				state.updateLists(false, false)
				*queue = append(*queue, handleTitle(str, group.depth))
			case '-', '.':
				li := handleItem(str)
				if li.Type == orderedList{
					state.updateLists(true, false)
				} else{
					state.updateLists(false, true)
				}
				if li.Behavior == groupListItem{
					groupAsArgument = NotexListItem{}
					continue
				}
				*queue = append(*queue, li)
			case '<':
				*queue = append(*queue, ts)
				
			case '\\':
				state.updateLists(false, false)
				if len(str) < 2{
					continue
				}
				p := handlePara(str[1:])
				*queue = append(*queue, p)
			default:
				state.updateLists(false, false)
				p := handlePara(str)
				*queue = append(*queue, p)
			}
		}
	}
	state.updateLists(false, false)
}

func removeComments(inp_str string) string {
    var strBuilder strings.Builder
    write := true
    ignore := false
    inQuote := false
    inp := []rune(inp_str)
    for _, c := range inp {
        switch c {
        case '"':
            inQuote = !inQuote
        case '%':
            if !ignore && !inQuote {
                write = false
            }
        case '\\':
            ignore = true
        case '\n':
            write = true
            inQuote = false
        }
        if write {
            strBuilder.WriteRune(c)
        }
        if c != '\\' {
            ignore = false
        }
    }
    return strBuilder.String()
}



