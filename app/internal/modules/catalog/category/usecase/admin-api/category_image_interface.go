package usecase

import (
	"context"

	"github.com/mondegor/go-storage/mrentity"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	CategoryImageUseCase interface {
		// GetFile - WARNING you don't forget to call item.Body.Close()
		GetFile(ctx context.Context, categoryID mrtype.KeyInt32) (mrtype.Image, error)
		StoreFile(ctx context.Context, categoryID mrtype.KeyInt32, image mrtype.Image) error
		RemoveFile(ctx context.Context, categoryID mrtype.KeyInt32) error
	}

	CategoryImageStorage interface {
		FetchMeta(ctx context.Context, categoryID mrtype.KeyInt32) (mrentity.ImageMeta, error)
		UpdateMeta(ctx context.Context, categoryID mrtype.KeyInt32, meta mrentity.ImageMeta) error
		DeleteMeta(ctx context.Context, categoryID mrtype.KeyInt32) error
	}
)
