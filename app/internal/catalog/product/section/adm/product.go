package adm

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-sample/internal/catalog/product/section/adm/entity"
)

type (
	// ProductUseCase - comment interface.
	ProductUseCase interface {
		GetList(ctx context.Context, params entity.ProductParams) ([]entity.Product, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Product, error)
		Create(ctx context.Context, item entity.Product) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Product) error
		ChangeStatus(ctx context.Context, item entity.Product) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
		MoveAfterID(ctx context.Context, itemID, afterID mrtype.KeyInt32) error
	}

	// ProductStorage - comment interface.
	ProductStorage interface {
		NewOrderMeta(categoryID uuid.UUID) mrstorage.MetaGetter
		NewSelectParams(params entity.ProductParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Product, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Product, error)
		FetchIDByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, rowID mrtype.KeyInt32) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Product) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Product) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Product) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
