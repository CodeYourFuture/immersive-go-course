package server

import (
	"context"
	"database/sql"
	"log"
	"os"
	"server-database/internal/images/service"
	"server-database/internal/images/storer"
	"time"
)

type Server struct {
	Logger       *log.Logger
	DB           *sql.DB
	ImageService service.Imager
}

func (s *Server) MountLogger() error {
	s.Logger = log.New(os.Stdout, "images-api", log.LstdFlags)
	return nil
}

func (s *Server) MountImageService() error {
	imageStoreManager := storer.NewManager(s.DB)
	imageService := service.New(s.Logger, imageStoreManager)
	s.ImageService = imageService
	return nil
}

func (s *Server) MountDB() (*sql.DB, error) {
	db, err := openDB()
	if err != nil {
		s.Logger.Printf("error opening db: %v", err)
		os.Exit(1)
	}
	s.DB = db
	return db, err
}

func openDB() (*sql.DB, error) {
	// TODO: move connection string to the env file
	db, err := sql.Open("postgres", "postgres://root:password@localhost:5432/go-server-database?sslmode=disable")
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
