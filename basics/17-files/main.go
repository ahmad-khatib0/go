package main

import (
	"fmt"
	"io"
	"os"
)

func main() {

	content := "Plcae this inside a file with golang"
	file, err := os.Create("./content.txt")

	checkErrors(err)
	length, err := io.WriteString(file, content)

	checkErrors(err)
	fmt.Println("length is: ", length)
	defer file.Close()

	readFile("./content.txt")
}

func readFile(filename string) {
	dataBytes, err := os.ReadFile(filename)
	checkErrors(err)

	fmt.Println("data in the file: ", string(dataBytes))

}

func checkErrors(err error) {
	if err != nil {
		panic(err)
	}

}
