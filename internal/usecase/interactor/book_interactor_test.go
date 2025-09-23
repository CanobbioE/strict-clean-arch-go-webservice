package interactor_test

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
	"golang.org/x/text/language"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/domain"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/test/mocks"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/test/testlog"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/usecase/interactor"
)

func TestBookInteractor_CreateBook(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockBookRepository := mocks.NewMockBookRepository(mockCtl)
	logger := testlog.NewTestLogger()
	book := &domain.Book{Title: "A book"}

	tests := []struct {
		name             string
		wantErr          bool
		compareErr       func(error) bool
		mockExpectations func()
	}{
		{
			name: "fails",
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Create(&domain.Book{
						Title:       book.Title,
						LanguageTag: language.English.String(),
					}).
					Return(errors.New("oops"))
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return strings.Contains(err.Error(), "oops")
			},
		},
		{
			name: "succeeds",
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Create(&domain.Book{
						Title:       book.Title,
						LanguageTag: language.English.String(),
					}).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}
			bi := interactor.NewBookInteractor(logger, mockBookRepository)
			err := bi.CreateBook(book)
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("CreateBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBookInteractor_GetBook(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockBookRepository := mocks.NewMockBookRepository(mockCtl)
	logger := testlog.NewTestLogger()
	book := &domain.Book{Title: "A book", ID: uuid.New()}

	tests := []struct {
		name             string
		id               string
		want             *domain.Book
		wantErr          bool
		compareErr       func(error) bool
		mockExpectations func()
	}{
		{
			name:    "fails to parse id",
			id:      "invalid",
			wantErr: true,
			compareErr: func(err error) bool {
				return errors.Is(err, domain.ErrInvalidBookID)
			},
		},
		{
			name: "fails with book not found error",
			id:   book.ID.String(),
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					ReadByID(book.ID).
					Return(nil, errors.New("book not found"))
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return errors.Is(err, domain.ErrBookNotFound)
			},
		},
		{
			name: "fails with generic error",
			id:   book.ID.String(),
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					ReadByID(book.ID).
					Return(nil, errors.New("something broke"))
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return strings.Contains(err.Error(), "something broke")
			},
		},
		{
			name: "succeeds",
			id:   book.ID.String(),
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					ReadByID(book.ID).
					Return(book, nil)
			},
			want: book,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}
			bi := interactor.NewBookInteractor(logger, mockBookRepository)
			got, err := bi.GetBook(tt.id)
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("GetBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBook() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookInteractor_ListBooks(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockBookRepository := mocks.NewMockBookRepository(mockCtl)
	logger := testlog.NewTestLogger()
	books := []*domain.Book{{Title: "A book", ID: uuid.New()}}

	tests := []struct {
		name             string
		want             []*domain.Book
		wantErr          bool
		compareErr       func(error) bool
		mockExpectations func()
	}{
		{
			name: "fails",
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					ReadAll().
					Return(nil, errors.New("oops"))
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return strings.Contains(err.Error(), "oops")
			},
		},
		{
			name: "succeeds",
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					ReadAll().
					Return(books, nil)
			},
			want: books,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}
			bi := interactor.NewBookInteractor(logger, mockBookRepository)
			got, err := bi.ListBooks()
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("GetBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListBooks() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBookInteractor_UpdateBook(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockBookRepository := mocks.NewMockBookRepository(mockCtl)
	logger := testlog.NewTestLogger()
	book := &domain.Book{Title: "A book", ID: uuid.New()}

	tests := []struct {
		name             string
		book             *domain.Book
		wantErr          bool
		compareErr       func(error) bool
		mockExpectations func()
	}{
		{
			name: "fails with book not found error",
			book: book,
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Update(book).
					Return(errors.New("book not found"))
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return errors.Is(err, domain.ErrBookNotFound)
			},
		},
		{
			name: "fails with generic error",
			book: book,
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Update(book).
					Return(errors.New("something broke"))
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return strings.Contains(err.Error(), "something broke")
			},
		},
		{
			name: "succeeds",
			book: book,
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Update(book).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}
			bi := interactor.NewBookInteractor(logger, mockBookRepository)
			err := bi.UpdateBook(tt.book)
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("UpdateBook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestBookInteractor_DeleteBook(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockBookRepository := mocks.NewMockBookRepository(mockCtl)
	logger := testlog.NewTestLogger()
	book := &domain.Book{Title: "A book", ID: uuid.New()}

	tests := []struct {
		name             string
		id               string
		wantErr          bool
		compareErr       func(error) bool
		mockExpectations func()
	}{
		{
			name:    "fails to parse id",
			id:      "invalid",
			wantErr: true,
			compareErr: func(err error) bool {
				return errors.Is(err, domain.ErrInvalidBookID)
			},
		},
		{
			name: "succeeds with book not found error",
			id:   book.ID.String(),
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Delete(book.ID).
					Return(errors.New("book not found"))
			},
		},
		{
			name: "fails with generic error",
			id:   book.ID.String(),
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Delete(book.ID).
					Return(errors.New("something broke"))
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return strings.Contains(err.Error(), "something broke")
			},
		},
		{
			name: "succeeds",
			id:   book.ID.String(),
			mockExpectations: func() {
				mockBookRepository.EXPECT().
					Delete(book.ID).
					Return(nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockExpectations != nil {
				tt.mockExpectations()
			}
			bi := interactor.NewBookInteractor(logger, mockBookRepository)
			err := bi.DeleteBook(tt.id)
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
