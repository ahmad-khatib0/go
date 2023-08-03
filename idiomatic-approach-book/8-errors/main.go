package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"internal/itoa"
	"os"
	"reflect"
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

	err = isAndAs("non_existnet.txt")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("this file is not existed")
		}
	}

	// Now we can find, for example, all errors that refer to the database, no matter the code:
	if errors.Is(err, ResourceErr{Resource: "Database"}) {
		fmt.Println("the database is broken", err)
	}

	var coder interface {
		Code() int
	}
	if errors.As(err, &coder) {
		fmt.Println(coder.Code())
	}

	// doPanic(os.Args[0]) // prints stacktrace
	for _, val := range []int{1, 2, 0, 6} {
		div60Recover(val) //  60 \n 30 \n runtime error: integer divide by zero \n 10
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

// *********************************   Is and As  *********************************
func isAndAs(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("in wrappingErrors: %w", err)
	}

	f.Close()
	return nil
}

// ▲
// █   Custom Is checker for (non-comparable type)  because Is By default
// █   uses == to compare each wrapped error with the specified error
// ▼
type MyErr struct {
	Codes []int
}

func (me MyErr) Error() string {
	return fmt.Sprintf("codes: %v", me.Codes)
}

func (me MyErr) Is(target error) bool {
	if me2, ok := target.(MyErr); ok {
		return reflect.DeepEqual(me, me2)
	}
	return false
}

// ▲
// █   Another use for defining your own Is method is to allow
// █   comparisons against errors that aren’t identical instances.
// ▼
type ResourceErr struct {
	Resource string
	Code     int
}

func (re ResourceErr) Error() string {
	return fmt.Sprintf("%s: %d", re.Resource, re.Code)
}

func (re ResourceErr) Is(target error) bool {
	if other, ok := target.(ResourceErr); ok {
		ignoreResource := other.Resource == ""
		ignoreCode := other.Code == 0
		matchResource := other.Resource == re.Resource
		matchCode := other.Code == re.Code
		return matchResource && matchCode || matchResource && ignoreCode || ignoreResource && matchCode
	}
	return false
}

//  ▲
//  █  The errors.As function allows you to check if a returned
//  █  error (or any error it wraps) matches a specific type.
//  ▼

// ****************************** Wrapping Errors with defer *************************

func wrappingErrorsWithDefer(val1 int, val2 string) (_ string, err error) {
	// NOTE: We have to name our return values, so that we can refer to err in the deferred function.
	defer func() {
		if err != nil {
			err = fmt.Errorf("in wrappingErrorsWithDefer: %w", err)
		}
	}()

	val3, err := doThing1(val1)
	if err != nil {
		return "", err
	}
	val4, err := doThing2(val2)
	if err != nil {
		return "", err
	}
	return doThing3(itoa.Itoa(val3), val4)
}

func doThing1(val int) (int, error)                    { return 0, nil }
func doThing2(val string) (string, error)              { return "", nil }
func doThing3(val string, val2 string) (string, error) { return "", nil }

// the above cleaner code is equivalent to:
func DoSomeThings(val1 int, val2 string) (string, error) {
	val3, err := doThing1(val1)
	if err != nil {
		return "", fmt.Errorf("in DoSomeThings: %w", err)
	}
	val4, err := doThing2(val2)
	if err != nil {
		return "", fmt.Errorf("in DoSomeThings: %w", err)
	}
	result, err := doThing3(itoa.Itoa(val3), val4)
	if err != nil {
		return "", fmt.Errorf("in DoSomeThings: %w", err)
	}
	return result, nil
}

// ****************************** panic and recover *************************
func doPanic(msg string) {
	panic(msg)
}

// Go provides a way to capture a panic in order to provide a more graceful shutdown or to prevent shutdown at all.
func div60Recover(i int) {
	defer func() {
		if v := recover(); v != nil {
			fmt.Println(v)
		}
	}()

	fmt.Println(60 / i)
}
