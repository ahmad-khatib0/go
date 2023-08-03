package decorator

import "fmt"

type Shape interface {
	Render() string
}

type Circle struct {
	Radius float32
}

func (c *Circle) Render() string {
	return fmt.Sprintf("Circle of radius %f", c.Radius)
}

func (c *Circle) Resize(factor float32) {
	c.Radius *= factor
}

type Square struct {
	Side float32
}

func (s *Square) Render() string {
	return fmt.Sprintf("Square with side %f", s.Side)
}

// possible, but not generic enough
type ColoredSquare struct {
	Square
	Color string
}

// this is a decorator
type ColoredShape struct {
	Shape Shape
	Color string
}

func (c *ColoredShape) Render() string {
	return fmt.Sprintf("%s has the color %s", c.Shape.Render(), c.Color)
}

// one upside is the decorators can be composed, which means you can
// apply decorators to decorators. There's no problem doing this.
type TransparentShape struct {
	Shape        Shape
	Transparency float32
}

func (t *TransparentShape) Render() string {
	return fmt.Sprintf("%s has %f%% transparency", t.Shape.Render(), t.Transparency*100.0)
}

func main() {
	circle := Circle{2}
	fmt.Println(circle.Render())

	// The problem is that once you've made a decorator, once you've put colored shape over the
	// ordinary shape, what you cannot do is you cannot say a redCircle.Resize() because
	// the resize method is no longer available. Unfortunately, and there is no real solution
	// to this because you are not aggregating anything. You've lost that particular method.
	// The only way you can restore it, so to speak, is if you add it again.
	// The problem is that how do you added without also adding it to the interface?
	//  Because remember, it's only the circle type that has the resize method. The square type,
	// for example, does not have a precise method. So you cannot add this to the interface.
	//  Unfortunately, that's a real life limitation of the decorator approach.
	redCircle := ColoredShape{&circle, "Red"}
	fmt.Println(redCircle.Render())

	// we can use a transparency decorator over the colour
	// shape decorator so we can apply a decorator over another decorator.
	rhsCircle := TransparentShape{&redCircle, 0.5}
	fmt.Println(rhsCircle.Render())
}
