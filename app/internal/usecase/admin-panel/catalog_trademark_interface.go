package usecase

import (
    "context"
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    CatalogTrademarkService interface {
        GetList(ctx context.Context, params entity.CatalogTrademarkParams) ([]entity.CatalogTrademark, int64, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogTrademark, error)
        Create(ctx context.Context, item *entity.CatalogTrademark) error
        Store(ctx context.Context, item *entity.CatalogTrademark) error
        ChangeStatus(ctx context.Context, item *entity.CatalogTrademark) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
    }

    CatalogTrademarkStorage interface {
        NewFetchParams(params entity.CatalogTrademarkParams) mrstorage.SqlSelectParams
        Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CatalogTrademark, error)
        FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
        LoadOne(ctx context.Context, row *entity.CatalogTrademark) error
        FetchStatus(ctx context.Context, row *entity.CatalogTrademark) (mrenum.ItemStatus, error)
        IsExists(ctx context.Context, id mrentity.KeyInt32) error
        Insert(ctx context.Context, row *entity.CatalogTrademark) error
        Update(ctx context.Context, row *entity.CatalogTrademark) error
        UpdateStatus(ctx context.Context, row *entity.CatalogTrademark) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
