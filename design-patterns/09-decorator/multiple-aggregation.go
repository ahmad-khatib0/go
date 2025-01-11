package main

import "fmt"

/*
type Bird struct {
  Age int
}

func (b *Bird) Fly() {
  if b.Age >= 10 {
    fmt.Println("Flying!")
  }
}

type Lizard struct {
  Age int
}

func (l *Lizard) Crawl() {
  if l.Age < 10 {
    fmt.Println("Crawling!")
  }
}

type Dragon struct {
  Bird
  Lizard
}
*/

//  ╔═════════════════════════════════════════════════════════════════════════════════╗
//  ║ Unfortunately, this can cause certain problems. You can see that                ║
//  ║ in this case, both the bird and the lizard have an edge field. And that's       ║
//  ║ going to be a problem as we try to sort of combine the operations of these two. ║
//  ║ for example if you did       1-   d := Dragon{}    2-  d.Age = 10               ║
//  ║ this will throw an error about ambiguous selector compiling error               ║
//  ║ Now, the problem with this scenario is that you can introduce a really nasty    ║
//  ║ inconsistency into the behaviour of the dragon if you said the different ages   ║
//  ║ to different values. And after all, you don't need to separate fields.          ║
//  ║ It's a single age. You want to keep it in a single field.                       ║
//  ╚═════════════════════════════════════════════════════════════════════════════════╝

type Aged interface {
	Age() int
	SetAge(age int)
}

type Bird struct {
	age int
}

func (b *Bird) Age() int {
	return b.age
}

func (b *Bird) SetAge(age int) {
	b.age = age
}

func (b *Bird) Fly() {
	if b.age >= 10 {
		fmt.Println("Flying!")
	}
}

type Lizard struct {
	age int
}

func (l *Lizard) Age() int {
	return l.age
}

func (l *Lizard) SetAge(age int) {
	l.age = age
}

func (l *Lizard) Crawl() {
	if l.age < 10 {
		fmt.Println("Crawling!")
	}
}

// So what is the situation here and how does it relate to the decorator?
// Well, in the dragon class that you see here, the dragon structure that we have
// constructed, a decorator, we have constructed an object which extends the behaviors
// of the types that we have right here. But what it's doing really is it's providing
// better access to the underlying fields of both the bird and the lizard,
type Dragon struct {
	bird   Bird
	lizard Lizard
}

func (d *Dragon) Age() int {
	return d.bird.age
}

func (d *Dragon) SetAge(age int) {
	d.bird.SetAge(age)
	d.lizard.SetAge(age)
}

func (d *Dragon) Fly() {
	d.bird.Fly()
}

func (d *Dragon) Crawl() {
	d.lizard.Crawl()
}

func NewDragon() *Dragon {
	return &Dragon{Bird{}, Lizard{}}
}

func main() {
	//d := Dragon{}
	//d.Bird.Age = 10
	//fmt.Println(d.Lizard.Age)
	//d.Fly()
	//d.Crawl()

	d := NewDragon()
	d.SetAge(10)
	d.Fly()
	d.Crawl()
}
