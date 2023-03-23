package interactor

import (
	"context"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/domain"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/interface/repository"
	"time"
)

type FileInteractor interface {
	Upload(context.Context, *domain.File) (*domain.File, error)
	Retrieve(context.Context, string) (*domain.File, error)
	Delete(context.Context, string) (*domain.File, error)
}

type File struct {
	Repo repository.FileRepository
}

func (fi *File) Upload(ctx context.Context, f *domain.File) (*domain.File, error) {
	now := time.Now()
	f.Created = &now
	return fi.Repo.Store(ctx, f)
}

func (fi *File) Retrieve(ctx context.Context, id string) (*domain.File, error) {
	return fi.Repo.Get(ctx, &domain.File{Id: &id})
}

func (fi *File) Delete(ctx context.Context, id string) (*domain.File, error) {
	return fi.Repo.Delete(ctx, &domain.File{Id: &id})
}
