package db_test

import (
	"reflect"
	"testing"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/infrastructure/db"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/test/testlog"
)

func TestNewInMemoryBookRepo(t *testing.T) {
	logger := testlog.NewTestLogger()
	t.Run("crud operations", func(t *testing.T) {
		repo := db.NewInMemoryBookRepo(logger)
		book := &domain.Book{
			Title:       "A Book",
			Author:      "An Author",
			LanguageTag: "en",
			Price:       10,
		}

		err := repo.Create(book)
		if err != nil {
			t.Fatalf("error creating book: %v", err)
		}

		readResult, err := repo.ReadByID(book.ID)
		if err != nil {
			t.Fatalf("error reading book: %v", err)
		}
		if !reflect.DeepEqual(book, readResult) {
			t.Errorf("expected %+v, got %+v", book, readResult)
		}

		updatedBook := &domain.Book{
			ID:          book.ID,
			Title:       book.Title,
			Author:      book.Author,
			LanguageTag: book.LanguageTag,
			Price:       42,
		}

		err = repo.Update(updatedBook)
		if err != nil {
			t.Fatalf("error updating book: %v", err)
		}

		readAllResult, err := repo.ReadAll()
		if err != nil {
			t.Fatalf("error reading all books: %v", err)
		}
		if !reflect.DeepEqual(updatedBook, readAllResult[0]) {
			t.Errorf("expected %+v, got %+v", updatedBook, readAllResult[0])
		}

		err = repo.Delete(updatedBook.ID)
		if err != nil {
			t.Fatalf("error deleting book: %v", err)
		}

		_, err = repo.ReadByID(updatedBook.ID)
		if !db.IsNotFoundError(err) {
			t.Fatalf("ReadByID() expected not found error, got: %v", err)
		}
		err = repo.Update(updatedBook)
		if !db.IsNotFoundError(err) {
			t.Fatalf("Update() expected not found error, got: %v", err)
		}
		err = repo.Delete(updatedBook.ID)
		if !db.IsNotFoundError(err) {
			t.Fatalf("Delete() expected not found error, got: %v", err)
		}
	})
}
