package main

import "fmt"

// imagine that you're working on some sort of graphical application and this should be capable of
// printing different objects like circles and rectangles and squares ...
// However, it needs to be able to render them in different ways.

// So if you have shapes like CIRCLE and SQUARE and if for example, you have,
// RASTER RENDERER && a VECTOR RENDERER
// then if you plan your application badly, you'd end up with lots of different types:
//  ╒═════════════════════════════════════════════════════════════════╕
//    RasterCircle  , RasterSquare , VectorCircle, VectorSquare .....
//  └─────────────────────────────────────────────────────────────────┘

//  ╔═════════════════════════════════════════════════════════════════════════════════════╗
//  ║ How can we actually reduce the number of times that we need to introduce?           ║
//  ║ And the answer is that you would typically do this using the bridge design pattern. ║
//  ╚═════════════════════════════════════════════════════════════════════════════════════╝

type Renderer interface {
	RenderCircle(radius float32)
}

type VectorRenderer struct{}

func (v *VectorRenderer) RenderCircle(radius float32) {
	fmt.Println("Drawing a circle of radius", radius)
}

type RasterRenderer struct {
	Dpi int
}

func (r *RasterRenderer) RenderCircle(radius float32) {
	fmt.Println("Drawing pixels for circle of radius", radius)
}

type Circle struct {
	renderer Renderer
	radius   float32
}

func (c *Circle) Draw() {
	c.renderer.RenderCircle(c.radius)
}

func NewCircle(renderer Renderer, radius float32) *Circle {
	return &Circle{renderer: renderer, radius: radius}
}

func (c *Circle) Resize(factor float32) {
	c.radius *= factor
}

func main() {
	raster := RasterRenderer{}
	vector := VectorRenderer{}

	circle1 := NewCircle(&vector, 5)
	circle2 := NewCircle(&raster, 10)

	circle1.Draw() // Drawing a circle of radius 5
	circle2.Draw() // Drawing pixels for circle of radius 10
}
