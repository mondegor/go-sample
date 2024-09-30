package factory

import (
	"context"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrlog"

	"github.com/mondegor/go-sample/config"
)

// NewFileProviderPool - создаёт объект mrstorage.FileProviderPool.
func NewFileProviderPool(ctx context.Context, cfg config.Config) (*mrstorage.FileProviderPool, error) {
	mrlog.Ctx(ctx).Info().Msg("Create and init file provider pool")

	pool := mrstorage.NewFileProviderPool()

	// fsAdapter := NewFileSystem(ctx, cfg)

	// if err := RegisterFileImageStorage(ctx, cfg, pool, fsAdapter); err != nil {
	//   return nil, err
	// }

	minioAdapter, err := NewS3Minio(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err = RegisterS3ImageStorage(ctx, cfg, pool, minioAdapter); err != nil {
		return nil, err
	}

	return pool, pool.Ping(ctx)
}
