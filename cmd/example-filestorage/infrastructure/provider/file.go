package provider

import (
	"context"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/infrastructure/db"
	"time"
)

type File struct {
	Id      *string
	Name    *string
	Content *string
	Owner   *string
	Created *time.Time
}

func (f *File) ID() string {
	if f.Id != nil {
		return *f.Id
	}
	return ""
}

func (f *File) SetID(id string) {
	f.Id = &id
}

type FileProvider interface {
	Get(context.Context, *File) (*File, error)
	Create(context.Context, *File) (*File, error)
	Delete(context.Context, *File) (*File, error)
}

type DbFileProvider struct {
	dbHandler db.Handler
}

func NewDbFileProvider(handler db.Handler) FileProvider {
	return &DbFileProvider{dbHandler: handler}
}

func (dfp *DbFileProvider) Get(ctx context.Context, f *File) (*File, error) {
	res, err := dfp.dbHandler.First(ctx, f)
	if err != nil {
		return nil, err
	}

	return res.(*File), nil
}

func (dfp *DbFileProvider) Create(ctx context.Context, f *File) (*File, error) {
	err := dfp.dbHandler.Save(ctx, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (dfp *DbFileProvider) Delete(ctx context.Context, f *File) (*File, error) {
	err := dfp.dbHandler.Delete(ctx, f)
	if err != nil {
		return nil, err
	}

	return f, nil
}
