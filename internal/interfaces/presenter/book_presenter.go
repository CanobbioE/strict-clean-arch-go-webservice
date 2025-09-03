// Package presenter formats output representation for the external world (JSON payload shape, view models).
// Keeps controllers thin.
package presenter

import (
	"log/slog"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
)

// BookPresenter prepares a domain.Book to be returned to an http interface.
type BookPresenter struct {
	logger *slog.Logger
}

// NewBookPresenter creates a new instance of BookPresenter.
func NewBookPresenter(logger *slog.Logger) *BookPresenter {
	return &BookPresenter{logger: logger}
}

// Present returns the map representation of a domain.Book.
// Present prefixes the book ID with the resource type (book:) and
// transform the title to Title Case based on the book language.
func (p *BookPresenter) Present(book *domain.Book) map[string]any {
	return map[string]any{
		"id":     "book:" + book.ID.String(),
		"title":  p.title(book),
		"author": book.Author,
		"price":  book.Price,
	}
}

func (p *BookPresenter) title(book *domain.Book) string {
	lt, err := language.Parse(book.LanguageTag)
	if err != nil {
		p.logger.
			With("book_id", book.ID, "book_language", book.LanguageTag).
			Warn("failed to parse tag language")
		return book.Title
	}

	return cases.Title(lt).String(book.Title)
}
