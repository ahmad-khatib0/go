package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)
	go func() {
		// We need to use a buffered channel here because signal.Notify() does not wait for a
		// receiver to be available when sending a signal to the quit channel. If we had used a regular
		// (non-buffered) channel here instead, a signal could be ‘missed’ if our quit channel is not
		// ready to receive at the exact moment that the signal is sent. By using a buffered channel, we
		// avoid this problem and ensure that we never miss a signal.
		quit := make(chan os.Signal, 1)
		// Use signal.Notify() to listen for incoming SIGINT and SIGTERM signals and relay them to the quit channel
		// Any other signals will not be caught by signal.Notify() and will retain their default behavior.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// Read the signal from the quit channel. This code will block until a signal is received.
		s := <-quit

		app.logger.PrintInfo("caught signal", map[string]string{"signal": s.String()})
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)

		// Exit the application with a 0 (success) status code.
		// os.Exit(0)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	err := srv.ListenAndServe()
	// Calling Shutdown() on our server will cause ListenAndServe() to immediately
	// return a http.ErrServerClosed error. So if we see this error, it is actually a
	// good thing and an indication that the graceful shutdown has started. So we check
	// specifically for this, only returning the error if it is NOT http.ErrServerClosed.
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Otherwise, we wait to receive the return value from Shutdown() on the shutdownError channel. If return
	// value is an error, we know that there was a problem with the graceful shutdown and we return the error.
	err = <-shutdownError
	if err != nil {
		return err
	}

	// here we know that the graceful shutdown completed successfully
	app.logger.PrintInfo("stopped server", map[string]string{"addr": srv.Addr})

	return nil
}
