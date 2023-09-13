package usecase

import (
    "context"
    "go-sample/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogTrademarkService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogTrademarkListFilter) ([]entity.CatalogTrademark, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogTrademark, error)
        Create(ctx context.Context, item *entity.CatalogTrademark) error
        Store(ctx context.Context, item *entity.CatalogTrademark) error
        ChangeStatus(ctx context.Context, item *entity.CatalogTrademark) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogTrademarkStorage interface {
        LoadAll(ctx context.Context, listFilter *entity.CatalogTrademarkListFilter, rows *[]entity.CatalogTrademark) error
        LoadOne(ctx context.Context, row *entity.CatalogTrademark) error
        FetchStatus(ctx context.Context, row *entity.CatalogTrademark) (mrcom.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.CatalogTrademark) error
        Update(ctx context.Context, row *entity.CatalogTrademark) error
        UpdateStatus(ctx context.Context, row *entity.CatalogTrademark) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
