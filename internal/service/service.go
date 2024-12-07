package service

import (
	"context"

	"github.com/SemenShakhray/list-of-song/internal/models"
	"github.com/SemenShakhray/list-of-song/internal/storage"
)

type Service struct {
	storage storage.Storer
}

type Servicer interface {
	GetAll(ctx context.Context, filtres models.Filters) ([]models.Song, error)
	AddSong(ctx context.Context, song models.Song) error
	Update(ctx context.Context, song models.Song) error
	Delete(ctx context.Context, id int) error
	GetText(ctx context.Context, filtres models.Filters, id int) (string, error)
}

func NewService(store storage.Storer) Servicer {
	return &Service{
		storage: store,
	}
}

func (s *Service) AddSong(ctx context.Context, song models.Song) error {
	return s.storage.AddSong(ctx, song)
}

func (s *Service) GetAll(ctx context.Context, filters models.Filters) ([]models.Song, error) {
	return s.storage.GetAll(ctx, filters)
}

func (s *Service) Update(ctx context.Context, song models.Song) error {
	return s.storage.Update(ctx, song)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.storage.Delete(ctx, id)
}

func (s *Service) GetText(ctx context.Context, filters models.Filters, id int) (string, error) {
	return s.storage.GetText(ctx, filters, id)
}
