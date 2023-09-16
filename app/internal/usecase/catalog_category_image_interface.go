package usecase

import (
    "context"
    "go-sample/internal/entity"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogCategoryImageService interface {
        // Load - WARNING you don't forget to call item.File.Body.Close()
        Load(ctx context.Context, item *entity.CatalogCategoryImageObject) error
        Store(ctx context.Context, item *entity.CatalogCategoryImageObject) error
        Remove(ctx context.Context, categoryId mrentity.KeyInt32) error
    }

    CatalogCategoryImageStorage interface {
        Fetch(ctx context.Context, categoryId mrentity.KeyInt32) (string, error)
        Update(ctx context.Context, categoryId mrentity.KeyInt32, imagePath string) error
        Delete(ctx context.Context, categoryId mrentity.KeyInt32) error
    }
)
