package main

import "fmt"

//  ╒═══════════════════════════════════════════════════════╕
//
//    one way of extending a builder that you already have
//    is by using a functional programming approach.
//
//  └───────────────────────────────────────────────────────┘

type Person struct {
	name, position string
}

type personMod func(*Person)
type PersonBuilder struct {
	actions []personMod
}

func (b *PersonBuilder) Called(name string) *PersonBuilder {
	b.actions = append(b.actions, func(p *Person) {
		p.name = name
	})
	return b
}

func (b *PersonBuilder) WorksAsA(position string) *PersonBuilder {
	b.actions = append(b.actions, func(p *Person) {
		p.position = position
	})

	return b
}

func (b *PersonBuilder) Build() *Person {
	p := Person{}
	for _, a := range b.actions {
		a(&p)
	}
	return &p
}

func main() {
	b := PersonBuilder{}
	p := b.Called("ahmad").WorksAsA("programmer").Build()

	fmt.Println(p)
}

//  ╒═════════════════════════════════════════════════════════════════════════════════╕
//    The benefit of the setup is that it's very easy to extend the builder with
//    additional build actions without messing about, with making new builders which
//    aggregate the current builder and so on and so forth.
//  └─────────────────────────────────────────────────────────────────────────────────┘

//  ╒════════════════════════════════════════════════════════════════════════════════════════════════╕
//    So this set up, all it does is it illustrates that effectively what you can do is you
//    can have a kind of delayed application of all of those modifications. So your builder,
//    instead of just doing the modifications in place, it can keep a list of actions, a list
//    of changes to perform upon the object that's being constructed. And then when you call
//    build, what you do is you create just a default implementation of the object and then you
//    go through every single action and you apply that action to that object that you are returning
//    and then subsequently just return that object.
//  └────────────────────────────────────────────────────────────────────────────────────────────────┘
