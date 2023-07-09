package storer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"server-database/internal/images"
)

const (
	tableName = "public.images"
)

type Store interface {
	Insert(*images.Image) error
	Get(ctx context.Context, id int) (*images.Image, error)
	Delete(ctx context.Context, id int) error
}

type Manager struct {
	Reader
	Writer
}

type Reader struct {
	db *sql.DB
}

type Writer struct {
	db *sql.DB
}

func NewManager(db *sql.DB) *Manager {
	r := Reader{
		db: db,
	}

	w := Writer{
		db: db,
	}

	return &Manager{
		r, w,
	}
}

func (m *Reader) Get(ctx context.Context, id int) (*images.Image, error) {
	if id < 1 {
		return nil, errors.New("record not found")
	}

	query := `
	SELECT id, title, url, alt_text, created_at, resolution 
	FROM public.images
	WHERE id = $1`

	var image images.Image

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := m.db.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.Title,
		&image.URL,
		&image.AltText,
		&image.CreatedAt,
		&image.Resolution,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, images.ImagesNotFound
		default:
			return nil, fmt.Errorf("unexpected error happened: %w", err)
		}
	}

	return &image, nil
}

func (m *Reader) List() error {
	return nil
}

func (m *Writer) Insert(image *images.Image) error {
	query := fmt.Sprintf(`
	INSERT INTO public.images(title, url, alt_text, resolution) 
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at
	`)

	args := []any{image.Title, image.URL, image.AltText, image.Resolution}
	return m.db.QueryRow(query, args...).Scan(&image.ID, &image.CreatedAt)
}

func (m *Writer) Update() error {
	return nil
}

func (m *Writer) Delete(ctx context.Context, id int) error {
	query := fmt.Sprintf(`
	DELETE FROM public.images 
	WHERE id = $1
	`)

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	result, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return images.ImagesNotFound
	}

	return nil
}
