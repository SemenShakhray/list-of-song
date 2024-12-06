package storage

import (
	"context"

	"listsongs/internal/models"
)

type Storer interface {
	GetAll(ctx context.Context, filtres models.Filters) ([]models.Song, error)
	AddSong(ctx context.Context, song models.Song) error
	Update(ctx context.Context, song models.Song, id int) error
	Delete(ctx context.Context, id int) error
	GetText(ctx context.Context, filtres models.Filters, id int) (string, error)
}
