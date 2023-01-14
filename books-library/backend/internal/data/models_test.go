package data

import "testing"

func Test_Ping(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("failed to ping db")
	}
}

func TestBook_GetAll(t *testing.T) {
	all, err := models.Book.GetAll()
	if err != nil {
		t.Error("failed to get all books ", err)
	}

	if len(all) != 1 {
		t.Error("failed to get the correct number of books")
	}
}

func TestBook_GetOneById(t *testing.T) {
	b, err := models.Book.GetOneById(1)
	if err != nil {
		t.Error("failed to get the book by id", err)
	}

	if b.Title != "My Book" {
		t.Errorf("expected title to be My Book but got %s ", b.Title)
	}
}

func TestBook_GetOneBySlug(t *testing.T) {
	b, err := models.Book.GetOneBySlug("my-book")
	if err != nil {
		t.Error("failed to get the book by slug", err)
	}

	if b.Title != "My Book" {
		t.Errorf("expected title to be My Book but got %s ", b.Title)
	}

	_, err = models.Book.GetOneBySlug("bad-slug")
	if err == nil {
		t.Error("did not get an error when attemting to fetch non existing book")
	}
}
