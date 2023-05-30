package main

import "fmt"

func main() {

	fmt.Println("arrays")

	var fruitList [4]string // 4 elements
	fruitList[0] = "banana"
	fruitList[1] = "apple"
	fruitList[3] = "tomatoes"

	fmt.Println("value of fruitList array with a white space between the 3th and 4th element", fruitList)
	fmt.Println("length of the array is 4, the total reserved indices, NOT how many elements are there", len(fruitList)) // 4
	//[banana apple  tomatoes]  NOTICE: notice how there is a whitespace

	var vegetablesList = [3]string{"cucumber", "potato's", "tomatoes"}
	fmt.Println("vegetables list is: ", vegetablesList)

}
