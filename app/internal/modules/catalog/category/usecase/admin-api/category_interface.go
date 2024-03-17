package usecase

import (
	"context"
	"go-sample/internal/modules/catalog/category/entity/admin-api"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
)

type (
	CategoryUseCase interface {
		GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error)
		GetItem(ctx context.Context, itemID uuid.UUID) (entity.Category, error)
		Create(ctx context.Context, item entity.Category) (uuid.UUID, error)
		Store(ctx context.Context, item entity.Category) error
		ChangeStatus(ctx context.Context, item entity.Category) error
		Remove(ctx context.Context, itemID uuid.UUID) error
	}

	CategoryStorage interface {
		NewFetchParams(params entity.CategoryParams) mrstorage.SqlSelectParams
		Fetch(ctx context.Context, params mrstorage.SqlSelectParams) ([]entity.Category, error)
		FetchTotal(ctx context.Context, where mrstorage.SqlBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error)
		FetchStatus(ctx context.Context, row entity.Category) (mrenum.ItemStatus, error)
		IsExists(ctx context.Context, rowID uuid.UUID) error
		Insert(ctx context.Context, row entity.Category) (uuid.UUID, error)
		Update(ctx context.Context, row entity.Category) (int32, error)
		UpdateStatus(ctx context.Context, row entity.Category) (int32, error)
		Delete(ctx context.Context, rowID uuid.UUID) error
	}
)
