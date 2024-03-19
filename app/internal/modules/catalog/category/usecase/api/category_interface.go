package usecase_api

import (
	"context"

	"github.com/google/uuid"
)

type (
	CategoryStorage interface {
		IsExists(ctx context.Context, rowID uuid.UUID) error
	}
)
