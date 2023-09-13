package usecase

import (
    "context"
    "go-sample/internal/entity"

    "github.com/mondegor/go-components/mrcom"
    mrcom_orderer "github.com/mondegor/go-components/mrcom/orderer"
    "github.com/mondegor/go-storage/mrentity"
)

type (
    CatalogProductService interface {
        GetList(ctx context.Context, listFilter *entity.CatalogProductListFilter) ([]entity.CatalogProduct, error)
        GetItem(ctx context.Context, id mrentity.KeyInt32, categoryId mrentity.KeyInt32) (*entity.CatalogProduct, error)
        Create(ctx context.Context, item *entity.CatalogProduct) error
        Store(ctx context.Context, item *entity.CatalogProduct) error
        ChangeStatus(ctx context.Context, item *entity.CatalogProduct) error
        Remove(ctx context.Context, id mrentity.KeyInt32) error
        MoveAfterId(ctx context.Context, id mrentity.KeyInt32, afterId mrentity.KeyInt32, categoryId mrentity.KeyInt32) error
    }

    CatalogProductStorage interface {
        GetMetaData(formId mrentity.KeyInt32) mrcom_orderer.EntityMeta
        LoadAll(ctx context.Context, listFilter *entity.CatalogProductListFilter, rows *[]entity.CatalogProduct) error
        LoadOne(ctx context.Context, row *entity.CatalogProduct) error
        FetchIdByArticle(ctx context.Context, article string) (mrentity.KeyInt32, error)
        FetchStatus(ctx context.Context, row *entity.CatalogProduct) (mrcom.ItemStatus, error)
        Insert(ctx context.Context, row *entity.CatalogProduct) error
        Update(ctx context.Context, row *entity.CatalogProduct) error
        UpdateStatus(ctx context.Context, row *entity.CatalogProduct) error
        Delete(ctx context.Context, id mrentity.KeyInt32) error
    }
)
