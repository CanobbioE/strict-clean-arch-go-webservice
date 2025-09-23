package webservice_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/infrastructure/webservice"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/test/mocks"
)

func TestNewHandler(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockBooksController := mocks.NewMockBookController(mockCtl)
	handler := webservice.NewHandler(mockBooksController)
	server := httptest.NewServer(handler)

	tests := []struct {
		name             string
		method           string
		endpoint         string
		mockExpectations func()
	}{
		{
			name:     "PUT /v1/books",
			method:   http.MethodPut,
			endpoint: "/v1/books",
			mockExpectations: func() {
				mockBooksController.EXPECT().CreateBook(gomock.Any(), gomock.Any())
			},
		},
		{
			name:     "GET /v1/books/{id}",
			method:   http.MethodGet,
			endpoint: "/v1/books/book-id",
			mockExpectations: func() {
				mockBooksController.EXPECT().GetBook(gomock.Any(), gomock.Any())
			},
		},
		{
			name:     "GET /v1/books",
			method:   http.MethodGet,
			endpoint: "/v1/books",
			mockExpectations: func() {
				mockBooksController.EXPECT().ListBooks(gomock.Any(), gomock.Any())
			},
		},
		{
			name:     "PATCH /v1/books",
			method:   http.MethodPatch,
			endpoint: "/v1/books",
			mockExpectations: func() {
				mockBooksController.EXPECT().UpdateBook(gomock.Any(), gomock.Any())
			},
		},
		{
			name:     "DELETE /v1/books/{id}",
			method:   http.MethodDelete,
			endpoint: "/v1/books/book-id",
			mockExpectations: func() {
				mockBooksController.EXPECT().DeleteBook(gomock.Any(), gomock.Any())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExpectations()
			req, err := http.NewRequestWithContext(context.Background(), tt.method, server.URL+tt.endpoint, http.NoBody)
			if err != nil {
				t.Fatal("failed to create request:", err)
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal("failed to make request:", err)
			}
			defer func() {
				err = res.Body.Close()
				if err != nil {
					t.Fatal("failed to close response body:", err)
				}
			}()

			if res.StatusCode != http.StatusOK {
				t.Errorf("status code is not 200: %d", res.StatusCode)
			}
		})
	}
}
