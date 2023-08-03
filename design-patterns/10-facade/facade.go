package main

import "log"

// the facade design pattern is basically a way of providing a simple, easy to use interface,
// so a set of methods or functions over a large and sophisticated body of code.
// So behind the scenes you could have lots and lots of structure and functions and whatnot.
// But the facade can provide you a simple API where everything is more or less self descriptive
// and intuitive and you don't have to know about all the complicated internal details.

type Buffer struct {
	width, height int
	buffer        []rune
}

func NewBuffer(width, height int) *Buffer {
	return &Buffer{width, height, make([]rune, width*height)}
}

func (b *Buffer) At(index int) rune {
	return b.buffer[index]
}

// a viewport basically shows you just a part of the buffer as a
// particular offset, starting at a particular line.
type Viewport struct {
	buffer *Buffer
	offset int
}

func NewViewport(buffer *Buffer) *Viewport {
	return &Viewport{buffer: buffer}
}

func (v *Viewport) GetCharacterAt(index int) rune {
	return v.buffer.At(v.offset + index)
}

// So as you can see, we have a situation where we have a buffer and viewport and you
// can imagine a console, a multi buffer console being a kind of combination. So you would
// have lots of U. Ports and you would have lots of buffers. But you also want a simple API
// for just creating a console which contains all of these viewport and buffers behind the scenes.
// And that is where you would build a facade.
// a facade over buffers and viewports
type Console struct {
	buffers   []*Buffer
	viewports []*Viewport
	offset    int
}

// an initial which creates a default scenario. Now, the default scenario is
// where a console has just a single buffer and a single viewport.
// And if you look at terminals in Windows and Mac OS and Linux, the default implementation of
// a console or a terminal is the one where you have just one buffer and one viewport.
// Not very exciting, but this is the kind of simplified API that you would expect a facade to.
func NewConsole() *Console {
	b := NewBuffer(10, 10)
	v := NewViewport(b)
	return &Console{[]*Buffer{b}, []*Viewport{v}, 0}
}

// kind of function for figuring out a character's position at a particular point in the console.
func (c *Console) GetCharacterAt(index int) rune {
	return c.viewports[0].GetCharacterAt(index)
}

func main() {
	c := NewConsole()
	u := c.GetCharacterAt(1)
	log.Println(u)
}
