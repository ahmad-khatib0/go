package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

// func Test_updateMessage(t *testing.T) {
// 	msg = "Hello, world!"

// 	wg.Add(2)
// 	go updateMessage("x")
// 	go updateMessage("Goodbye, cruel world!")
// 	wg.Wait()

// 	if msg != "Goodbye, cruel world!" {
// 		t.Error("incorrect value in msg")
// 	}
// }

func Test_main(t *testing.T) {
	stdOut := os.Stdout
	r, w, _ := os.Pipe()

	os.Stdout = w

	main()

	_ = w.Close()

	result, _ := io.ReadAll(r)
	output := string(result)

	os.Stdout = stdOut

	if !strings.Contains(output, "$34320.00") {
		t.Error("wrong balance returned")
	}

}
