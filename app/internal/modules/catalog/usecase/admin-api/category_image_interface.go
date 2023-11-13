package usecase

import (
    "context"

    "github.com/mondegor/go-webcore/mrtype"
)

type (
    CategoryImageService interface {
        // Get - WARNING you don't forget to call item.Body.Close()
        Get(ctx context.Context, categoryID mrtype.KeyInt32) (*mrtype.File, error)
        GetInfoByPath(ctx context.Context, imagePath string) (*mrtype.FileInfo, error)
        Store(ctx context.Context, categoryID mrtype.KeyInt32, file *mrtype.File) error
        Remove(ctx context.Context, categoryID mrtype.KeyInt32) error
    }

    CategoryImageStorage interface {
        FetchPath(ctx context.Context, categoryID mrtype.KeyInt32) (string, error)
        Update(ctx context.Context, categoryID mrtype.KeyInt32, path string) error
        Delete(ctx context.Context, categoryID mrtype.KeyInt32) error
    }
)
