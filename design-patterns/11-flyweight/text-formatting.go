package main

import (
	"fmt"
	"strings"
	"unicode"
)

// this boolean array, the underlying boolean there is going to be of the same size as the number
// of letters in plain text. So we're going to basically allocate boolean value for every single
// letter indicating whether or not this letter has to be capitalized and.

type FormattedText struct {
	plainText  string
	capitalize []bool
}

func (f *FormattedText) String() string {
	sb := strings.Builder{}
	for i := 0; i < len(f.plainText); i++ {
		c := f.plainText[i]
		if f.capitalize[i] {
			sb.WriteRune(unicode.ToUpper(rune(c)))
		} else {
			sb.WriteRune(rune(c))
		}
	}
	return sb.String()
}

func NewFormattedText(plainText string) *FormattedText {
	return &FormattedText{plainText, make([]bool, len(plainText))}
}

func (f *FormattedText) Capitalize(start, end int) {
	for i := start; i <= end; i++ {
		f.capitalize[i] = true
	}
}

// this approach is extremely inefficient. It's inefficient because we are specifying
// a huge boolean in the rate one element for every single character inside plaintext.
// ********************************************
// ********************************************
// ********************************************

type TextRange struct {
	Start, End               int
	Capitalize, Bold, Italic bool
}

func (t *TextRange) Covers(position int) bool {
	return position >= t.Start && position <= t.End
}

type BetterFormattedText struct {
	plainText  string
	formatting []*TextRange
}

func (b *BetterFormattedText) String() string {
	sb := strings.Builder{}

	for i := 0; i < len(b.plainText); i++ {
		c := b.plainText[i]
		for _, r := range b.formatting {
			if r.Covers(i) && r.Capitalize {
				c = uint8(unicode.ToUpper(rune(c)))
			}
		}
		sb.WriteRune(rune(c))
	}

	return sb.String()
}

func NewBetterFormattedText(plainText string) *BetterFormattedText {
	return &BetterFormattedText{plainText: plainText}
}

func (b *BetterFormattedText) Range(start, end int) *TextRange {
	r := &TextRange{start, end, false, false, false}
	b.formatting = append(b.formatting, r)
	return r
}

func main() {
	text := "This is a brave new world"

	ft := NewFormattedText(text)
	ft.Capitalize(10, 15) // brave
	fmt.Println(ft.String())

	bft := NewBetterFormattedText(text)
	bft.Range(16, 19).Capitalize = true // new
	fmt.Println(bft.String())

	// the critical distinction between the first and second approach is, of course,
	// the savings in memory. And this is what the flyweight pattern is actually for
	// is essentially a trick to try to avoid using too much memory and to instead
	// maybe introduce additional sort of temporary objects like we're doing with text ranges.
	//  But these temporary objects allow us to save a lot of memory
}
