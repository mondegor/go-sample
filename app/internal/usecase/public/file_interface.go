package usecase

import (
    "context"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    FileItemService interface {
        Get(ctx context.Context, path string) (*mrentity.File, error)
    }
)
