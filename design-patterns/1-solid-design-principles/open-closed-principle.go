package main

import "fmt"

// OCP
// Open for extension , closed for modification

type Color int

const (
	red Color = iota
	green
	blue
)

type Size int

const (
	small Size = iota
	medium
	large
)

type Product struct {
	Name  string
	Color Color
	Size  Size
}

type Filter struct {
}

// *********************************

func (f *Filter) FilterByColor(products []Product, color Color) []*Product {
	res := make([]*Product, 0)
	for i, v := range products {
		if v.Color == color {
			res = append(res, &products[i])
		}
	}

	return res
}

func (f *Filter) FilterBySize(products []Product, size Size) []*Product {
	res := make([]*Product, 0)
	for i, v := range products {
		if v.Size == size {
			res = append(res, &products[i])
		}
	}
	return res
}

func (f *Filter) filterBySizeAndColor(
	products []Product, size Size,
	color Color) []*Product {
	result := make([]*Product, 0)

	for i, v := range products {
		if v.Size == size && v.Color == color {
			result = append(result, &products[i])
		}
	}

	return result
}

// so the demonstration here is a violation of the open, close principle, because what we're
// doing as we're going back into the product type and we're modifying we're adding additional methods
// on the product, we are sort of interfering with something that's already been written and already
// been tested. And the open closed principle is all about being open for extension.
// So you want to be able to extend a scenario, but by maybe adding additional types, maybe
// additional, just freestanding functions, but without modifying something that you've already written
// and you've already tested. so bescially you really want to leave the filter type
// alone, you want to leave the filter type alone. You don't want to come back to it and keep adding
// more and more methods to it and and all that sort of thing. You want to basically have some sort
// of extendible set up, and that's exactly what we can do, what we can get if we use (THE SPECIFICATION PATTERN).
// *********************************

// #################################
// (THE SPECIFICATION PATTERN).
// So the specification pattern is somewhat different because it has a bunch of interfaces.
// It has a bunch of additional elements done for flexibility.

// the idea behind the specification interface is you are testing whether or
//
// not a product specified here via pointer satisfies some criteria.
type Specification interface {
	IsSatisfied(p *Product) bool
}

type ColorSpecification struct {
	color Color
}

func (c ColorSpecification) IsSatisfied(p *Product) bool {
	return p.Color == c.color
}

type SizeSpecification struct {
	size Size
}

func (s SizeSpecification) IsSatisfied(p *Product) bool {
	return p.Size == s.size
}

type BetterFilter struct{}

func (b *BetterFilter) Filter(products []Product, spec Specification) []*Product {
	result := make([]*Product, 0)
	for i, v := range products {
		if spec.IsSatisfied(&v) {
			result = append(result, &products[i])
		}
	}
	return result
}

//  #################################
// the approach with the specification pattern gives you more flexibility, because
// if you want to filter by a particular new type, all you have to do is you have to make a new specification.
// So, for example, here we have a color specification, but you decide that you want to filter by size.
// So all you have to do is make a size specification and make sure that it conforms to the specification
// interface. That's pretty much all the would have to do, and that follows the open close principle.
// So the times in this case, the interface type is open for extension, meaning you can implement this
// interface, but it's close to a modification, which means that you are unlikely to ever modify the
// specification interface and in a similar fashion, you are unlikely to ever modify BetterFilter because
// there's no reason for us to do so. It's very flexible. It takes a bunch of products and a specification
// and that's pretty much all there is to it.

// #################################
// COMPOSITE SPECIFICATION is just a combinator. It just combines two different specifications.

type AndSpecification struct {
	first, second Specification
}

func (a AndSpecification) IsSatisfied(p *Product) bool {
	return a.first.IsSatisfied(p) && a.second.IsSatisfied(p)
}

func main() {
	apple := Product{"Apple", green, small}
	tree := Product{"Tree", green, large}
	house := Product{"House", blue, large}

	products := []Product{apple, tree, house}
	fmt.Printf("Green products (old) \n ")

	f := Filter{}
	for _, v := range f.FilterByColor(products, green) {
		fmt.Printf(" - %s is green \n", v.Name)
	}

	fmt.Printf("Green products (new ) \n")
	greenSpec := ColorSpecification{green}
	b := BetterFilter{}
	for _, v := range b.Filter(products, greenSpec) {
		fmt.Printf(" - %s is green \n", v.Name)
	}

	largeSpec := SizeSpecification{large}
	largeGreenSpec := AndSpecification{largeSpec, greenSpec}
	fmt.Print("Large blue items:\n")
	for _, v := range b.Filter(products, largeGreenSpec) {
		fmt.Printf(" - %s is large and green\n", v.Name)
	}
}
