package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

func Middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		ctx := req.Context() //  extract the existing context from the request using the Context method
		req = req.WithContext(ctx)
		// we create a new request based on the old request and the now-populated context using the WithContext method
		handler.ServeHTTP(rw, req)

	})
}

func logic(ctx context.Context, data string) (string, error) { return "", nil }
func handler(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	err := req.ParseForm()
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	data := req.FormValue("data")
	result, err := logic(ctx, data)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	rw.Write([]byte(result))
}

// Passing context to make an http request to another server
type ServiceCaller struct {
	client *http.Client
}

func (sc ServiceCaller) callAnotherService(ctx context.Context, data string) (string, error) {
	req, err := http.NewRequest(http.MethodGet,
		"http://example.com?data="+data, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	resp, err := sc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Unexpected status code %d",
			resp.StatusCode)
	}
	// do the rest of the stuff to process the response
	return "done", err
}

// Timers

func timers() {

	ctx := context.Background()
	parent, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	child, cancel2 := context.WithTimeout(parent, 3*time.Second)
	defer cancel2()

	start := time.Now()
	<-child.Done()
	// Done returns channel of struct{} ( Why struct ?  because  an empty struct uses no memory.)

	end := time.Now()
	fmt.Println(end.Sub(start))
}

func longRunningThing(ctx context.Context, data string) (string, error) { return "", nil }

func longRunningThingManager(ctx context.Context, data string) (string, error) {
	type wrapper struct {
		result string
		err    error
	}
	ch := make(chan wrapper, 1)
	go func() {
		// do the long running thing
		result, err := longRunningThing(ctx, data)
		ch <- wrapper{result, err}
	}()
	select {
	case data := <-ch:
		return data.result, data.err
	case <-ctx.Done(): // triggered if the context is canceled
		return "", ctx.Err()
	}
}

// Value  (pass values between contexts)

type userKey int

const key userKey = 1

func ContextWithUser(ctx context.Context, user string) context.Context {
	return context.WithValue(ctx, key, user)
}

func UserFromContext(ctx context.Context) (string, bool) {
	user, ok := ctx.Value(key).(string)
	return user, ok
}

func extractUser(req *http.Request) (string, error) {
	userCookie, err := req.Cookie("user")
	if err != nil {
		return "", err
	}
	return userCookie.Value, nil
}

func MyMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		user, err := extractUser(req)
		if err != nil {
			rw.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := req.Context()
		ctx = ContextWithUser(ctx, user)
		req = req.WithContext(ctx)
		h.ServeHTTP(rw, req)
	})
}
