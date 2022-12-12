package main

import "fmt"

func main() {

	languages := make(map[string]string)
	languages["js"] = "javascript"
	languages["css"] = "style sheet"
	languages["rb"] = "ruby"
	languages["cpp"] = "c plus plus"
	fmt.Println("list of languages: ", languages) // map[cpp:c plus plus css:style sheet js:javascript rb:ruby]
	fmt.Println("js is the format of: ", languages["js"])

	delete(languages, "css") //delete
	fmt.Println("list of languages: ", languages)

	for key, value := range languages {
		fmt.Printf("for key %v, value is %v\n", key, value)
	}
}
