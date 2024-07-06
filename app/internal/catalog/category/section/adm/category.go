package adm

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/entity"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	// CategoryUseCase - comment interface.
	CategoryUseCase interface {
		GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error)
		GetItem(ctx context.Context, itemID uuid.UUID) (entity.Category, error)
		Create(ctx context.Context, item entity.Category) (uuid.UUID, error)
		Store(ctx context.Context, item entity.Category) error
		ChangeStatus(ctx context.Context, item entity.Category) error
		Remove(ctx context.Context, itemID uuid.UUID) error
	}

	// CategoryStorage - comment interface.
	CategoryStorage interface {
		NewSelectParams(params entity.CategoryParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Category, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error)
		FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Category) (uuid.UUID, error)
		Update(ctx context.Context, row entity.Category) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Category) (int32, error)
		Delete(ctx context.Context, rowID uuid.UUID) error
	}
)
