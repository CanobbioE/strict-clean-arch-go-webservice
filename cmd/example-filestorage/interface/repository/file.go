package repository

import (
	"context"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/domain"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/infrastructure/provider"
)

type FileRepository interface {
	Store(context.Context, *domain.File) (*domain.File, error)
	Get(context.Context, *domain.File) (*domain.File, error)
	Delete(context.Context, *domain.File) (*domain.File, error)
}

type Files struct {
	provider provider.FileProvider
}

func NewFilesRepository(fp provider.FileProvider) FileRepository {
	return &Files{provider: fp}
}

func (fr *Files) Store(ctx context.Context, f *domain.File) (*domain.File, error) {

	created, err := fr.provider.Create(ctx, fr.toProvider(f))
	if err != nil {
		return nil, err
	}

	return fr.fromProvider(created), nil
}
func (fr *Files) Get(ctx context.Context, f *domain.File) (*domain.File, error) {
	match, err := fr.provider.Get(ctx, fr.toProvider(f))
	if err != nil {
		return nil, err
	}

	return fr.fromProvider(match), nil
}

func (fr *Files) Delete(ctx context.Context, f *domain.File) (*domain.File, error) {
	created, err := fr.provider.Delete(ctx, fr.toProvider(f))
	if err != nil {
		return nil, err
	}

	return fr.fromProvider(created), nil
}

func (fr *Files) toProvider(f *domain.File) *provider.File {
	return &provider.File{
		Id:      f.Id,
		Name:    f.Name,
		Content: f.Content,
		Owner:   f.Owner,
		Created: f.Created,
	}
}

func (fr *Files) fromProvider(f *provider.File) *domain.File {
	return &domain.File{
		Id:      f.Id,
		Name:    f.Name,
		Content: f.Content,
		Owner:   f.Owner,
		Created: f.Created,
	}
}
