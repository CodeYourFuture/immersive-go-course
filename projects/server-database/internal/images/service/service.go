package service

import (
	"context"
	"fmt"
	"log"

	"github.com/lib/pq"

	"server-database/internal/images"
	"server-database/internal/images/storer"
)

type Imager interface {
	Create(ctx context.Context, payload *images.CreateImagePayload) (int, error)
	Get(ctx context.Context, id int) (*images.Image, error)
	Delete(ctx context.Context, id int) error
}

type Service struct {
	log   *log.Logger
	store storer.Store
}

func New(logger *log.Logger, store storer.Store) *Service {
	return &Service{
		log:   logger,
		store: store,
	}
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}

func (s *Service) Get(ctx context.Context, id int) (*images.Image, error) {
	return s.store.Get(ctx, id)
}

func (s *Service) Create(ctx context.Context, payload *images.CreateImagePayload) (int, error) {
	if err := payload.Validate(); err != nil {
		return 0, err
	}

	img := &images.Image{
		Title:   payload.Title,
		AltText: payload.AltText,
		URL:     payload.URL,
	}

	err := s.store.Insert(img)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return 0, fmt.Errorf("%w: %w", images.ImagesUniqueCodeViolation, err)
			default:
				return 0, err
			}
		}
	}

	return img.ID, nil
}
