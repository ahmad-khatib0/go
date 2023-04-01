package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func main() {
	httpRequest()
}

func httpRequest() {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://jsonplaceholder.typicode.com/todos/1", nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req.Header.Add("X-My-Client", "Learning Go")

	res, err := client.Do(req) // send it
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("unexpected status: got %v", res.Status))
	}

	fmt.Println(res.Header.Get("Content-Type")) // application/json; charset=utf-8
	var data struct {
		UserID    int    `json:"userId"`
		ID        int    `json:"id"`
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", data) // { UserID:1 ID:1 Title:delectus aut autem Completed:false }
}
