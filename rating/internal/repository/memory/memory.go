package memory

import (
	"context"

	"movieexamplekhubaib.com/rating/internal/repository"
	"movieexamplekhubaib.com/rating/pkg/model"
)

// Object holding data
type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func New() *Repository {
	return &Repository{map[model.RecordType]map[model.RecordID][]model.Rating{}}
}

// Pointer receiver most likely to actually get values
func (r *Repository) Get(ctx context.Context, recordId model.RecordID, recordType model.RecordType) ([]model.Rating, error) {
	// Check if we have records for the type
	if _, ok := r.data[recordType]; !ok {
		return nil, repository.ErrNotFound
	}
	// Check in the type if we have that id
	if ratings, ok := r.data[recordType][recordId]; !ok || len(ratings) == 0 {
		return nil, repository.ErrNotFound
	}

	return r.data[recordType][recordId], nil
}

// Insert a rating for a record. Pointer receiver used to change the repository
func (r *Repository) Put(ctx context.Context, recordId model.RecordID, recordType model.RecordType, rating *model.Rating) error {
	return nil
}