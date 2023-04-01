package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// ********************************* Encoding *********************************
//  ┌
//  │       If a field should be ignored when marshaling or unmarshaling, use - for the name.
//  │ If the field should be left out of the output when it is empty, add ,omitempty after the name.
//  └

type Order struct {
	ID          string    `json:"id"`
	DateOrdered time.Time `json:"date_ordered"`
	CustomerID  string    `json:"customer_id"`
	Items       []Item    `json:"items"`
}
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// convert data to a struct of type Order :
func Unmarshaling(data string) error {
	var o Order
	err := json.Unmarshal([]byte(data), &o) // If o is nil or not a pointer, Unmarshal returns an InvalidUnmarshalError
	if err != nil {
		return err
	}
	return nil
}

// The json.Decoder and json.Encoder types read from and write to anything that meets the
// io.Reader and io.Writer interfaces, respectively
func toFile() {
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	toFile := Person{
		Name: "Fred",
		Age:  40,
	}

	tmpFile, err := ioutil.TempFile(os.TempDir(), "sample-")
	if err != nil {
		panic(err)
	}

	defer os.Remove(tmpFile.Name())

	err = json.NewEncoder(tmpFile).Encode(toFile)
	if err != nil {
		panic(err)
	}
	err = tmpFile.Close()
	if err != nil {
		panic(err)
	}

	// Now  we can read the JSON back in by passing a reference to the temp file to json.NewDecoder
	tmpFile2, err := os.Open(tmpFile.Name())
	if err != nil {
		panic(err)
	}
	var fromFile Person
	err = json.NewDecoder(tmpFile2).Decode(&fromFile)
	if err != nil {
		panic(err)
	}
	err = tmpFile2.Close()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", fromFile) // => {Name:Fred Age:40}

}
