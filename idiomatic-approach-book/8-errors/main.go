package main

import (
	"archive/zip"
	"bytes"
	"fmt"
)

func main() {
	res, err := doubleEven(3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	zipReader()
}

func doubleEven(i int) (int, error) {
	if i%2 != 0 {
		return 0, fmt.Errorf("%d is not an even number", i)
	}

	return i * 2, nil
}

// *********************************   Sentinel errors *********************************
// Sentinel errors are usually used to indicate that you cannot start or continue processing
func zipReader() {
	data := []byte("this is not a zip file or data")
	notAZipFile := bytes.NewReader(data)
	_, err := zip.NewReader(notAZipFile, int64(len(data)))
	if err == zip.ErrFormat {
		fmt.Println("not a zip file")
	}

}
