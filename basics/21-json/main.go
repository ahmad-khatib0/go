package main

import (
	"encoding/json"
	"fmt"
)

type course struct {
	Name     string `json:"coursename"` // renaming
	Price    int
	Platform string   `json:"website"`
	Password string   `json:"-"`              // hide password
	Tags     []string `json:"tags,omitempty"` // hide it if its nil
}

func main() {
	// EncodeJson()
	DecodeJson()
}

func EncodeJson() {

	lcoCourses := []course{
		{"ReactJS BootCamp", 299, "LearnCodeOnline.in", "abc123", []string{"web-dev", "js"}},
		{"MERN BootCamp", 199, "LearnCodeOnline.in", "bcd123", []string{"full-stack", "js"}},
		{"Angular BootCamp", 299, "LearnCodeOnline.in", "hit123", nil},
	}

	finalJson, _ := json.MarshalIndent(lcoCourses, "", "\t")
	fmt.Printf("%s\n", finalJson)
}

func DecodeJson() {
	data := []byte(`
			{
				"coursename": "ReactJS BootCamp",
				"Price": 299,
				"website": "LearnCodeOnline.in",
				"Tags": [ "web-dev", "js" ]
			}
  `)

	var lcoCourse course
	isJsonValid := json.Valid(data)

	if isJsonValid {
		json.Unmarshal(data, &lcoCourse)
		fmt.Printf("%#v\n", lcoCourse)
	} else {
		fmt.Println("json is not valid")
	}

	var onlineData map[string]interface{}
	json.Unmarshal(data, &onlineData)
	fmt.Printf("%#v\n", onlineData)

	for k, v := range onlineData {
		fmt.Printf("key is %v and value is %v , type is: %T\n", k, v, v)
	}
}
