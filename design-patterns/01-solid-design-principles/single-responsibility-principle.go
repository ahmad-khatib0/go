package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

var entryCount = 0

type Journal struct {
	entries []string
}

func (j *Journal) String() string {
	return strings.Join(j.entries, "\n")
}

func (j *Journal) AddEntry(text string) int {
	entryCount++
	entry := fmt.Sprintf("%d: %s",
		entryCount,
		text)
	j.entries = append(j.entries, entry)
	return entryCount
}

func (j *Journal) RemoveEntry(index int) {
	// ...
}

// separation of concerns
// what we're doing below is we're breaking the single responsibility principle because
// the responsibility of the journal is to do with the management of the entries.
// The responsibility of the journal is not to deal with persistence. Persistence can be
// handled by a separate component, whether it's a separate package or ...

func (j *Journal) Save(filename string) {
	_ = ioutil.WriteFile(filename,
		[]byte(j.String()), 0644)
}

func (j *Journal) Load(filename string) {}

func (j *Journal) LoadFromWeb(url *url.URL) {}

// below is a potintioal solution for this break
var lineSeparator = "\n"

func SaveToFile(j *Journal, filename string) {
	_ = ioutil.WriteFile(filename,
		[]byte(strings.Join(j.entries, lineSeparator)), 0644)
}

type Persistence struct {
	lineSeparator string
}

func (p *Persistence) saveToFile(j *Journal, filename string) {
	_ = ioutil.WriteFile(filename,
		[]byte(strings.Join(j.entries, p.lineSeparator)), 0644)
}

func main_() {
	j := Journal{}
	j.AddEntry("I cried today.")
	j.AddEntry("I ate a bug")
	fmt.Println(strings.Join(j.entries, "\n"))

	// separate function
	SaveToFile(&j, "journal.txt")

	//
	p := Persistence{"\n"}
	p.saveToFile(&j, "journal.txt")
}
