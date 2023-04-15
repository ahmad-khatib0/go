package main

import "fmt"

func main() {

	days := []string{"Sunaday", "Monday", "Wednesday", "Friday", "Saturday"}

	for d := 0; d < len(days); d++ {
		fmt.Println(days[d])
	}

	for i := range days { // i is the index not the value
		fmt.Println(days[i])
	}

	for i, val := range days { // i is the index not the value
		fmt.Printf("index is %v and value is %v\n ", i, val)
	}

	likeWhile := 1
	for likeWhile < 10 {

		if likeWhile == 3 {
			fmt.Println("value is now", likeWhile)
			likeWhile++
			continue
		}

		if likeWhile == 6 {
			goto jumpToMe
		}

		if likeWhile == 8 {
			break
		}
		fmt.Println("value is: ", likeWhile)
		likeWhile++
	}

jumpToMe:
	fmt.Println("jumpped to it inside the for loop and it will break the loop after jumpping")
}
