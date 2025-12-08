package util

import (
	"context"
	"github.com/google/uuid"
)

func NewContext() context.Context {
	ctx := context.Background()
	return context.WithValue(ctx, "request_id", uuid.New().String())
}
