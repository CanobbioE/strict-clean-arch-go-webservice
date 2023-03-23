package presenter

import (
	"context"
	"fmt"
	"github.com/CanobbioE/strict-clean-arch-go-webservice/cmd/example-filestorage/domain"
)

type FilePresenter interface {
	AnonymizeOwner(context.Context, *domain.File) (*domain.File, error)
	RemoveID(context.Context, *domain.File) (*domain.File, error)
}

type File struct{}

func (fp *File) AnonymizeOwner(ctx context.Context, f *domain.File) (*domain.File, error) {
	var replacer string

	if f.Owner == nil {
		return f, nil
	}

	for _, c := range *f.Owner {
		if c == ' ' {
			replacer = fmt.Sprintf("%s ", replacer)
			continue
		}
		replacer = fmt.Sprintf("%sx", replacer)
	}

	f.Owner = &replacer

	return f, nil
}

func (fp *File) RemoveID(ctx context.Context, f *domain.File) (*domain.File, error) {
	f.Id = nil
	return f, nil
}
