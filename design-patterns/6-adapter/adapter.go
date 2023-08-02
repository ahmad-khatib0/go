package main

import (
	"fmt"
	"strings"
)

type Line struct {
	X1, Y1, X2, Y2 int // 1 start , 2 end
}

type VectorImage struct {
	Lines []Line
}

func NewRectangle(width, height int) *VectorImage {
	width -= 1
	height -= 1
	return &VectorImage{[]Line{
		{0, 0, width, 0},
		{0, 0, 0, height},
		{width, 0, width, height},
		{0, height, width, height},
	}}
}

// ↑↑↑ the interface you're given

// ↓↓↓ the interface you have

type Point struct {
	X, Y int
}

// a roster image is an image that's defined by the points, the pixels on the screen somewhere.
type RasterImage interface {
	GetPoints() []Point
}

// so this introduces an obvious problem in the entire system because we are given the interface right
// here. So the only way to create a rectangle is by making a vector image.
// But unfortunately the only way to print something is by providing a roster image.
// So what do you need in this case? And the answer is, well, obviously you need an
// adapter, you need something that takes a vector image and somehow adapted into
// something which has a bunch of points in it so that those points can be fed into
// the roster image and subsequently fed into the drop points function.
//  ╒══════════════════════════════════════════════════════════════════════════╕
//    problem: I want to print a RasterImage but I can only make a VectorImage
//  └──────────────────────────────────────────────────────────────────────────┘

func DrawPoints(owner RasterImage) string {
	maxX, maxY := 0, 0
	points := owner.GetPoints()
	for _, pixel := range points {
		if pixel.X > maxX {
			maxX = pixel.X
		}
		if pixel.Y > maxY {
			maxY = pixel.Y
		}
	}
	maxX += 1
	maxY += 1

	// preallocate
	data := make([][]rune, maxY)
	for i := 0; i < maxY; i++ {
		data[i] = make([]rune, maxX)
		for j := range data[i] {
			data[i][j] = ' '
		}
	}

	for _, point := range points {
		data[point.Y][point.X] = '*'
	}

	b := strings.Builder{}
	for _, line := range data {
		b.WriteString(string(line))
		b.WriteRune('\n')
	}

	return b.String()
}

// the solution

type vectorToRasterAdapter struct {
	points []Point
}

func (a vectorToRasterAdapter) GetPoints() []Point {
	return a.points
}

func VectorToRaster(vi *VectorImage) RasterImage {
	adapter := vectorToRasterAdapter{}
	for _, line := range vi.Lines {
		adapter.addLine(line)
	}

	return adapter // as RasterImage
}

func (a *vectorToRasterAdapter) addLine(line Line) {
	left, right := minmax(line.X1, line.X2)
	top, bottom := minmax(line.Y1, line.Y2)
	dx := right - left
	dy := line.Y2 - line.Y1

	if dx == 0 {
		for y := top; y <= bottom; y++ {
			a.points = append(a.points, Point{left, y})
		}
	} else if dy == 0 {
		for x := left; x <= right; x++ {
			a.points = append(a.points, Point{x, top})
		}
	}

	fmt.Println("generated", len(a.points), "points")
}

func minmax(a, b int) (int, int) {
	if a < b {
		return a, b
	} else {
		return b, a
	}
}

func main() {
	rc := NewRectangle(6, 4)
	a := VectorToRaster(rc) // adapter!
	_ = VectorToRaster(rc)  // adapter!
	fmt.Print(DrawPoints(a))
}
