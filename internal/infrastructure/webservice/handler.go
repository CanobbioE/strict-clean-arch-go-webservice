// Package webservice defines HTTP routes and maps them to [controller] methods.
// Abstracts away the web server framework.
package webservice

import (
	"net/http"
)

// BookController is the interface abstraction of an HTTP controller.
type BookController interface {
	// CreateBook handles create requests over http.
	CreateBook(w http.ResponseWriter, r *http.Request)
	// GetBook handles read book by ID requests over http.
	GetBook(w http.ResponseWriter, r *http.Request)
	// ListBooks handles read books requests over http.
	ListBooks(w http.ResponseWriter, _ *http.Request)
	// UpdateBook handles update requests over http.
	UpdateBook(w http.ResponseWriter, r *http.Request)
	// DeleteBook handles delete book by ID requests over http.
	DeleteBook(w http.ResponseWriter, r *http.Request)
}

// NewHandler creates a new webservice serving CRUD operation on the /books endpoint.
func NewHandler(bc BookController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /v1/books", bc.CreateBook)
	mux.HandleFunc("GET /v1/books/{id}", bc.GetBook)
	mux.HandleFunc("GET /v1/books", bc.ListBooks)
	mux.HandleFunc("PATCH /v1/books", bc.UpdateBook)
	mux.HandleFunc("DELETE /v1/books/{id}", bc.DeleteBook)
	return mux
}
