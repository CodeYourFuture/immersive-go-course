package api

import (
	"context"

	"github.com/CodeYourFuture/immersive-go-course/buggy-app/util"
)

type ApiService struct {
	util.Service
}

type Config struct {
	util.Config
}

func (as *ApiService) Run(ctx context.Context, config Config) error {
	return nil
}

func NewApiService() ApiService {
	return ApiService{}
}
