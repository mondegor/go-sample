package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/category/entity/admin-api"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryService interface {
		GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error)
		GetItem(ctx context.Context, id mrtype.KeyInt32) (*entity.Category, error)
		Create(ctx context.Context, item *entity.Category) error
		Store(ctx context.Context, item *entity.Category) error
		ChangeStatus(ctx context.Context, item *entity.Category) error
		Remove(ctx context.Context, id mrtype.KeyInt32) error
	}

	CategoryStorage interface {
		NewFetchParams(params entity.CategoryParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		LoadOne(ctx context.Context, row *entity.Category) error
		FetchStatus(ctx context.Context, row *entity.Category) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, id mrtype.KeyInt32) error
		Insert(ctx context.Context, row *entity.Category) error
		Update(ctx context.Context, row *entity.Category) (int32, error)
		UpdateStatus(ctx context.Context, row *entity.Category) (int32, error)
		Delete(ctx context.Context, id mrtype.KeyInt32) error
	}
)
