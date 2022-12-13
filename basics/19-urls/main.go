package main

import (
	"fmt"
	"net/url"
)

const api string = "https://lco.dev:3000/learn?coursename=reactjs"

func main() {

	result, _ := url.Parse(api)
	fmt.Println(result.Scheme)
	fmt.Println(result.Host)
	fmt.Println(result.Path)
	fmt.Println(result.RawQuery)
	fmt.Println(result.Port())

	queryParams := result.Query()
	fmt.Println(queryParams["coursename"])

	for key, val := range queryParams {
		fmt.Printf("value is: %v , for the key: %v\n ", val, key)
	}

	partsOfUrl := &url.URL{
		Scheme:  "https",
		Host:    "lco.dev",
		Path:    "/learn",
		RawPath: "user=ahmad",
	}

	constructedUrl := partsOfUrl.String()
	fmt.Println(constructedUrl)
}
