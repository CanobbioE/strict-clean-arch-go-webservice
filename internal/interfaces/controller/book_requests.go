package controller

import (
	"errors"

	"github.com/google/uuid"
)

// CreateBookRequest defines the expected request for create.
type CreateBookRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Price  int    `json:"price"`
}

// Validate a CreateBookRequest.
func (r *CreateBookRequest) Validate() error {
	switch {
	case r.Title == "":
		return errors.New("title is required")
	case r.Author == "":
		return errors.New("author is required")
	case r.Price <= 0:
		return errors.New("price must be greater than zero")
	}
	return nil
}

// UpdateBookRequest defines the expected request for update.
type UpdateBookRequest struct {
	ID    string `json:"id"`
	Price int    `json:"price"`
}

// Validate an UpdateBookRequest.
func (r *UpdateBookRequest) Validate() error {
	switch {
	case r.ID == "":
		return errors.New("id is required")
	case r.Price <= 0:
		return errors.New("price must be greater than zero")
	}

	if _, err := uuid.Parse(r.ID); err != nil {
		return errors.New("invalid id")
	}
	return nil
}
