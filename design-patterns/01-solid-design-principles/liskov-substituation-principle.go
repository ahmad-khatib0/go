package main

import "fmt"

type Sized interface {
	GetWidth() int
	SetWidth(width int)
	GetHeight() int
	SetHeight(height int)
}

type Rectangle struct {
	width, height int
}

func (r *Rectangle) GetWidth() int {
	return r.width
}

func (r *Rectangle) SetWidth(width int) {
	r.width = width
}

func (r *Rectangle) GetHeight() int {
	return r.height
}

func (r *Rectangle) SetHeight(height int) {
	r.height = height
}

// modified LSP
// If a function takes an interface and works with a type T that implements this
// interface, any structure that aggregates T should also be usable in that function.
type Square struct {
	Rectangle
}

// So the risk of substitution principle basically states that if you continue to use generalisations like
// interfaces, for example, then you should not have inherited or you should not have implementors of
// those generalisations break some of the assumptions which are set up at the higher level.
// So at the higher level, we kind of assume that if you have a set object and you said it's height,
// you are just setting its height, not both the height and the width.

func (s *Square) SetWidth(width int) {
	s.width = width
	s.height = width
}

func (s *Square) SetHeight(height int) {
	s.width = height
	s.height = height
}

func NewSquare(size int) *Square {
	sq := Square{}
	sq.width = size
	sq.height = size
	return &sq
}

// So you should be able to continue taking sized objects instead of somehow
// figuring out in here, for example, by doing type checks whether you have a rectangle or
// a square, it should still work in the generalised case.
func UseIt(sized Sized) {
	width := sized.GetWidth()
	sized.SetHeight(10)

	expectedArea := width * 10
	actualArea := sized.GetWidth() * sized.GetHeight()
	fmt.Printf("Expected an area of: %d , but got %d  \n", expectedArea, actualArea)
}

// ################################# a potintioal solutin $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$

// to recap, basically, the idea of the risk of substitution principle is that
// the behaviour of implementors of a particular type, like in this case, the size the
// interface should not break the core fundamental behaviors that you rely on.

type Square2 struct {
	size int
}

func (s *Square2) Rectangle() Rectangle {
	return Rectangle{s.size, s.size}
}

func main() {
	rc := &Rectangle{2, 3}
	UseIt(rc)
}
