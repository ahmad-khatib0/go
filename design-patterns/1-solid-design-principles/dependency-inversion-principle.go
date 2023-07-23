package main

import "fmt"

//  ┌──────────────────────────────────────────────────────────────────────┐
//
//     Dependency Inversion Principle: HLM (high level modules) should not
//     depend on LLM (low level modules ) Both should depend on abstractions
//
//
//  └──────────────────────────────────────────────────────────────────────┘

//  +-----------------------------------------------------------------------------------+
//  | By abstractions, we typically mean interfaces, at least in go in other languages. |
//  | You would talk about abstract classes and  base classes                           |
//  +-----------------------------------------------------------------------------------+

type Relationship int

const (
	Parent Relationship = iota
	Child
	Sibling
)

type Person struct {
	name string
	// other useful stuff here
}

// how do we model relationships between different people?
type Info struct {
	from         *Person
	relationship Relationship
	to           *Person
}

// low level module
// The reason why it's low level is because it kind of storage. so, this could be in
// a database or something. It could be on the web or somewhere. So it's basically the storage mechanism
type Relationships struct {
	relations []Info
}

func (r *Relationships) AddParentAndChild(parent, child *Person) {
	r.relations = append(r.relations, Info{parent, Parent, child})
	r.relations = append(r.relations, Info{child, Child, parent})
}

// high level module
// this is the high level module designed to operate on data and perform some sort of research.
type Research struct {
	// break DIP  (HLM should not depend on LLM )
	// relationships Relationship

	// the solution
	browser RelationshipBrowser
}

func (r *Research) Investigate() {

	//  +----------------------------------------------------------------------------------------------+
	//  | imagine if relationships, the low level module decides to change the storage mechanic from a |
	//  | slice to, let's say, a database. So what happens then? And the answer is that the code,      |
	//  | which depends on the low level module, actually breaks.                                      |
	//  | So all of this is going to break because, for example, you can no longer use a for loop here.|
	//  | So that's obviously something we want to avoid. And that's what the dependency inversion     |
	//  | principle is trying to protect us from these situations where everything breaks down.        |
	//  +----------------------------------------------------------------------------------------------+

	//relations := r.relationships.relations
	//for _, rel := range relations {
	//	if rel.from.name == "John" &&
	//		rel.relationship == Parent {
	//		fmt.Println("John has a child called", rel.to.name)
	//	}
	//}

	for _, p := range r.browser.FindAllChildrenOf("John") {
		fmt.Println("John has a child called", p.name)
	}
}

type RelationshipBrowser interface {
	FindAllChildrenOf(name string) []*Person
}

func (rs *Relationships) FindAllChildrenOf(name string) []*Person {
	result := make([]*Person, 0)

	for i, v := range rs.relations {
		if v.relationship == Parent && v.from.name == name {
			result = append(result, rs.relations[i].to)
		}
	}

	return result
}

func main() {

	parent := Person{"John"}
	child1 := Person{"Chris"}
	child2 := Person{"Matt"}

	// low-level module
	relationships := Relationships{}
	relationships.AddParentAndChild(&parent, &child1)
	relationships.AddParentAndChild(&parent, &child2)

	research := Research{&relationships}
	research.Investigate()
}
