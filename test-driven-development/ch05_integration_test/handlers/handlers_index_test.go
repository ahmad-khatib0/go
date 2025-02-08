package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/ahmad-khatib0/go/test-driven-development/ch05_integration_test/db"
	"github.com/ahmad-khatib0/go/test-driven-development/ch05_integration_test/handlers"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// The Ginkgo equivalent of our Index integration test
//
// - We make use of closures to set up our spec hierarchy. The Describe function allows us to create
//   container nodes. Specs must begin with a top-level Describe node.
//
// - The BeforeEach function creates setup nodes that run before tests. They are used for extracting
//   common setups, allowing us to streamline our tests.
//
// - The AfterEach function creates setup nodes that run after tests. They allow us to clean up after
//   our specs have run, ensuring that critical resources are cleaned up correctly.
//
// - We can further define container nodes inside the top-level nodes as required to organize our
//   specs and their scenarios.
//
// - The Context function is an alias for Describe that allows us to add extra information to
//   our specs to help people understand them. It also creates container nodes but can be used to
//   organize information.
//
// - The It function allows us to define subject nodes. These nodes contain the assertions of the
//   subject under test and cannot contain any other nested nodes.
//
// - The assertions inside the subject nodes are written with the gomega assertion library. These can
//   be nested just like the assertions of testify but take a human-readable form. All assertions
//   must begin with the Expect function, which wraps an actual value.

var _ = Describe("Handlers integration", func() {
	var svr *httptest.Server
	var eb db.Book

	BeforeEach(func() {
		eb = db.Book{
			ID:     uuid.New().String(),
			Name:   "My first integration test",
			Status: db.Available.String(),
		}

		bs := db.NewBookService([]db.Book{eb}, nil)
		ha := handlers.NewHandler(bs, nil)
		svr = httptest.NewServer(http.HandlerFunc(ha.Index))
	})

	AfterEach(func() {
		svr.Close()
	})

	Describe("Index endpoint", func() {
		Context("with one existing book", func() {
			It("should return book", func() {
				r, err := http.Get(svr.URL)
				Expect(err).To(BeNil())
				Expect(r.StatusCode).To(Equal(http.StatusOK))

				body, err := io.ReadAll(r.Body)
				r.Body.Close()
				Expect(err).To(BeNil())

				var resp handlers.Response
				err = json.Unmarshal(body, &resp)

				Expect(err).To(BeNil())
				Expect(len(resp.Books)).To(Equal(1))
				Expect(resp.Books).To(ContainElement(eb))
			})
		})
	})
})
