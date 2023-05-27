package main

import (
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Ahmadkhatib0/go/greenlight/internal/data"
	"github.com/felixge/httpsnoop"
	"github.com/pascaldekloe/jwt"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// deferred function (which will always be run in the event of a panic as Go unwinds the stack).
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu sync.Mutex
		// Update the map so the values are pointers to a client struct.
		clients = make(map[string]*client)
	)

	// Launch a background goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to prevent any rate limiter checks from happening while the cleanup is taking place.
			mu.Lock()

			// If they haven't been seen within the last three minutes, delete it
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only carry out the check if rate limiting is enabled.
		if app.config.limiter.enabled {
			ip := realip.FromRequest(r)

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			// Very importantly, unlock the mutex before calling the next handler in the chain. Notice that we
			// ( DON'T USE DEFER TO UNLOCK THE MUTEX), as that would mean that the mutex isn't unlocked until
			// ALL the handlers downstream of this middleware have also returned.
			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		// Parse the JWT and extract the claims. This will return an error if the JWT contents doesn't
		// match the signature (i.e. the token has been tampered with) or the algorithm isn't valid.
		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.jwt.secret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check if the JWT is still valid at this moment in time.
		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check that the issuer is our application.
		if claims.Issuer != "greenlight.alexedwards.net" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Check that our application is in the expected audiences for the JWT.
		if !claims.AcceptAudience("greenlight.alexedwards.net") {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// At this point, we know that the JWT is all OK and we can trust the data in it.
		// We extract the user ID from the claims subject and convert it from a string into an int64.
		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Lookup the user record from the database.
		user, err := app.models.Users.Get(userID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Add the user record to the request context and continue as normal.
		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})

}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	// Rather than returning this http.HandlerFunc we assign it to the variable fn.
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})

	// Wrap fn with the requireAuthenticatedUser() middleware before returning it.
	return app.requireAuthenticatedUser(fn)
}

// requireActivatedUser() middleware automatically calls the requireAuthenticatedUser()
// middleware before being executed itself
func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	// Wrap this with the requireActivatedUser() middleware before returning it.
	return app.requireActivatedUser(fn)
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")
		// vary because response will be different depending on whether or not this header exists in the request.
		w.Header().Add("Vary", "Access-Control-Request-Method")
		origin := r.Header.Get("Origin")
		if origin != "" {
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// Write the headers along with a 200 OK status and return from the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
					break
				}

			}
		}

		next.ServeHTTP(w, r)
	})
}

// func (app *application) metrics(next http.Handler) http.Handler {
// 	// Initialize the new expvar variables when the middleware chain is first built.
// 	totalRequestsReceived := expvar.NewInt("total_requests_received")
// 	totalResponsesSent := expvar.NewInt("total_responses_sent")
// 	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

// 	// The following code will be run for every request...
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()
// 		totalRequestsReceived.Add(1)

// 		next.ServeHTTP(w, r)
// 		totalResponsesSent.Add(1)

// 		// Calculate the number of microseconds since we began to process the request,
// 		// then increment the total processing time by this amount.
// 		duration := time.Since(start).Microseconds()
// 		totalProcessingTimeMicroseconds.Add(duration)
// 	})
// }

func (app *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

	// Declare a new expvar map to hold the count of responses for each HTTP status code.
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequestsReceived.Add(1)

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalResponsesSent.Add(1)

		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}
