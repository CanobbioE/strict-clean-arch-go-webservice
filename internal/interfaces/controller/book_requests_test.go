package controller_test

import (
	"testing"

	"github.com/google/uuid"

	"github.com/CanobbioE/strict-clean-arch-go-webservice/internal/interfaces/controller"
)

func TestCreateBookRequest_Validate(t *testing.T) {
	type fields struct {
		Title  string
		Author string
		Price  int
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		compareErr func(error) bool
	}{
		{
			name: "fails if missing title",
			fields: fields{
				Author: "John Doe",
				Price:  42,
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return err.Error() == "title is required"
			},
		},
		{
			name: "fails if missing Author",
			fields: fields{
				Title: "Test Book",
				Price: 42,
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return err.Error() == "author is required"
			},
		},
		{
			name: "fails if invalid price",
			fields: fields{
				Title:  "Test Book",
				Author: "John Doe",
				Price:  -100,
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return err.Error() == "price must be greater than zero"
			},
		},
		{
			name: "succeeds",
			fields: fields{
				Title:  "Test Book",
				Author: "John Doe",
				Price:  42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &controller.CreateBookRequest{
				Title:  tt.fields.Title,
				Author: tt.fields.Author,
				Price:  tt.fields.Price,
			}
			err := r.Validate()
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateBookRequest_Validate(t *testing.T) {
	type fields struct {
		ID    string
		Price int
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		compareErr func(error) bool
	}{
		{
			name: "fails if missing id",
			fields: fields{
				Price: 42,
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return err.Error() == "id is required"
			},
		},
		{
			name: "fails if invalid price",
			fields: fields{
				ID:    uuid.NewString(),
				Price: 0,
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return err.Error() == "price must be greater than zero"
			},
		},
		{
			name: "fails if invalid id",
			fields: fields{
				ID:    "book-id",
				Price: 42,
			},
			wantErr: true,
			compareErr: func(err error) bool {
				return err.Error() == "invalid id"
			},
		},
		{
			name: "succeeds",
			fields: fields{
				ID:    uuid.NewString(),
				Price: 42,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &controller.UpdateBookRequest{
				ID:    tt.fields.ID,
				Price: tt.fields.Price,
			}
			err := r.Validate()
			if (tt.wantErr != (err != nil)) || (tt.wantErr && !tt.compareErr(err)) {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
