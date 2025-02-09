package db_test

import (
	"testing"

	"github.com/ahmad-khatib0/go/test-driven-development/ch09_concurrency_test/db"
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
