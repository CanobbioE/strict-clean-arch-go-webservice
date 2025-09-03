// Package webservice defines HTTP routes and maps them to [controller] methods.
// Abstracts away the web server framework.
package webservice

import (
	"net/http"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/interfaces/controller"
)

// NewRouter creates a new webservice serving CRUD operation on the /books endpoint.
func NewRouter(bc *controller.BookController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /books", bc.CreateBook)
	mux.HandleFunc("GET /books/{id}", bc.GetBook)
	mux.HandleFunc("GET /books", bc.ListBooks)
	mux.HandleFunc("POST /books", bc.UpdateBook)
	mux.HandleFunc("DELETE /books/{id}", bc.DeleteBook)
	return mux
}
