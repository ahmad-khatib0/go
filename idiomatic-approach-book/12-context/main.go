package main

import (
	"context"
	"fmt"
	"net/http"
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
