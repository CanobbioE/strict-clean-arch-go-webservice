// Package db provides concrete repository implementation (in-memory or DB).
// Bridges actual persistence technology and the [domain] interface.
package db

import (
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
)

var errNotFound = errors.New("book not found")

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
	return nil, errNotFound
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
		return errNotFound
	}
	r.books[book.ID].Price = book.Price
	return nil
}

// Delete a single book, matched by ID.
func (r *InMemoryBookRepo) Delete(id uuid.UUID) error {
	if _, ok := r.books[id]; !ok {
		return errNotFound
	}
	delete(r.books, id)
	return nil
}

// IsNotFoundError return true if the error is not nil and matches the internal errNotFound.
func IsNotFoundError(err error) bool {
	return errors.Is(err, errNotFound)
}
