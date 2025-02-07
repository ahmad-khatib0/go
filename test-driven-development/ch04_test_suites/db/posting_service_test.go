package db_test

import (
	"testing"

	"github.com/ahmad-khatib0/go/test-driven-development/ch04_test_suites/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	ps := db.NewPostingService()
	b := db.Book{
		ID: uuid.New().String(),
	}
	err := ps.NewOrder(b)
	assert.Nil(t, err)
}
