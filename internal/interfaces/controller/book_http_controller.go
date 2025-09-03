// Package controller handles HTTP request/response flow.
// Converts external input (JSON, URL params) into domain objects and calls use cases.
// Returns results or errors to clients.
package controller

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
)

// BookInteractor is the interface an interactor must implement
// to be used by the BookController to execute business logic.
type BookInteractor interface {
	// CreateBook sends the book to be created to the underlying repository.
	CreateBook(book *domain.Book) error
	// GetBook retrieves a domain.Book by its ID.
	GetBook(id string) (*domain.Book, error)
	// ListBooks retrieves a list of books.
	ListBooks() ([]*domain.Book, error)
	// UpdateBook updates a single book by its ID.
	UpdateBook(book *domain.Book) error
	// DeleteBook removes a book from the repository.
	DeleteBook(id string) error
}

// BookPresenter is the interface a presenter must implement
// to be used by the BookController to return successful responses.
type BookPresenter interface {
	// Present prepares the domain.Book message to be returned.
	Present(book *domain.Book) map[string]any
}

// ErrorPresenter is the interface a presenter must implement
// to be used by the BookController to return error responses.
type ErrorPresenter interface {
	// Present prepares the error message to be returned through w.
	Present(w http.ResponseWriter, err error)
}

// BookController handles http requests, validates them and transform them into domain objects.
// The domain objects are then passed to the usecase layer, executing the business logic.
type BookController struct {
	interactor    BookInteractor
	bookPresenter BookPresenter
	errPresenter  ErrorPresenter
	logger        *slog.Logger
}

// NewBookController creates a new instance of BookController.
func NewBookController(
	logger *slog.Logger,
	interactor BookInteractor,
	bookPresenter BookPresenter,
	errPresenter ErrorPresenter,
) *BookController {
	return &BookController{
		interactor:    interactor,
		logger:        logger,
		bookPresenter: bookPresenter,
		errPresenter:  errPresenter,
	}
}

// CreateBook handles CreateBookRequest over http.
func (bc *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	var b CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		bc.logger.With("error", err).Error("unable to decode request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := b.validate(); err != nil {
		bc.logger.With("error", err).Error("invalid request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := bc.interactor.CreateBook(&domain.Book{
		Title:  b.Title,
		Author: b.Author,
		Price:  b.Price,
	}); err != nil {
		bc.logger.With("error", err).Error("unable to create book")
		bc.errPresenter.Present(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetBook handles read book by ID requests over http.
func (bc *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	l := bc.logger.With("book_id", id)

	if id == "" {
		http.Error(w, "book id required", http.StatusBadRequest)
		return
	}

	book, err := bc.interactor.GetBook(id)
	if err != nil {
		l.With("error", err).Error("error getting book")
		bc.errPresenter.Present(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(bc.bookPresenter.Present(book))
	if err != nil {
		l.With("error", err).Error("error presenting book")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ListBooks handles read books requests over http.
func (bc *BookController) ListBooks(w http.ResponseWriter, _ *http.Request) {
	books, err := bc.interactor.ListBooks()
	if err != nil {
		bc.logger.With("error", err).Error("error listing books")
		bc.errPresenter.Present(w, err)
		return
	}

	res := make([]map[string]any, len(books))
	for i, book := range books {
		res[i] = bc.bookPresenter.Present(book)
	}

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		bc.logger.With("error", err).Error("error presenting books")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// UpdateBook handles UpdateBookRequest over http.
func (bc *BookController) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var b UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		bc.logger.With("error", err).Error("unable to decode request body")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := b.validate(); err != nil {
		bc.logger.With("error", err).Error("invalid request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := bc.interactor.UpdateBook(&domain.Book{
		ID:    uuid.MustParse(b.ID),
		Price: b.Price,
	})
	if err != nil {
		bc.logger.With("book_id", b.ID).With("error", err).Error("error updating book")
		bc.errPresenter.Present(w, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// DeleteBook handles delete book by ID requests over http.
func (bc *BookController) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	l := bc.logger.With("book_id", id)

	if id == "" {
		http.Error(w, "book id required", http.StatusBadRequest)
		return
	}

	err := bc.interactor.DeleteBook(id)
	if err != nil {
		l.With("error", err).Error("error deleting book")
		bc.errPresenter.Present(w, err)
		return
	}
}
