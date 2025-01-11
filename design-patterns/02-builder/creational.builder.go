package main

import (
	"fmt"
	"strings"
)

func main() {

	//*********************************
	hello := "hello"
	sb := strings.Builder{}

	sb.WriteString("<p>")
	sb.WriteString(hello)
	sb.WriteString("<p>")
	fmt.Println(sb.String())

	words := []string{"hello", "world"}
	sb.Reset()

	sb.WriteString("<ul>")
	for _, v := range words {
		sb.WriteString("<li>")
		sb.WriteString(v)
		sb.WriteString("</li>")
	}
	sb.WriteString("</ul>")
	fmt.Println(sb.String())
	// this is cofusing and not have clear steps to serve a porpuse of building for example an html element
	//*********************************

	// #################################
	// working with build now
	b := NewHtmlBuilder("ul")

	b.AddChildFluent("li", "hello").
		AddChildFluent("li", "world")
	fmt.Println(b.String())
}

const (
	indentSize = 2
)

type HtmlElement struct {
	name, text string
	elements   []HtmlElement
}

func NewHtmlBuilder(rootName string) *HtmlBuilder {
	return &HtmlBuilder{rootName, HtmlElement{rootName, "", []HtmlElement{}}}
}

func (e *HtmlElement) String() string {
	return e.string(0)
}

func (e *HtmlElement) string(indent int) string {
	sb := strings.Builder{}
	i := strings.Repeat(" ", indentSize*indent)

	sb.WriteString(fmt.Sprintf("%s<%s>\n", i, e.name))

	if len(e.text) > 0 {
		sb.WriteString(strings.Repeat(" ", indentSize*(indent+1)))
		sb.WriteString(e.text)
		sb.WriteString("\n")
	}

	for _, el := range e.elements {
		sb.WriteString(el.string(indent + 1))
	}

	sb.WriteString(fmt.Sprintf("%s</%s>\n", i, e.name))
	return sb.String()
}

type HtmlBuilder struct {
	rootName string
	root     HtmlElement
}

func (b *HtmlBuilder) String() string {
	return b.root.String()
}

func (b *HtmlBuilder) AddChild(childName, childText string) {
	e := HtmlElement{childName, childText, []HtmlElement{}}
	b.root.elements = append(b.root.elements, e)
}

// fluent interfaces shows up quite a bit inside the builder pattern
// So a fluent interface basically is an interface that allows you to chain calls together.
// Now, changing calls in go isn't really that convenient because  you leave the dot hanging
// at the end instead of at the beginning.
func (b *HtmlBuilder) AddChildFluent(chName, chText string) *HtmlBuilder {
	e := HtmlElement{chName, chText, []HtmlElement{}}
	b.root.elements = append(b.root.elements, e)

	// so fluent interfaces, basically returning the receiver appointed to the
	// receiver at the end of the method
	return b // we return the builder again in the fluent interface style
}
