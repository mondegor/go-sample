package usecase

import (
    "context"
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    CatalogCategoryService interface {
        GetList(ctx context.Context, params entity.CatalogCategoryParams) ([]entity.CatalogCategory, int64, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogCategory, error)
        CheckAvailability(ctx context.Context, id mrentity.KeyInt32) error
        Create(ctx context.Context, item *entity.CatalogCategory) error
        Store(ctx context.Context, item *entity.CatalogCategory) error
        ChangeStatus(ctx context.Context, item *entity.CatalogCategory) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogCategoryStorage interface {
        NewFetchParams(params entity.CatalogCategoryParams) mrstorage.SqlSelectParams
        Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CatalogCategory, error)
        FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
        LoadOne(ctx context.Context, row *entity.CatalogCategory) error
        FetchStatus(ctx context.Context, row *entity.CatalogCategory) (mrenum.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.CatalogCategory) error
        Update(ctx context.Context, row *entity.CatalogCategory) error
        UpdateStatus(ctx context.Context, row *entity.CatalogCategory) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
