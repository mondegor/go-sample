package usecase

import (
    "context"
    "go-sample/internal/entity/public"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type (
    CatalogCategoryService interface {
        GetList(ctx context.Context, params entity.CatalogCategoryParams) ([]entity.CatalogCategory, int64, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogCategory, error)
    }

    CatalogCategoryStorage interface {
        NewFetchParams(params entity.CatalogCategoryParams) mrstorage.SqlSelectParams
        Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CatalogCategory, error)
        FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
        LoadOne(ctx context.Context, row *entity.CatalogCategory) error
    }
)
