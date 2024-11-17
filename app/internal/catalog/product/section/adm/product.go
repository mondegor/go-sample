package adm

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
)

type (
	// ProductUseCase - comment interface.
	ProductUseCase interface {
		GetList(ctx context.Context, params entity.ProductParams) (items []entity.Product, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uint64) (entity.Product, error)
		Create(ctx context.Context, item entity.Product) (itemID uint64, err error)
		Store(ctx context.Context, item entity.Product) error
		ChangeStatus(ctx context.Context, item entity.Product) error
		Remove(ctx context.Context, itemID uint64) error
		MoveAfterID(ctx context.Context, itemID, afterID uint64) error
	}

	// ProductStorage - comment interface.
	ProductStorage interface {
		NewCondition(categoryID uuid.UUID) mrstorage.SQLPartFunc
		FetchWithTotal(ctx context.Context, params entity.ProductParams) (rows []entity.Product, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uint64) (entity.Product, error)
		FetchIDByArticle(ctx context.Context, article string) (rowID uint64, err error)
		FetchStatus(ctx context.Context, rowID uint64) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Product) (rowID uint64, err error)
		Update(ctx context.Context, row entity.Product) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.Product) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uint64) error
	}
)
