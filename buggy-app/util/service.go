package util

import (
	"context"
)

type Service interface {
	Run(ctx context.Context, config any) error
}
