// Package interactor contains application-specific business logic.
// Coordinates between [domain] entities and repository interfaces to perform CRUD operations.
// Right now the interactor looks like a pass-through because this is a minimal CRUD example.
package interactor

import (
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/text/language"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/infrastructure/db"
)

// BookInteractor handles business logic.
type BookInteractor struct {
	repo   domain.BookRepository
	logger *slog.Logger
}

// NewBookInteractor creates a new BookInteractor.
func NewBookInteractor(logger *slog.Logger, repo domain.BookRepository) *BookInteractor {
	return &BookInteractor{repo: repo, logger: logger}
}

// CreateBook sends the book to be created to the underlying repository.
// Sets the language tag to english.
func (bi *BookInteractor) CreateBook(book *domain.Book) error {
	// in the real world, CreateBook might, for example,
	// trigger inventory updates, sends events or validates stock.
	book.LanguageTag = language.English.String()
	return bi.repo.Create(book)
}

// GetBook retrieves a domain.Book by its ID.
// Validates the given id is a valid UUID.
func (bi *BookInteractor) GetBook(id string) (*domain.Book, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrInvalidBookID
	}
	b, err := bi.repo.ReadByID(uid)
	if db.IsNotFoundError(err) {
		return nil, domain.ErrBookNotFound
	}
	return b, err
}

// ListBooks retrieves a list of books.
// Does not fail if nothing is found.
func (bi *BookInteractor) ListBooks() ([]*domain.Book, error) {
	return bi.repo.ReadAll()
}

// UpdateBook updates a single book by its ID.
func (bi *BookInteractor) UpdateBook(book *domain.Book) error {
	err := bi.repo.Update(book)
	if db.IsNotFoundError(err) {
		return domain.ErrBookNotFound
	}
	return err
}

// DeleteBook removes a book from the repository.
// Does not fail if nothing is found.
func (bi *BookInteractor) DeleteBook(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrInvalidBookID
	}
	err = bi.repo.Delete(uid)
	if db.IsNotFoundError(err) {
		return nil
	}
	return err
}
