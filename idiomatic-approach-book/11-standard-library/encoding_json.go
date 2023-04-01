package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {
	// encodingAndDecodingStreams()

	customJsonParsing()
}

// ********************************* Encoding *********************************
//  ┌
//  │       If a field should be ignored when marshaling or unmarshaling, use - for the name.
//  │ If the field should be left out of the output when it is empty, add ,omitempty after the name.
//  └

type Order struct {
	ID          string      `json:"id"`
	DateOrdered RFC822ZTime `json:"date_ordered"`
	CustomerID  string      `json:"customer_id"`
	Items       []Item      `json:"items"`
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

// Encoding and Decoding JSON Streams
func encodingAndDecodingStreams() {
	const data = `
		{"name": "Fred", "age": 40}
		{"name": "Mary", "age": 21}
		{"name": "Pat", "age": 30}
	`
	var t struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	dec := json.NewDecoder(strings.NewReader(data))
	var b bytes.Buffer
	enc := json.NewEncoder(&b)

	// dec.More() This lets us read the data in, one JSON object at a time.
	for dec.More() {
		err := dec.Decode(&t)
		if err != nil {
			panic(err)
		}
		fmt.Println(t) // {Fred 40} then:  {Mary 21} then: {Pat 30} sequentially
		err = enc.Encode(t)
		// Writing out multiple values with the json.Encoder works just like using it to write out a single value
		if err != nil {
			panic(err)
		}
	}
	out := b.String()
	fmt.Println(out)
}

// Custom JSON Parsing
type RFC822ZTime struct {
	time.Time
}

func (rt RFC822ZTime) MarshalJSON() ([]byte, error) {
	out := rt.Time.Format(time.RFC822Z)
	return []byte(`"` + out + `"`), nil
}

func (rt *RFC822ZTime) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		return nil
	}
	t, err := time.Parse(`"`+time.RFC822Z+`"`, string(b))
	//  the method that reads the time value is declared on a value receiver
	if err != nil {
		return err
	}
	*rt = RFC822ZTime{t}
	// while the method that modifies the the time value is declared on a pointer receiver.
	return nil
}

func customJsonParsing() {
	data := `
	{
		"id": "12345",
		"items": [
			{
				"id": "xyz123",
				"name": "Thing 1"
			},
			{
				"id": "abc789",
				"name": "Thing 2"
			}
		],
		"date_ordered": "01 May 20 13:01 +0000",
		"customer_id": "3"
	}`

	var o Order
	err := json.Unmarshal([]byte(data), &o)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", o)
	fmt.Println(o.DateOrdered.Month())
	out, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
}
