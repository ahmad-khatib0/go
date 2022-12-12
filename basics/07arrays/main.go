package main

import "fmt"

func main() {

	fmt.Println("arrays")

	var fruitList [4]string // 4 elements
	fruitList[0] = "bananna"
	fruitList[1] = "apple"
	fruitList[3] = "tomatos"

	fmt.Println("value of fruitList array with a white space between the 3th and 4th element", fruitList)
	fmt.Println("length of the array is 4, the total reserved indices, NOT how many elements are there", len(fruitList)) // 4
	//[bananna apple  tomatos]  NOTICE: notice how there is a whitespace

	var vegatablesList = [3]string{"cucumber", "potattos", "tomatos"}
	fmt.Println("vegatables list is: ", vegatablesList)

}
