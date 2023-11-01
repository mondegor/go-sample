package usecase

import (
    "context"
    "go-sample/internal/entity/admin-panel"

    "github.com/mondegor/go-components/mrorderer"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrenum"
)

type (
    CatalogProductService interface {
        GetList(ctx context.Context, params entity.CatalogProductParams) ([]entity.CatalogProduct, int64, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32) (*entity.CatalogProduct, error)
        Create(ctx context.Context, item *entity.CatalogProduct) error
        Store(ctx context.Context, item *entity.CatalogProduct) error
        ChangeStatus(ctx context.Context, item *entity.CatalogProduct) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
        MoveAfterId(ctx context.Context, id mrentity.KeyInt32, afterId mrentity.KeyInt32) error
    }

    CatalogProductStorage interface {
        GetMetaData(categoryId mrentity.KeyInt32) mrorderer.EntityMeta
        NewFetchParams(params entity.CatalogProductParams) mrstorage.SqlSelectParams
        Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.CatalogProduct, error)
        FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
        LoadOne(ctx context.Context, row *entity.CatalogProduct) error
        FetchIdByArticle(ctx context.Context, article string) (mrentity.KeyInt32, error)
        FetchStatus(ctx context.Context, row *entity.CatalogProduct) (mrenum.ItemStatus, error)
        Insert(ctx context.Context, row *entity.CatalogProduct) error
        Update(ctx context.Context, row *entity.CatalogProduct) error
        UpdateStatus(ctx context.Context, row *entity.CatalogProduct) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
