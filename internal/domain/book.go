// Package domain defines core entities (Book) with only fields and domain-level validation if needed.
// It also defines interfaces for persistence (BookRepository); pure abstraction, no implementation details
package domain

import (
	"github.com/google/uuid"
)

// Book represents a book entity in the system.
type Book struct {
	Title       string
	Author      string
	LanguageTag string
	ID          uuid.UUID
	Price       int
}

// BookRepository defines repository behavior for Book entities.
type BookRepository interface {
	// Create a new book entry.
	Create(book *Book) error
	// ReadByID return a single book that matches the given ID.
	ReadByID(id uuid.UUID) (*Book, error)
	// ReadAll return a list of books.
	ReadAll() ([]*Book, error)
	// Update a book by ID.
	Update(book *Book) error
	// Delete a single book, matched by ID.
	Delete(id uuid.UUID) error
}
