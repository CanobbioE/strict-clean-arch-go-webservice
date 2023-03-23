package controller

import (
	"context"
	"encoding/json"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/domain"
	"io"
	"io/ioutil"
)

type FileController interface {
	ParseRequestBody(context.Context, io.ReadCloser) (*domain.File, error)
	ValidateRequest(context.Context, *domain.File) error
}

type File struct{}

func (fc *File) ParseRequestBody(ctx context.Context, rc io.ReadCloser) (*domain.File, error) {
	var req domain.File

	body, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(body, &req); err != nil {
		return nil, err
	}

	return &req, nil
}

func (fc *File) ValidateRequest(context.Context, *domain.File) error {
	return nil
}
