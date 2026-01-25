package metadata

import (
	"context"
	"errors"

	"movieexamplekhubaib.com/metadata/internal/repository"
	"movieexamplekhubaib.com/metadata/pkg/model"
)

// Controller meant to encapsulate the repository logic

var ErrNotFound = errors.New("Record not found")

type metadataRepository interface {
	Get (ctx context.Context, id string) (*model.Metadata, error)
}

type Controller struct {
	repo metadataRepository
}

func New (repo metadataRepository) *Controller {
	return &Controller{repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	res, err := c.repo.Get(ctx, id)
	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return nil, ErrNotFound
	}
	return res, err
}
