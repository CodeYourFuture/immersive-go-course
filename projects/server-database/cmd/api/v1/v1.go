package v1

import (
	"net/http"
	"server-database/cmd/api/server"
)

func Register(mux *http.ServeMux, svr *server.Server) {
	images := NewImages(svr.Logger, svr.ImageService)
	mux.Handle("/images", images)
}
