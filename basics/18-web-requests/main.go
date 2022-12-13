package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

const url = "https://lco.dev"

func main() {
	fmt.Println("web request")

	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	fmt.Printf("response is of type:  %T\n", res)
	defer res.Body.Close() // developer responsibility to close the connection ALWAYS

	dataBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dataBytes))

}
