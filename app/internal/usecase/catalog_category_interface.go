package usecase

import (
    "context"
    "go-sample/internal/entity"

    mrcom_status "github.com/mondegor/go-components/mrcom/status"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogCategoryService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogCategoryListFilter) ([]entity.CatalogCategory, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogCategory, error)
        CheckAvailability(ctx context.Context, id mrentity.KeyInt32) error
        Create(ctx context.Context, item *entity.CatalogCategory) error
        Store(ctx context.Context, item *entity.CatalogCategory) error
        ChangeStatus(ctx context.Context, item *entity.CatalogCategory) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogCategoryStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogCategoryListFilter, rows *[]entity.CatalogCategory) error
        LoadOne(ctx context.Context, row *entity.CatalogCategory) error
        FetchStatus(ctx context.Context, row *entity.CatalogCategory) (mrcom_status.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.CatalogCategory) error
        Update(ctx context.Context, row *entity.CatalogCategory) error
        UpdateStatus(ctx context.Context, row *entity.CatalogCategory) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
