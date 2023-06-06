package main

import "fmt"

func callByValue() {
	p := person{}
	i := 2
	s := "hello"
	modifyFails(i, s, p)
	fmt.Println(i, s, p)

	m := map[int]string{1: "first", 2: "second"}
	modMap(m)
	fmt.Println(m) // map[2:hello 3:goodbye]
	e := []int{1, 2, 3}
	modSlice(e)
	fmt.Println(e) // [2 4 6]
}

type person struct {
	age  int
	name string
}

func modifyFails(i int, s string, p person) {
	i = i * 2
	s = "Goodbye"
	p.name = "Bob"
}

func modMap(m map[int]string) {
	m[2] = "hello"
	m[3] = "goodbye"
	delete(m, 1)
}
func modSlice(s []int) {
	for k, v := range s {
		s[k] = v * 2
	}
	s = append(s, 10)
}
