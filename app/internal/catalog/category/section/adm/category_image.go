package adm

import (
	"context"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// CategoryImageUseCase - comment interface.
	CategoryImageUseCase interface {
		// GetFile - WARNING you don't forget to call item.Body.Close()
		GetFile(ctx context.Context, categoryID uuid.UUID) (mrtype.Image, error)
		StoreFile(ctx context.Context, categoryID uuid.UUID, image mrtype.Image) error
		RemoveFile(ctx context.Context, categoryID uuid.UUID) error
	}

	// CategoryImageStorage - comment interface.
	CategoryImageStorage interface {
		FetchMeta(ctx context.Context, categoryID uuid.UUID) (mrentity.ImageMeta, error)
		UpdateMeta(ctx context.Context, categoryID uuid.UUID, meta mrentity.ImageMeta) error
		DeleteMeta(ctx context.Context, categoryID uuid.UUID) error
	}
)
