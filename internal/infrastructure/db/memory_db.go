// Package db provides concrete repository implementation (in-memory or DB).
// Bridges actual persistence technology and the [domain] interface.
package db

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/google/uuid"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
)

// InMemoryBookRepo implements domain.BookRepository as an in-memory database.
// The repository is wiped with each restart.
type InMemoryBookRepo struct {
	books  map[uuid.UUID]*domain.Book
	logger *slog.Logger
}

// NewInMemoryBookRepo creates a new instance of InMemoryBookRepo, implementing domain.BookRepository.
func NewInMemoryBookRepo(logger *slog.Logger) *InMemoryBookRepo {
	return &InMemoryBookRepo{books: make(map[uuid.UUID]*domain.Book), logger: logger}
}

// Create a new book entry.
func (r *InMemoryBookRepo) Create(book *domain.Book) error {
	book.ID = uuid.New()
	r.books[book.ID] = book
	return nil
}

// ReadByID return a single book that matches the given ID.
func (r *InMemoryBookRepo) ReadByID(id uuid.UUID) (*domain.Book, error) {
	if b, ok := r.books[id]; ok {
		return b, nil
	}
	return nil, errors.New("book not found")

}

// ReadAll return a list of books.
func (r *InMemoryBookRepo) ReadAll() ([]*domain.Book, error) {
	var list []*domain.Book
	for _, b := range r.books {
		list = append(list, b)
	}
	return list, nil
}

// Update the price of a book.
func (r *InMemoryBookRepo) Update(book *domain.Book) error {
	if _, ok := r.books[book.ID]; !ok {
		return errors.New("book not found")
	}
	r.books[book.ID].Price = book.Price
	return nil
}

// Delete a single book, matched by ID.
func (r *InMemoryBookRepo) Delete(id uuid.UUID) error {
	if _, ok := r.books[id]; !ok {
		return errors.New("book not found")
	}
	delete(r.books, id)
	return nil
}

// IsNotFoundError return true if the error is not nil and is a not found error.
// This is more useful with real DBs, where errors are a bit more cryptic
// (e.g. [-106] Row to DELETE not found).
func IsNotFoundError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "book not found")
}
