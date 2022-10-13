package util

import "context"

type Config interface{}

type Service interface {
	Run(ctx context.Context, config Config) error
}
