package pub

import (
	"context"

	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// CategoryUseCase - comment interface.
	CategoryUseCase interface {
		GetList(ctx context.Context, params entity.CategoryParams) ([]entity.Category, int64, error)
		GetItem(ctx context.Context, itemID uuid.UUID, languageID uint16) (entity.Category, error)
	}

	// CategoryStorage - comment interface.
	CategoryStorage interface {
		NewSelectParams(params entity.CategoryParams) mrstorage.SQLSelectParams
		Fetch(ctx context.Context, params mrstorage.SQLSelectParams) ([]entity.Category, error)
		FetchTotal(ctx context.Context, where mrstorage.SQLBuilderPart) (int64, error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error)
	}
)
