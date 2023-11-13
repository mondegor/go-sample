package usecase

import (
    "context"

    "github.com/mondegor/go-webcore/mrtype"
)

type (
    FileProviderAdapterService interface {
        Get(ctx context.Context, path string) (*mrtype.File, error)
    }
)
