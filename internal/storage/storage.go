package storage

import (
	"context"

	"github.com/SemenShakhray/list-of-song/internal/models"
)

type Storer interface {
	GetAll(ctx context.Context, filtres models.Filters) ([]models.Song, error)
	AddSong(ctx context.Context, song models.Song) error
	Update(ctx context.Context, song models.Song) error
	Delete(ctx context.Context, id int) error
	GetText(ctx context.Context, filtres models.Filters, id int) (string, error)
}
