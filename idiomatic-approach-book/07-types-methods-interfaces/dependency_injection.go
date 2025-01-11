package main

import (
	"errors"
	"fmt"
	"net/http"
)

// ╭──────────────────────────────────────────────────────╮
// │ Implicit Interfaces Make Dependency Injection Easier │
// ╰──────────────────────────────────────────────────────╯

func LogOutput(message string) {
	fmt.Println(message)
}

type SimpleDataStore struct {
	userData map[string]string
}

func (sds SimpleDataStore) UserNameForID(userID string) (string, bool) {
	name, ok := sds.userData[userID]
	return name, ok
}

// a factory function to create an instance of a SimpleDataStore
func NewSimpleDataStore() SimpleDataStore {
	return SimpleDataStore{
		userData: map[string]string{
			"1": "Fred",
			"2": "Mary",
			"3": "Pat",
		},
	}
}

// these interfaces to avoid depending on n LogOutput or SimpleDataStore
type DataStore interface {
	UserNameForID(userID string) (string, bool)
}
type Logger interface {
	Log(message string)
}

// making the LogOutput meets this interface (Logger)
type LoggerAdapter func(message string)

func (lg LoggerAdapter) Log(message string) {
	lg(message)
}

// The business logic
type SimpleLogic struct {
	l  Logger
	ds DataStore
}

// We have a struct with two fields, one a Logger, the other a DataStore. There’s nothing in our
// SimpleLogic that mentions the concrete types, so there’s no dependency on them. There’s no
// problem if we later swap in new implementations from an entirely different provider,

func (sl SimpleLogic) SayHello(userID string) (string, error) {
	sl.l.Log("in SayHello for " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("unknown user")
	}
	return "Hello, " + name, nil
}
func (sl SimpleLogic) SayGoodbye(userID string) (string, error) {
	sl.l.Log("in SayGoodbye for " + userID)
	name, ok := sl.ds.UserNameForID(userID)
	if !ok {
		return "", errors.New("unknown user")
	}
	return "Goodbye, " + name, nil
}

// When we want a SimpleLogic instance, we call a
// factory function, passing in interfaces and returning a struct
func NewSimpleLogic(l Logger, ds DataStore) SimpleLogic {
	return SimpleLogic{
		l:  l,
		ds: ds,
	}
}

// Our controller needs business logic that says hello, so we define an interface for that:
type LogicInterface interface {
	SayHello(userID string) (string, error)
}

type Controller struct {
	l     Logger
	logic LogicInterface
}

func (c Controller) HandleGreeting(w http.ResponseWriter, r *http.Request) {
	c.l.Log("In SayHello")
	userID := r.URL.Query().Get("user_id")
	message, err := c.logic.SayHello(userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(message))
}

func NewController(l Logger, logic LogicInterface) Controller {
	return Controller{
		l:     l,
		logic: logic,
	}
}

func main1() {
	// The main function is the only part of the code that knows what all the concrete types actually are.
	// If we want to swap in different implementations, this is the only place that needs to change
	l := LoggerAdapter(LogOutput)
	ds := NewSimpleDataStore()

	logic := NewSimpleLogic(l, ds)
	c := NewController(l, logic)

	http.HandleFunc("/hello", c.HandleGreeting) // NOTE: this should be SayHello, but it gives an error
	http.ListenAndServe(":8080", nil)
}
