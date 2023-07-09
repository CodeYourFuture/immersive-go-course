package v1

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"server-database/internal/images"
	"server-database/internal/images/service"
)

type queryParams = string

const (
	indent queryParams = "indent"
)

type Images struct {
	logger  *log.Logger
	service service.Imager
}

// NewImages is a constructor of the images
func NewImages(log *log.Logger, svc service.Imager) *Images {
	return &Images{
		logger:  log,
		service: svc,
	}
}

func (i *Images) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		i.Get(w, req)
	case http.MethodPost:
		i.Post(w, req)
	case http.MethodDelete:
		i.Delete(w, req)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (i *Images) Delete(w http.ResponseWriter, request *http.Request) {
	id, err := fetchId(request)
	if err != nil {
		i.logger.Printf("error fetching id: %v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = i.service.Delete(request.Context(), id)
	switch err {
	case nil:
		Respond(i.logger, w, nil, http.StatusNoContent, nil)
		return
	case images.ImagesNotFound:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (i *Images) Get(w http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	image, err := i.service.Get(request.Context(), idInt)
	switch err {
	case nil:
	case images.ImagesNotFound:
		http.Error(w, "not found", http.StatusNotFound)
		return
	default:
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	queryParams := request.URL.Query()
	indent, ok := queryParams[indent]
	if ok {
		indent, err := strconv.Atoi(indent[0])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		identStr := strings.Repeat(" ", indent)
		Respond(i.logger, w, image, http.StatusOK, &identStr)
		return
	}

	Respond(i.logger, w, image, http.StatusOK, nil)
}

func (i *Images) Post(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()

	var payload images.CreateImagePayload
	if err := decoder.Decode(&payload); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	id, err := i.service.Create(ctx, &payload)
	switch {
	case err == nil:
		Respond(i.logger, w, map[string]int{
			"id": id,
		}, http.StatusCreated, nil)
	case errors.Is(err, images.ImagesUniqueCodeViolation):
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func fetchId(request *http.Request) (int, error) {
	id := request.URL.Query().Get("id")
	if id == "" {
		return 0, errors.New("id not found")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	return idInt, nil
}
