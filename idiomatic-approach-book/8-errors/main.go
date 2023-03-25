package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"os"
)

func main() {
	res, err := doubleEven(3)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)

	zipReader()

	errorsAreValues("3345", "pwd", "filename")

	err = uninitializedInstance(true)
	fmt.Println(err != nil) // prints true , which is expected
	err = uninitializedInstance(false)
	fmt.Println(err != nil) // prints false!, which is wrong

	//  ╔══════════════════════════════════════════════════════════════════════════════════════════╗
	//  ║   The reason why err is non-nil is that error is a interface.  for an interface to be    ║
	//  ║      considered nil, both the underlying type and the underlying value must be nil.      ║
	//  ║ Whether or not genErr is a pointer, the underlying type part of the interface is not nil ║
	//  ╚══════════════════════════════════════════════════════════════════════════════════════════╝

	err = wrappingErrors("non_existnet.txt")
	if err != nil {
		fmt.Println(err) // in wrappingErrors: open non_existnet.txt: no such file or directory
		if wrappedError := errors.Unwrap(err); wrappedError != nil {
			fmt.Println(err) // in wrappingErrors: open non_existnet.txt: no such file or directory
		}
	}

}

func doubleEven(i int) (int, error) {
	if i%2 != 0 {
		return 0, fmt.Errorf("%d is not an even number", i)
	}

	return i * 2, nil
}

//     *********************************   Sentinel errors *********************************
//  ▲
//  █ Sentinel errors are usually used to indicate that you cannot start or continue processing
//  ▼

func zipReader() {
	data := []byte("this is not a zip file or data")
	notAZipFile := bytes.NewReader(data)
	_, err := zip.NewReader(notAZipFile, int64(len(data)))
	if err == zip.ErrFormat {
		fmt.Println("not a zip file")
	}

}

// ********************************* Errors are values *********************************
type Status int

const (
	InvalidLogin Status = iota + 1
	NotFound
)

type StatusErr struct {
	Status  Status
	Message string
	Err     error
}

func (se StatusErr) Error() string {
	return se.Message
}

func errorsAreValues(uid, pwd, file string) ([]byte, error) {
	err := login(uid, pwd)
	if err != nil {
		return nil, StatusErr{
			Status:  InvalidLogin,
			Message: fmt.Sprintf("invalid credentials for user %s", uid),
			Err:     err,
		}
	}
	data, err := getData(file)
	if err != nil {
		return nil, StatusErr{
			Status:  NotFound,
			Message: fmt.Sprintf("file %s not found", file),
			Err:     err,
		}
	}
	return data, nil
}

func login(uid, pwd string) error {
	return nil
}
func getData(file string) ([]byte, error) {
	return nil, errors.New("file not found")
}

// DON'T use this pattern  (DON’T RETURN AN UNINITIALIZED INSTANCE)
func uninitializedInstance(flag bool) error {
	var genErr StatusErr
	if flag {
		genErr = StatusErr{
			Status: NotFound,
		}
	}
	return genErr
}

// *********************************    Wrapping Errors  *********************************
func wrappingErrors(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("in wrappingErrors: %w", err)
	}

	f.Close()
	return nil
}

// -- Wrap an error with custom error type
//  ▲
//  █ If you want to wrap an error with your custom error type, your error type needs to implement the method Unwrap
//  ▼

func (se StatusErr) Unwrap() error {
	return se.Err
}
