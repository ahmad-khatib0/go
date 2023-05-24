package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ahmadkhatib0/go/snippetbox/pkg/models/mock"
	"github.com/golangcollege/sessions"
)

// testServer type which anonymously embeds a httptest.Server instance.
type testServer struct {
	*httptest.Server
}

// newTestApplication() returns an instance of our application struct containing mocked dependencies.
func newTestApplication(t *testing.T) *application {

	templateCache, err := newTemplateCache("./../../ui/html/")
	if err != nil {
		t.Fatal(err)
	}

	session := sessions.New([]byte("3dSm5MnygFHh7XidAtbskXrjbwfoJcbJ"))
	session.Lifetime = 12 * time.Hour
	session.Secure = true

	// the reason for mocking these and writing to ioutil.Discard is to avoid clogging
	// up our test output with unnecessary log messages when we run go test -v.
	return &application{
		errorLog:      log.New(ioutil.Discard, "", 0),
		infoLog:       log.New(ioutil.Discard, "", 0),
		session:       session,
		templateCache: templateCache,
		snippets:      &mock.SnippetModel{},
		users:         &mock.UserModel{},
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	// we added jar, so now response cookies are stored and then sent with subsequent requests.
	ts.Client().Jar = jar

	// Disable redirect-following for the client. Essentially this function is called after
	// a 3xx response is received by the client, and returning the http.ErrUseLastResponse error
	// forces it to immediately return the received response.
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, []byte) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	body, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, body
}
