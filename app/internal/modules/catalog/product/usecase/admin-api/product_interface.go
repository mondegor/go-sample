package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/product/entity/admin-api"

	"github.com/google/uuid"
	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ProductUseCase interface {
		GetList(ctx context.Context, params entity.ProductParams) ([]entity.Product, int64, error)
		GetItem(ctx context.Context, itemID mrtype.KeyInt32) (entity.Product, error)
		Create(ctx context.Context, item entity.Product) (mrtype.KeyInt32, error)
		Store(ctx context.Context, item entity.Product) error
		ChangeStatus(ctx context.Context, item entity.Product) error
		Remove(ctx context.Context, itemID mrtype.KeyInt32) error
		MoveAfterID(ctx context.Context, itemID mrtype.KeyInt32, afterID mrtype.KeyInt32) error
	}

	ProductStorage interface {
		NewOrderMeta(categoryID uuid.UUID) mrorderer.EntityMeta
		NewSelectParams(params entity.ProductParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Product, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID mrtype.KeyInt32) (entity.Product, error)
		FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, row entity.Product) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, rowID mrtype.KeyInt32) error
		Insert(ctx context.Context, row entity.Product) (mrtype.KeyInt32, error)
		Update(ctx context.Context, row entity.Product) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Product) (int32, error)
		Delete(ctx context.Context, rowID mrtype.KeyInt32) error
	}
)
