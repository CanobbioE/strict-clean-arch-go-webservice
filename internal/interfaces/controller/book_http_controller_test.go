package controller_test

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/interfaces/controller"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/interfaces/presenter"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/test/mocks"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/test/testlog"
)

type controllerFields struct {
	interactor    controller.BookInteractor
	bookPresenter controller.BookPresenter
	errPresenter  controller.ErrorPresenter
	logger        *slog.Logger
}

func TestBookController_CreateBook(t *testing.T) {
	logger := testlog.NewTestLogger()
	mockCtl := gomock.NewController(t)
	mockBookInteractor := mocks.NewMockBookInteractor(mockCtl)
	commonFields := controllerFields{
		interactor:    mockBookInteractor,
		bookPresenter: presenter.NewBookPresenter(logger),
		errPresenter:  presenter.NewErrorPresenter(logger),
		logger:        logger,
	}

	tests := []struct {
		name             string
		body             string
		mockExpectations func()
		expect           func(*httptest.ResponseRecorder)
	}{
		{
			name: "fails to decode request body",
			body: `}} broken json body`,
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusBadRequest {
					t.Errorf("want status: %d, got status %d", http.StatusBadRequest, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"invalid character '}' looking for beginning of value","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "fails to validate request",
			body: `
				{
					"title": "a book",
					"author": "someone",
					"price": -10
				}`,
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusBadRequest {
					t.Errorf("want status: %d, got status %d", http.StatusBadRequest, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"price must be greater than zero","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "fails to create new book",
			body: `
				{
					"title": "a book",
					"author": "someone",
					"price": 42
				}`,
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					CreateBook(&domain.Book{
						Title:  "a book",
						Author: "someone",
						Price:  42,
					}).
					Return(errors.New("oops"))
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusInternalServerError {
					t.Errorf("want status: %d, got status %d", http.StatusInternalServerError, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"internal server error","status":"Internal Server Error"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "succeeds",
			body: `
				{
					"title": "a book",
					"author": "someone",
					"price": 42
				}`,
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					CreateBook(&domain.Book{
						Title:  "a book",
						Author: "someone",
						Price:  42,
					})
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusCreated {
					t.Errorf("want status: %d, got status %d", http.StatusCreated, res.Code)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}

			bc := controller.NewBookController(
				commonFields.logger,
				commonFields.interactor,
				commonFields.bookPresenter,
				commonFields.errPresenter,
			)
			r := httptest.NewRequest(http.MethodPut, "/v1/books", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			bc.CreateBook(w, r)
			if tt.expect != nil {
				tt.expect(w)
			}
		})
	}
}

func TestBookController_GetBook(t *testing.T) {
	logger := testlog.NewTestLogger()
	mockCtl := gomock.NewController(t)
	mockBookInteractor := mocks.NewMockBookInteractor(mockCtl)
	bookID := uuid.New()
	commonFields := controllerFields{
		interactor:    mockBookInteractor,
		bookPresenter: presenter.NewBookPresenter(logger),
		errPresenter:  presenter.NewErrorPresenter(logger),
		logger:        logger,
	}

	tests := []struct {
		name             string
		id               string
		mockExpectations func()
		expect           func(*httptest.ResponseRecorder)
	}{
		{
			name: "fails if missing ID",
			id:   "",
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusBadRequest {
					t.Errorf("want status: %d, got status %d", http.StatusBadRequest, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"book id is required","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "fails to get book",
			id:   bookID.String(),
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					GetBook(bookID.String()).
					Return(nil, domain.ErrBookNotFound)
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusNotFound {
					t.Errorf("want status: %d, got status %d", http.StatusNotFound, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"book not found","status":"Not Found"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "succeeds",
			id:   bookID.String(),
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					GetBook(bookID.String()).
					Return(&domain.Book{
						ID:          bookID,
						Title:       "a book",
						Author:      "someone",
						Price:       42,
						LanguageTag: language.Italian.String(),
					}, nil)
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusOK {
					t.Errorf("want status: %d, got status %d", http.StatusCreated, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"author":"someone","id":"book:` + bookID.String() + `","price":42,"title":"A Book"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}

			bc := controller.NewBookController(
				commonFields.logger,
				commonFields.interactor,
				commonFields.bookPresenter,
				commonFields.errPresenter,
			)
			r := httptest.NewRequest(http.MethodGet, "/v1/books/"+tt.id, http.NoBody)
			r.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			bc.GetBook(w, r)
			if tt.expect != nil {
				tt.expect(w)
			}
		})
	}
}

func TestBookController_ListBooks(t *testing.T) {
	logger := testlog.NewTestLogger()
	mockCtl := gomock.NewController(t)
	mockBookInteractor := mocks.NewMockBookInteractor(mockCtl)
	bookID := uuid.New()
	commonFields := controllerFields{
		interactor:    mockBookInteractor,
		bookPresenter: presenter.NewBookPresenter(logger),
		errPresenter:  presenter.NewErrorPresenter(logger),
		logger:        logger,
	}

	tests := []struct {
		name             string
		id               string
		mockExpectations func()
		expect           func(*httptest.ResponseRecorder)
	}{
		{
			name: "fails to list book",
			id:   bookID.String(),
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					ListBooks().
					Return(nil, errors.New("oops"))
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusInternalServerError {
					t.Errorf("want status: %d, got status %d", http.StatusInternalServerError, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"internal server error","status":"Internal Server Error"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "succeeds",
			id:   bookID.String(),
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					ListBooks().
					Return([]*domain.Book{
						{
							ID:     bookID,
							Title:  "a book",
							Author: "someone",
							Price:  42,
						},
					}, nil)
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusOK {
					t.Errorf("want status: %d, got status %d", http.StatusCreated, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `[{"author":"someone","id":"book:` + bookID.String() + `","price":42,"title":"a book"}]`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}

			bc := controller.NewBookController(
				commonFields.logger,
				commonFields.interactor,
				commonFields.bookPresenter,
				commonFields.errPresenter,
			)
			r := httptest.NewRequest(http.MethodGet, "/v1/books", http.NoBody)
			w := httptest.NewRecorder()

			bc.ListBooks(w, r)
			if tt.expect != nil {
				tt.expect(w)
			}
		})
	}
}

func TestBookController_UpdateBook(t *testing.T) {
	logger := testlog.NewTestLogger()
	mockCtl := gomock.NewController(t)
	mockBookInteractor := mocks.NewMockBookInteractor(mockCtl)
	commonFields := controllerFields{
		interactor:    mockBookInteractor,
		bookPresenter: presenter.NewBookPresenter(logger),
		errPresenter:  presenter.NewErrorPresenter(logger),
		logger:        logger,
	}
	bookID := uuid.New()

	tests := []struct {
		name             string
		body             string
		mockExpectations func()
		expect           func(*httptest.ResponseRecorder)
	}{
		{
			name: "fails to decode request body",
			body: `}} broken json body`,
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusBadRequest {
					t.Errorf("want status: %d, got status %d", http.StatusBadRequest, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"invalid character '}' looking for beginning of value","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "fails to validate request",
			body: `{
					"id": "",
					"price": 10
				}`,
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusBadRequest {
					t.Errorf("want status: %d, got status %d", http.StatusBadRequest, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"id is required","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "fails to update book",
			body: `{
					"id": " ` + bookID.String() + ` ",
					"price": 10
				}`,
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					UpdateBook(&domain.Book{
						ID:    bookID,
						Price: 10,
					}).
					Return(errors.New("oops"))
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusInternalServerError {
					t.Errorf("want status: %d, got status %d", http.StatusInternalServerError, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"internal server error","status":"Internal Server Error"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "succeeds",
			body: `
				{
					"id": " ` + bookID.String() + ` ",
					"price": 10
				}`,
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					UpdateBook(&domain.Book{
						ID:    bookID,
						Price: 10,
					})
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusAccepted {
					t.Errorf("want status: %d, got status %d", http.StatusAccepted, res.Code)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}

			bc := controller.NewBookController(
				commonFields.logger,
				commonFields.interactor,
				commonFields.bookPresenter,
				commonFields.errPresenter,
			)
			r := httptest.NewRequest(http.MethodPatch, "/v1/books", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			bc.UpdateBook(w, r)
			if tt.expect != nil {
				tt.expect(w)
			}
		})
	}
}

func TestBookController_DeleteBook(t *testing.T) {
	logger := testlog.NewTestLogger()
	mockCtl := gomock.NewController(t)
	mockBookInteractor := mocks.NewMockBookInteractor(mockCtl)
	bookID := uuid.New()
	commonFields := controllerFields{
		interactor:    mockBookInteractor,
		bookPresenter: presenter.NewBookPresenter(logger),
		errPresenter:  presenter.NewErrorPresenter(logger),
		logger:        logger,
	}

	tests := []struct {
		name             string
		id               string
		mockExpectations func()
		expect           func(*httptest.ResponseRecorder)
	}{
		{
			name: "fails if missing ID",
			id:   "",
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusBadRequest {
					t.Errorf("want status: %d, got status %d", http.StatusBadRequest, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"book id is required","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "fails to delete book",
			id:   bookID.String(),
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					DeleteBook(bookID.String()).
					Return(domain.ErrInvalidBookID)
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusBadRequest {
					t.Errorf("want status: %d, got status %d", http.StatusBadRequest, res.Code)
				}
				got := strings.TrimSpace(res.Body.String())
				want := `{"message":"invalid book id","status":"Bad Request"}`
				if got != want {
					t.Errorf("want %s, got %s", want, got)
				}
			},
		},
		{
			name: "succeeds",
			id:   bookID.String(),
			mockExpectations: func() {
				mockBookInteractor.EXPECT().
					DeleteBook(bookID.String()).
					Return(nil)
			},
			expect: func(res *httptest.ResponseRecorder) {
				if res.Code != http.StatusOK {
					t.Errorf("want status: %d, got status %d", http.StatusCreated, res.Code)
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}

			bc := controller.NewBookController(
				commonFields.logger,
				commonFields.interactor,
				commonFields.bookPresenter,
				commonFields.errPresenter,
			)
			r := httptest.NewRequest(http.MethodDelete, "/v1/books/"+tt.id, http.NoBody)
			r.SetPathValue("id", tt.id)
			w := httptest.NewRecorder()

			bc.DeleteBook(w, r)
			if tt.expect != nil {
				tt.expect(w)
			}
		})
	}
}
