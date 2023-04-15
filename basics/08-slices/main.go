package main

import (
	"fmt"
	"sort"
)

func main() {

	// var fruitList = []string{} //or
	var fruitList = []string{"tomatos", "peach", "apple"}
	fmt.Printf("the value of slice is: %T", fruitList) //string[]

	fruitList = append(fruitList, "Banana", "Mango") // add
	fmt.Println("value is now after adding is: ", fruitList)

	fruitList = append(fruitList[1:]) // start from postion 1, and delete, also, [1:3] start from
	// postion and delete untill the postion 3, not including this postion
	fmt.Println("value is now after deleting is: ", fruitList)

	highScores := make([]int, 4)
	highScores[0] = 333
	highScores[1] = 433
	highScores[2] = 533
	highScores[3] = 633
	// highScores[4] = 633 // will breaks, because its out of range,
	highScores = append(highScores, 44, 55, 66) // will work, because it will realocate the address in memory
	fmt.Println(highScores)

	sort.Ints(highScores) //these methods are available in Slices, not in arrays
	fmt.Println(highScores)
	fmt.Println(sort.IntsAreSorted(highScores)) //true

	// remove from slices
	var courses = []string{"js", "py", "rb", "cpp"}
	var toBeRemoved int = 2
	courses = append(courses[:toBeRemoved], courses[toBeRemoved+1:]...) // like concatenating the tow parts of the slice
	fmt.Println(courses)                                                // [js py cpp]
}
