package v1

import (
	"github.com/go-chi/chi/v5"
	"server-database/cmd/api/server"
)

func Register(mux *chi.Mux, svr *server.Server) {
	//image := NewImage(svr.Logger, svr.ImageService)
	images := NewImage(svr.Logger, svr.ImageService)

	mux.Get("/images", images.Get)
	mux.Get("/images/", images.List)
	mux.Delete("/images", images.Delete)
	mux.Post("/images", images.Post)

}
