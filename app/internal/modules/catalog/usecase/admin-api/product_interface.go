package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/entity/admin-api"

	"github.com/mondegor/go-components/mrorderer"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ProductService interface {
		GetList(ctx context.Context, params entity.ProductParams) ([]entity.Product, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Product, error)
		Create(ctx context.Context, item *entity.Product) error
		Store(ctx context.Context, item *entity.Product) error
		ChangeStatus(ctx context.Context, item *entity.Product) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
		MoveAfterID(ctx context.Context, id mrtype.KeyInt32, afterID mrtype.KeyInt32) error
	}

	ProductStorage interface {
		GetMetaData(categoryID mrtype.KeyInt32) mrorderer.EntityMeta
		NewFetchParams(params entity.ProductParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Product, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.Product) error
		FetchIdByArticle(ctx context.Context, article string) (mrtype.KeyInt32, error)
		FetchStatus(ctx context.Context, row *entity.Product) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row *entity.Product) error
		Update(ctx context.Context, row *entity.Product) error
		UpdateStatus(ctx context.Context, row *entity.Product) error
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
