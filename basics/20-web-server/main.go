package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// getReq()
	// postReq()
	postForm()

}

func getReq() {
	const url string = "http://localhost:8000/get"
	res, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	fmt.Println("status code: ", res.StatusCode)
	fmt.Println("content length ", res.ContentLength)

	var resString strings.Builder
	content, _ := ioutil.ReadAll(res.Body)
	byteCount, _ := resString.Write(content)

	fmt.Println("byteCount is: ", byteCount)

	// fmt.Println(string(content))  // or
	fmt.Println(resString.String())

}

func postReq() {
	const myurl = "http://localhost:8000/post"

	requestBody := strings.NewReader(`
		{
			"coursename":"Let's go with golang",
			"price": 0,
			"platform":"learnCodeOnline.in"
		}
	`)

	response, err := http.Post(myurl, "application/json", requestBody)

	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	content, _ := ioutil.ReadAll(response.Body)

	fmt.Println(string(content))
}

func postForm() {

	const myurl = "http://localhost:8000/postform"

	data := url.Values{}
	data.Add("firstname", "ahmad")
	data.Add("lastname", "programmer")
	data.Add("email", "ahmad@test.com")

	res, err := http.PostForm(myurl, data)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	content, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(content))
}
