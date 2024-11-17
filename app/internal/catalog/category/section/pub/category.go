package pub

import (
	"context"

	"github.com/google/uuid"

	"github.com/mondegor/go-sample/internal/catalog/category/section/pub/entity"
)

type (
	// CategoryUseCase - comment interface.
	CategoryUseCase interface {
		GetList(ctx context.Context, params entity.CategoryParams) (items []entity.Category, countItems uint64, err error)
		GetItem(ctx context.Context, itemID uuid.UUID, languageID uint16) (entity.Category, error)
	}

	// CategoryStorage - comment interface.
	CategoryStorage interface {
		FetchWithTotal(ctx context.Context, params entity.CategoryParams) (rows []entity.Category, countRows uint64, err error)
		FetchOne(ctx context.Context, rowID uuid.UUID) (entity.Category, error)
	}
)
