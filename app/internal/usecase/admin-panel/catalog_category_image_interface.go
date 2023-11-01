package usecase

import (
    "context"
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogCategoryImageService interface {
        // Get - WARNING you don't forget to call item.File.Body.Close()
        Get(ctx context.Context, categoryId mrentity.KeyInt32) (*entity.CatalogCategoryImageObject, error)
        Store(ctx context.Context, item *entity.CatalogCategoryImageObject) error
        Remove(ctx context.Context, categoryId mrentity.KeyInt32) error
    }

    CatalogCategoryImageStorage interface {
        FetchOne(ctx context.Context, categoryId mrentity.KeyInt32) (string, error)
        Update(ctx context.Context, categoryId mrentity.KeyInt32, imagePath string) error
        Delete(ctx context.Context, categoryId mrentity.KeyInt32) error
    }
)
