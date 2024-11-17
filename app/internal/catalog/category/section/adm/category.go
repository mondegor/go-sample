package adm

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrenum"

	"github.com/mondegor/go-sample/internal/catalog/category/section/adm/entity"
)

type (
	// CategoryUseCase - comment interface.
	CategoryUseCase interface {
		GetList(ctx context.Context, params entity.CategoryParams) (items []entity.Category, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uuid.UUID) (entity.Category, error)
		Create(ctx context.Context, item entity.Category) (itemID uuid.UUID, err error)
		Store(ctx context.Context, item entity.Category) error
		ChangeStatus(ctx context.Context, item entity.Category) error
		Remove(ctx context.Context, itemID uuid.UUID) error
	}

	// CategoryStorage - comment interface.
	CategoryStorage interface {
		FetchWithTotal(ctx context.Context, params entity.CategoryParams) (rows []entity.Category, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error)
		FetchStatus(ctx context.Context, rowID uuid.UUID) (mrenum.ItemStatus, error)
		Insert(ctx context.Context, row entity.Category) (rowID uuid.UUID, err error)
		Update(ctx context.Context, row entity.Category) (tagVersion uint32, err error)
		UpdateStatus(ctx context.Context, row entity.Category) (tagVersion uint32, err error)
		Delete(ctx context.Context, rowID uuid.UUID) error
	}
)
